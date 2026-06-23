package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var IST = time.FixedZone("IST", 5*3600+30*60)

// StartCompetitionScheduler initializes the background loop and seeds initial competition records
func StartCompetitionScheduler() {
	log.Println("Starting Competition Scheduler...")
	SeedActiveCompetitions()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		var lastRunMin int = -1

		for range ticker.C {
			now := time.Now().In(IST)
			if now.Minute() == lastRunMin {
				continue
			}
			lastRunMin = now.Minute()

			// 1. Weekly Quest Rotation Check (Friday 20:00 IST)
			if now.Weekday() == time.Friday && now.Hour() == 20 && now.Minute() == 0 {
				RotateWeeklyQuest(now)
			}

			// 2. Weekly Quest Reminder (Friday 18:00 IST)
			if now.Weekday() == time.Friday && now.Hour() == 18 && now.Minute() == 0 {
				SendWeeklyQuestReminder()
			}

			// 3. Daily Challenge Rotation Check (Everyday 00:00 IST)
			if now.Hour() == 0 && now.Minute() == 0 {
				RotateDailyChallenge(now)
			}

			// 4. Daily Challenge Warning Check (Everyday 20:00 IST)
			if now.Hour() == 20 && now.Minute() == 0 {
				SendDailyChallengeWarning()
			}
		}
	}()
}

// SeedActiveCompetitions ensures there is always an active Daily Challenge and Weekly Quest for testing
func SeedActiveCompetitions() {
	now := time.Now().In(IST)

	// Seed Daily Challenge if missing
	var dailyCount int64
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)
	DB.Model(&DailyChallenge{}).Where("date = ?", todayMidnight).Count(&dailyCount)
	if dailyCount == 0 {
		flag := fmt.Sprintf("ANTARIKSHA{daily_%s}", now.Format("2006_01_02"))
		dc := DailyChallenge{
			ID:          uuid.New(),
			Date:        todayMidnight,
			Description: "Analyze the log entries in /var/log/proxy/access.log and find the source IP of Shiva's command payload.",
			FlagHash:    hashFlag(flag),
		}
		DB.Create(&dc)
		log.Printf("[SEED] Created active Daily Challenge. Flag: %s", flag)
	}

	// Seed Weekly Quest if missing
	var weeklyCount int64
	DB.Model(&WeeklyQuest{}).Where("starts_at <= ? AND ends_at >= ?", now, now).Count(&weeklyCount)
	if weeklyCount == 0 {
		// Create a quest starting last Friday 8 PM to next Friday 8 PM for sandbox testing
		starts := now.AddDate(0, 0, -int(now.Weekday())) // past Sunday or earlier
		starts = time.Date(starts.Year(), starts.Month(), starts.Day(), 20, 0, 0, 0, IST)
		ends := starts.AddDate(0, 0, 7) // next Friday

		flag := "ANTARIKSHA{weekly_satellite_link_restored}"
		wq := WeeklyQuest{
			ID:          uuid.New(),
			Title:       "Weekly Quest: Restore Compromised Satellite Link",
			Description: "A malware agent has disabled the primary downlink. Solve the proxy logs routing to restore signal.",
			FlagHash:    hashFlag(flag),
			StartsAt:    starts,
			EndsAt:      ends,
		}
		DB.Create(&wq)
		log.Printf("[SEED] Created active Weekly Quest. Flag: %s", flag)
	}
}

// RotateWeeklyQuest closes old quests and starts a new one
func RotateWeeklyQuest(now time.Time) {
	log.Println("[SCHEDULER] Rotating Weekly Quest...")

	// 1. Create new quest
	endsAt := now.Add(48 * time.Hour) // Closes Sunday 8:00 PM IST
	secretCode := uuid.New().String()[:8]
	flag := fmt.Sprintf("ANTARIKSHA{weekly_quest_%s_%s}", now.Format("2006_01_02"), secretCode)

	newQuest := WeeklyQuest{
		ID:          uuid.New(),
		Title:       fmt.Sprintf("Operation Antariksha - Weekly Quest %s", now.Format("2006-01-02")),
		Description: "Verify system integrity, track the Shiva SUID backdoor, and recover the telemetry signal coordinates.",
		FlagHash:    hashFlag(flag),
		StartsAt:    now,
		EndsAt:      endsAt,
	}

	if err := DB.Create(&newQuest).Error; err != nil {
		log.Printf("Failed to rotate Weekly Quest: %v", err)
		return
	}

	log.Printf("[SCHEDULER] Weekly Quest rotated successfully. Flag: %s", flag)
}

// RotateDailyChallenge creates a daily challenge for the new day
func RotateDailyChallenge(now time.Time) {
	log.Println("[SCHEDULER] Rotating Daily Challenge...")

	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)
	flag := fmt.Sprintf("ANTARIKSHA{daily_%s}", now.Format("2006_01_02"))

	newChallenge := DailyChallenge{
		ID:          uuid.New(),
		Date:        todayMidnight,
		Description: "Decrypt the Shiva beacon message from the local cron jobs payload.",
		FlagHash:    hashFlag(flag),
	}

	if err := DB.Create(&newChallenge).Error; err != nil {
		log.Printf("Failed to rotate Daily Challenge: %v", err)
		return
	}

	// Also perform streak decay check: users who didn't complete yesterday's challenge lose their streak
	yesterdayMidnight := todayMidnight.AddDate(0, 0, -1)
	var yesterdayChallenge DailyChallenge
	if err := DB.Where("date = ?", yesterdayMidnight).First(&yesterdayChallenge).Error; err == nil {
		// Find users who have streak > 0 but didn't complete yesterday's challenge
		var users []User
		DB.Where("streak > 0").Find(&users)
		for _, user := range users {
			var progress DailyChallengeProgress
			err := DB.Where("user_id = ? AND daily_challenge_id = ? AND completed = ?", user.ID, yesterdayChallenge.ID, true).First(&progress).Error
			if err != nil {
				// No completion found, reset streak!
				user.Streak = 0
				DB.Save(&user)
				log.Printf("[SCHEDULER] Streak reset for user: %s due to inactivity.", user.Username)
			}
		}
	}

	log.Printf("[SCHEDULER] Daily Challenge rotated successfully. Flag: %s", flag)
}

// SendWeeklyQuestReminder notifies participants about upcoming quest
func SendWeeklyQuestReminder() {
	var users []User
	DB.Find(&users)
	log.Printf("[MAIL QUEUE] Dispatched Weekly Quest Reminder to %d participants.", len(users))
	for _, u := range users {
		log.Printf("[EMAIL] To: %s | Subject: Operation Antariksha - Weekly Quest Begins at 20:00 IST!", u.Email)
	}
}

// SendDailyChallengeWarning alerts users whose daily streak is at risk
func SendDailyChallengeWarning() {
	now := time.Now().In(IST)
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)

	var activeChallenge DailyChallenge
	if err := DB.Where("date = ?", todayMidnight).First(&activeChallenge).Error; err != nil {
		return
	}

	var users []User
	DB.Where("streak > 0").Find(&users)

	for _, u := range users {
		var progress DailyChallengeProgress
		err := DB.Where("user_id = ? AND daily_challenge_id = ? AND completed = ?", u.ID, activeChallenge.ID, true).First(&progress).Error
		if err != nil {
			// Not completed today yet! Send streak warning.
			log.Printf("[EMAIL] To: %s | WARNING: Complete today's Daily Challenge to save your %d-day streak! Resets in 4 hours.", u.Email, u.Streak)
		}
	}
}

// CalculateQuestScore computes score according to user requirements
func CalculateQuestScore(startedAt time.Time, completedAt time.Time, hints int, wrongAttempts int) float64 {
	durationSeconds := completedAt.Sub(startedAt).Seconds()
	if durationSeconds < 1 {
		durationSeconds = 1
	}
	score := durationSeconds * (1.0 + 0.1*float64(hints)) * (1.0 + 0.05*float64(wrongAttempts))
	return score
}
