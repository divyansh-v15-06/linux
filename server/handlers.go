package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Difficulty rating mapping
func getDifficultyRating(difficulty string) int {
	switch strings.ToLower(difficulty) {
	case "newcomer":
		return 800
	case "beginner":
		return 1000
	case "beginner+":
		return 1100
	case "intermediate":
		return 1200
	case "intermediate-advanced":
		return 1400
	case "advanced":
		return 1600
	case "advanced-pro":
		return 1750
	case "pro":
		return 1900
	default:
		return 1000
	}
}

// Elo update calculator
func calculateEloChange(playerRating, challengeRating int, result string, hintsUsed int) int {
	K := 32.0
	expected := 1.0 / (1.0 + math.Pow(10, float64(challengeRating-playerRating)/400.0))
	score := 0.0
	if result == "win" {
		score = math.Max(0.5, 1.0-float64(hintsUsed)*0.1)
	}
	return playerRating + int(math.Round(K*(score-expected)))
}

// XP reward helper
func getXPReward(difficulty string) int {
	switch strings.ToLower(difficulty) {
	case "newcomer":
		return 300
	case "beginner", "beginner+":
		return 500
	case "intermediate", "intermediate-advanced":
		return 800
	case "advanced", "advanced-pro":
		return 1000
	case "pro":
		return 1200
	default:
		return 500
	}
}

// GetCampaigns lists active campaigns
func GetCampaigns(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user token"})
		return
	}

	// Calculate overall campaign progress
	var totalChapters int64
	DB.Model(&Chapter{}).Where("campaign_slug = ?", "antariksha").Count(&totalChapters)

	var completedChapters int64
	DB.Model(&UserProgress{}).
		Joins("join chapters on chapters.id = user_progresses.chapter_id").
		Where("user_progresses.user_id = ? and user_progresses.status = ? and chapters.campaign_slug = ?", userID, "complete", "antariksha").
		Count(&completedChapters)

	c.JSON(http.StatusOK, gin.H{
		"campaigns": []gin.H{
			{
				"name":        "Operation Antariksha",
				"slug":        "antariksha",
				"description": "A rogue AI designation S.H.I.V.A has hijacked India's satellite network.",
				"total_nodes": totalChapters,
				"completed":   completedChapters,
			},
		},
	})
}

// GetCampaignBySlug lists all chapters and user progress inside a campaign
func GetCampaignBySlug(c *gin.Context) {
	slug := c.Param("slug")
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user token"})
		return
	}

	var chapters []Chapter
	err = DB.Where("campaign_slug = ?", slug).Order("number asc").Find(&chapters).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chapters"})
		return
	}

	// Fetch progress map
	var progressRecords []UserProgress
	DB.Where("user_id = ?", userID).Find(&progressRecords)
	progressMap := make(map[string]string) // chapterID -> status
	for _, p := range progressRecords {
		progressMap[p.ChapterID.String()] = p.Status
	}

	type ChapterProgressResponse struct {
		Number     int    `json:"number"`
		Title      string `json:"title"`
		City       string `json:"city"`
		Difficulty string `json:"difficulty"`
		Status     string `json:"status"` // locked, active, complete
	}

	var list []ChapterProgressResponse
	for _, ch := range chapters {
		status := progressMap[ch.ID.String()]
		if status == "" {
			if ch.Number == 0 {
				status = "active" // Bootcamp starts unlocked
			} else {
				status = "locked"
			}
		}
		list = append(list, ChapterProgressResponse{
			Number:     ch.Number,
			Title:      ch.Title,
			City:       ch.City,
			Difficulty: ch.Difficulty,
			Status:     status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"campaign_slug": slug,
		"chapters":      list,
	})
}

// GetChapterDetail fetches details for a specific chapter without exposing the flag hash
func GetChapterDetail(c *gin.Context) {
	slug := c.Param("slug")
	numStr := c.Param("num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter number"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user token"})
		return
	}

	var chapter Chapter
	err = DB.Where("campaign_slug = ? AND number = ?", slug, num).First(&chapter).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	// Enforce progression check
	if num > 0 {
		var prevChapter Chapter
		err = DB.Where("campaign_slug = ? AND number = ?", slug, num-1).First(&prevChapter).Error
		if err == nil {
			var prevProgress UserProgress
			err = DB.Where("user_id = ? AND chapter_id = ?", userID, prevChapter.ID).First(&prevProgress).Error
			if err != nil || prevProgress.Status != "complete" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Chapter is locked. Complete the previous chapter first."})
				return
			}
		}
	}

	// Return data safely
	c.JSON(http.StatusOK, gin.H{
		"number":      chapter.Number,
		"title":       chapter.Title,
		"city":        chapter.City,
		"difficulty":  chapter.Difficulty,
		"commands":    strings.Split(chapter.Commands, ","),
		"story_text":  chapter.StoryText,
	})
}

// SubmitFlag validates submitted flag and awards XP / ELO
func SubmitFlag(c *gin.Context) {
	slug := c.Param("slug")
	numStr := c.Param("num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter number"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user token"})
		return
	}

	var body struct {
		Flag string `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flag parameter is required"})
		return
	}

	var chapter Chapter
	err = DB.Where("campaign_slug = ? AND number = ?", slug, num).First(&chapter).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	// Fetch or initialize progress
	var progress UserProgress
	err = DB.Where("user_id = ? AND chapter_id = ?", userID, chapter.ID).First(&progress).Error
	if err == gorm.ErrRecordNotFound {
		// Verify progression
		if num > 0 {
			var prevChapter Chapter
			if DB.Where("campaign_slug = ? AND number = ?", slug, num-1).First(&prevChapter).Error == nil {
				var prevProgress UserProgress
				if DB.Where("user_id = ? AND chapter_id = ? AND status = ?", userID, prevChapter.ID, "complete").First(&prevProgress).Error != nil {
					c.JSON(http.StatusForbidden, gin.H{"error": "Cannot submit flag for a locked chapter"})
					return
				}
			}
		}

		progress = UserProgress{
			UserID:    userID,
			ChapterID: chapter.ID,
			Status:    "active",
		}
		DB.Create(&progress)
	}

	submittedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(body.Flag)))
	isCorrect := submittedHash == chapter.FlagHash

	// Record submission
	submission := Submission{
		ID:            uuid.New(),
		UserID:        userID,
		ChapterID:     chapter.ID,
		FlagSubmitted: body.Flag,
		IsCorrect:     isCorrect,
		SubmittedAt:   time.Now(),
	}
	DB.Create(&submission)

	// Increment attempts
	progress.Attempts++
	DB.Save(&progress)

	if !isCorrect {
		c.JSON(http.StatusOK, gin.H{
			"correct": false,
			"message": "\x1b[1;31m[WRONG] Invalid flag hash. Review coordinates or inspect spelling.\x1b[0m",
		})
		return
	}

	// Check if already completed to avoid duplicate rewards
	wasComplete := progress.Status == "complete"
	var user User
	DB.First(&user, userID)

	xpEarned := 0
	eloDiff := 0
	oldElo := user.Elo

	if !wasComplete {
		progress.Status = "complete"
		now := time.Now()
		progress.CompletedAt = &now
		DB.Save(&progress)

		// Award rewards
		xpEarned = getXPReward(chapter.Difficulty)
		user.XP += xpEarned
		user.Level = 1 + (user.XP / 1000)

		challengeRating := getDifficultyRating(chapter.Difficulty)
		user.Elo = calculateEloChange(user.Elo, challengeRating, "win", progress.HintsUsed)
		eloDiff = user.Elo - oldElo

		user.LastActive = &now
		DB.Save(&user)

		// Unlock next chapter
		var nextChapter Chapter
		err = DB.Where("campaign_slug = ? AND number = ?", slug, num+1).First(&nextChapter).Error
		if err == nil {
			var nextProgress UserProgress
			err = DB.Where("user_id = ? AND chapter_id = ?", userID, nextChapter.ID).First(&nextProgress).Error
			if err == gorm.ErrRecordNotFound {
				nextProgress = UserProgress{
					UserID:    userID,
					ChapterID: nextChapter.ID,
					Status:    "active",
				}
				DB.Create(&nextProgress)
			}
		}

		// Send async badge notifications if user reaches ELO threshold
		go func() {
			if oldElo < 1200 && user.Elo >= 1200 {
				SendBadgeEmail(user.Email, user.Username, "Power User Rank")
			}
			if oldElo < 1600 && user.Elo >= 1600 {
				SendBadgeEmail(user.Email, user.Username, "System Administrator Rank")
			}
			if oldElo < 2000 && user.Elo >= 2000 {
				SendBadgeEmail(user.Email, user.Username, "Incident Response Engineer Rank")
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":   true,
		"xp_earned": xpEarned,
		"old_elo":   oldElo,
		"new_elo":   user.Elo,
		"elo_diff":  eloDiff,
		"message":   "\x1b[1;32m[CORRECT] Flag matches! Threat isolated. Node security restored.\x1b[0m",
	})
}

// StartSandbox validates progression and returns a signed Cloudflare R2 container signed URL
func StartSandbox(c *gin.Context) {
	var body struct {
		ChapterNumber int `json:"chapter_number"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chapter number is required"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Verify progression
	if body.ChapterNumber > 0 {
		var prevChapter Chapter
		err = DB.Where("campaign_slug = ? AND number = ?", "antariksha", body.ChapterNumber-1).First(&prevChapter).Error
		if err == nil {
			var prevProgress UserProgress
			err = DB.Where("user_id = ? AND chapter_id = ?", userID, prevChapter.ID).First(&prevProgress).Error
			if err != nil || prevProgress.Status != "complete" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Complete the previous chapter first to boot this sandbox."})
				return
			}
		}
	}

	// Generate S3 Pre-signed URL or fallback
	signedURL, err := getSignedR2URL(body.ChapterNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate signed download URL: %v", err)})
		return
	}

	// Record sandbox session in DB
	var chapter Chapter
	if DB.Where("campaign_slug = ? AND number = ?", "antariksha", body.ChapterNumber).First(&chapter).Error == nil {
		session := SandboxSession{
			ID:          uuid.New(),
			UserID:      userID,
			ChapterID:   chapter.ID,
			ContainerID: fmt.Sprintf("sandbox-ch%d", body.ChapterNumber),
			StartedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		}
		DB.Create(&session)
	}

	c.JSON(http.StatusOK, gin.H{
		"signed_url": signedURL,
		"expires_in": 3600, // 1 hour
	})
}

// GetLeaderboard returns the top 20 players by Elo rating
func GetLeaderboard(c *gin.Context) {
	var users []User
	err := DB.Order("elo desc, xp desc").Limit(20).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	type LeaderboardEntry struct {
		Username string `json:"username"`
		Elo      int    `json:"elo"`
		XP       int    `json:"xp"`
		Rank     string `json:"rank"`
	}

	var list []LeaderboardEntry
	for _, u := range users {
		rank := "Newcomer"
		if u.Elo >= 2400 {
			rank = "Wizard"
		} else if u.Elo >= 2000 {
			rank = "Engineer"
		} else if u.Elo >= 1600 {
			rank = "Sysadmin"
		} else if u.Elo >= 1200 {
			rank = "Power User"
		} else if u.Elo >= 800 {
			rank = "User"
		}

		list = append(list, LeaderboardEntry{
			Username: u.Username,
			Elo:      u.Elo,
			XP:       u.XP,
			Rank:     rank,
		})
	}

	c.JSON(http.StatusOK, list)
}

// GetUserProfile fetches user profile detail mapping (ELO, XP, badges, streaks)
func GetUserProfile(c *gin.Context) {
	username := c.Param("username")
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	// Calculate completed chapters count
	var completedCount int64
	DB.Model(&UserProgress{}).Where("user_id = ? AND status = ?", user.ID, "complete").Count(&completedCount)

	// Determine rank band
	rank := "Newcomer"
	if user.Elo >= 2400 {
		rank = "Wizard"
	} else if user.Elo >= 2000 {
		rank = "Engineer"
	} else if user.Elo >= 1600 {
		rank = "Sysadmin"
	} else if user.Elo >= 1200 {
		rank = "Power User"
	} else if user.Elo >= 800 {
		rank = "User"
	}

	c.JSON(http.StatusOK, gin.H{
		"username":     user.Username,
		"email":        user.Email,
		"elo":          user.Elo,
		"xp":           user.XP,
		"level":        user.Level,
		"streak":       user.Streak,
		"rank":         rank,
		"completed_ch": completedCount,
		"created_at":   user.CreatedAt,
	})
}

// GetUserProgress returns a list of completed chapter IDs/numbers for a user
func GetUserProgress(c *gin.Context) {
	username := c.Param("username")
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	var progress []UserProgress
	DB.Where("user_id = ?", user.ID).Find(&progress)

	type ProgressItem struct {
		ChapterID string `json:"chapter_id"`
		Status    string `json:"status"`
	}

	var list []ProgressItem
	for _, p := range progress {
		list = append(list, ProgressItem{
			ChapterID: p.ChapterID.String(),
			Status:    p.Status,
		})
	}

	c.JSON(http.StatusOK, list)
}

// getSignedR2URL generates a time-limited signed S3 URL for Cloudflare R2 images
func getSignedR2URL(chapterNum int) (string, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	key := fmt.Sprintf("ch%d.img", chapterNum)

	// If R2 details are not supplied or are placeholders, fallback to public link or mock URL
	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || accountID == "your-cloudflare-account-id" {
		if publicURL != "" && publicURL != "https://pub-xxxxxx.r2.dev" {
			return fmt.Sprintf("%s/%s", publicURL, key), nil
		}
		return "/images/ch_base.img", nil
	}

	// Create S3 pre-signer for R2 S3 compatibility endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
			SigningRegion: "auto",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return "", err
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

// GetWeeklyQuestDetail retrieves the currently active Weekly Quest
func GetWeeklyQuestDetail(c *gin.Context) {
	now := time.Now().In(IST)
	var quest WeeklyQuest
	err := DB.Where("starts_at <= ? AND ends_at >= ?", now, now).First(&quest).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active Weekly Quest at this time."})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	// Get or start user progress
	var progress WeeklyQuestProgress
	err = DB.Where("user_id = ? AND weekly_quest_id = ?", userID, quest.ID).First(&progress).Error
	if err != nil {
		progress = WeeklyQuestProgress{
			UserID:        userID,
			WeeklyQuestID: quest.ID,
			StartedAt:     time.Now(),
			Completed:     false,
		}
		DB.Create(&progress)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          quest.ID,
		"title":       quest.Title,
		"description": quest.Description,
		"starts_at":   quest.StartsAt,
		"ends_at":     quest.EndsAt,
		"completed":   progress.Completed,
		"attempts":    progress.Attempts,
		"hints_used":  progress.HintsUsed,
	})
}

// SubmitWeeklyQuestFlag validates the submitted flag for the Weekly Quest
func SubmitWeeklyQuestFlag(c *gin.Context) {
	var input struct {
		Flag string `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().In(IST)
	var quest WeeklyQuest
	err := DB.Where("starts_at <= ? AND ends_at >= ?", now, now).First(&quest).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active Weekly Quest to submit to."})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var progress WeeklyQuestProgress
	err = DB.Where("user_id = ? AND weekly_quest_id = ?", userID, quest.ID).First(&progress).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quest progress not initialized. Please retrieve quest details first."})
		return
	}

	if progress.Completed {
		c.JSON(http.StatusOK, gin.H{"message": "Quest already completed!", "correct": true})
		return
	}

	progress.Attempts++
	submittedHash := hashFlag(input.Flag)

	if submittedHash == quest.FlagHash {
		progress.Completed = true
		completedTime := time.Now()
		progress.CompletedAt = &completedTime
		progress.Score = CalculateQuestScore(progress.StartedAt, completedTime, progress.HintsUsed, progress.Attempts)

		// Award user ELO and XP
		var user User
		DB.First(&user, userID)
		user.XP += 200
		user.Elo += 25
		DB.Save(&user)

		DB.Save(&progress)

		BroadcastLeaderboardUpdate()

		c.JSON(http.StatusOK, gin.H{
			"message": "Correct! Weekly Quest completed successfully.",
			"correct": true,
			"score":   progress.Score,
		})
	} else {
		DB.Save(&progress)
		c.JSON(http.StatusOK, gin.H{
			"message": "Incorrect flag. Try again.",
			"correct": false,
		})
	}
}

// GetDailyChallengeDetail retrieves today's Daily Challenge
func GetDailyChallengeDetail(c *gin.Context) {
	now := time.Now().In(IST)
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)

	var challenge DailyChallenge
	err := DB.Where("date = ?", todayMidnight).First(&challenge).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Daily challenge not generated yet."})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var progress DailyChallengeProgress
	err = DB.Where("user_id = ? AND daily_challenge_id = ?", userID, challenge.ID).First(&progress).Error
	completed := false
	if err == nil {
		completed = progress.Completed
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          challenge.ID,
		"description": challenge.Description,
		"date":        challenge.Date,
		"completed":   completed,
	})
}

// SubmitDailyChallengeFlag validates the daily challenge flag
func SubmitDailyChallengeFlag(c *gin.Context) {
	var input struct {
		Flag string `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().In(IST)
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)

	var challenge DailyChallenge
	err := DB.Where("date = ?", todayMidnight).First(&challenge).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active Daily Challenge to submit to."})
		return
	}

	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var progress DailyChallengeProgress
	err = DB.Where("user_id = ? AND daily_challenge_id = ?", userID, challenge.ID).First(&progress).Error
	if err != nil {
		progress = DailyChallengeProgress{
			UserID:           userID,
			DailyChallengeID: challenge.ID,
			Completed:        false,
		}
		DB.Create(&progress)
	}

	if progress.Completed {
		c.JSON(http.StatusOK, gin.H{"message": "Daily challenge already completed!", "correct": true})
		return
	}

	submittedHash := hashFlag(input.Flag)
	if submittedHash == challenge.FlagHash {
		progress.Completed = true
		completedTime := time.Now()
		progress.CompletedAt = &completedTime
		DB.Save(&progress)

		// Award user XP & increment streak
		var user User
		DB.First(&user, userID)
		user.XP += 50
		user.Streak += 1
		DB.Save(&user)

		BroadcastLeaderboardUpdate()

		c.JSON(http.StatusOK, gin.H{
			"message": "Correct! Daily challenge completed. Streak incremented!",
			"correct": true,
			"streak":  user.Streak,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Incorrect flag. Try again.",
			"correct": false,
		})
	}
}

func getRatingBand(elo int) string {
	if elo >= 2400 {
		return "Wizard"
	} else if elo >= 2000 {
		return "Engineer"
	} else if elo >= 1600 {
		return "Sysadmin"
	} else if elo >= 1200 {
		return "Power User"
	} else if elo >= 800 {
		return "User"
	}
	return "Newcomer"
}

// GetUserBadgeSVG returns a dynamic vector SVG badge for the user
func GetUserBadgeSVG(c *gin.Context) {
	username := c.Param("username")
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.Header("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, `<svg xmlns="http://www.w3.org/2000/svg" width="220" height="20">
  <rect width="220" height="20" fill="#ff4444" rx="3"/>
  <text x="110" y="14" fill="#fff" font-family="sans-serif" font-size="11" text-anchor="middle">Operator Not Found</text>
</svg>`)
		return
	}

	rank := getRatingBand(user.Elo)
	var badgeColor string
	switch rank {
	case "Newcomer":
		badgeColor = "#808080"
	case "User":
		badgeColor = "#00bcd4"
	case "Power User":
		badgeColor = "#4caf50"
	case "Sysadmin":
		badgeColor = "#ff9800"
	case "Engineer":
		badgeColor = "#9c27b0"
	case "Wizard":
		badgeColor = "#e91e63"
	default:
		badgeColor = "#39ff14"
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="230" height="20">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="230" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <rect width="110" height="20" fill="#0d1117"/>
    <rect x="110" width="120" height="20" fill="%s"/>
    <rect width="230" height="20" fill="url(#b)"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="55" y="15" fill="#c9d1d9" fill-opacity=".3">LinuxQuest ELO</text>
    <text x="55" y="14" fill="#c9d1d9">LinuxQuest ELO</text>
    <text x="170" y="15" fill="#fff" fill-opacity=".3">%d (%s)</text>
    <text x="170" y="14" fill="#fff">%d (%s)</text>
  </g>
</svg>`, badgeColor, user.Elo, rank, user.Elo, rank)

	c.Header("Content-Type", "image/svg+xml")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.String(http.StatusOK, svg)
}

