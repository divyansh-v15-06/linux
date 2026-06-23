package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB holds the global database instance
var DB *gorm.DB

// Models
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email      string    `gorm:"uniqueIndex;not null"`
	Username   string    `gorm:"uniqueIndex;not null"`
	GoogleID   *string   `gorm:"uniqueIndex"`
	Elo        int       `gorm:"default:800"`
	XP         int       `gorm:"default:0"`
	Level      int       `gorm:"default:1"`
	Streak     int       `gorm:"default:0"`
	LastActive *time.Time
	CreatedAt  time.Time
}

type Chapter struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	CampaignSlug string    `gorm:"index"`
	Number       int       `gorm:"index"`
	Title        string    `gorm:"not null"`
	City         string
	Difficulty   string
	Commands     string    `gorm:"type:text"` // Comma-separated or JSON
	FlagHash     string    `gorm:"not null"`  // SHA-256 hash of the flag
	StoryText    string    `gorm:"type:text"`
}

type UserProgress struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ChapterID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status      string    `gorm:"type:varchar(20);default:'locked'"` // 'locked', 'active', 'complete'
	Attempts    int       `gorm:"default:0"`
	HintsUsed   int       `gorm:"default:0"`
	CompletedAt *time.Time
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Chapter     Chapter   `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}

type Submission struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;index"`
	ChapterID     uuid.UUID `gorm:"type:uuid;index"`
	FlagSubmitted string    `gorm:"type:text"`
	IsCorrect     bool
	SubmittedAt   time.Time
	User          User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Chapter       Chapter   `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}

type SandboxSession struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	ChapterID   uuid.UUID `gorm:"type:uuid;index"`
	ContainerID string    `gorm:"type:text"`
	StartedAt   time.Time
	ExpiresAt   time.Time
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Chapter     Chapter   `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}

type CommandHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	ChapterID uuid.UUID `gorm:"type:uuid;index"`
	Command   string    `gorm:"type:text"`
	RanAt     time.Time
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Chapter   Chapter   `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}

type WeeklyQuest struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	FlagHash    string    `gorm:"not null"`
	StartsAt    time.Time `gorm:"index"`
	EndsAt      time.Time `gorm:"index"`
}

type WeeklyQuestProgress struct {
	UserID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	WeeklyQuestID uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Completed     bool        `gorm:"default:false"`
	StartedAt     time.Time
	CompletedAt   *time.Time
	HintsUsed     int         `gorm:"default:0"`
	Attempts      int         `gorm:"default:0"`
	Score         float64     `gorm:"default:0"`
	User          User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	WeeklyQuest   WeeklyQuest `gorm:"foreignKey:WeeklyQuestID;constraint:OnDelete:CASCADE"`
}

type DailyChallenge struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date        time.Time `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"type:text;not null"`
	FlagHash    string    `gorm:"not null"`
}

type DailyChallengeProgress struct {
	UserID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	DailyChallengeID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Completed        bool           `gorm:"default:false"`
	CompletedAt      *time.Time
	User             User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	DailyChallenge   DailyChallenge `gorm:"foreignKey:DailyChallengeID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hooks for platform-agnostic ID and Timestamp generation
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	return
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

func (s *Submission) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.SubmittedAt = time.Now()
	return
}

func (ss *SandboxSession) BeforeCreate(tx *gorm.DB) (err error) {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	ss.StartedAt = time.Now()
	return
}

func (ch *CommandHistory) BeforeCreate(tx *gorm.DB) (err error) {
	if ch.ID == uuid.Nil {
		ch.ID = uuid.New()
	}
	ch.RanAt = time.Now()
	return
}

func (wq *WeeklyQuest) BeforeCreate(tx *gorm.DB) (err error) {
	if wq.ID == uuid.Nil {
		wq.ID = uuid.New()
	}
	return
}

func (dc *DailyChallenge) BeforeCreate(tx *gorm.DB) (err error) {
	if dc.ID == uuid.Nil {
		dc.ID = uuid.New()
	}
	return
}

// InitDB initializes GORM connection, auto-migrates, and seeds data
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error

	if dsn == "" {
		log.Println("DATABASE_URL environment variable is not set. Falling back to local SQLite database: linuxquest.db")
		DB, err = gorm.Open(sqlite.Open("linuxquest.db"), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto migrate schema
	err = DB.AutoMigrate(
		&User{},
		&Chapter{},
		&UserProgress{},
		&Submission{},
		&SandboxSession{},
		&CommandHistory{},
		&WeeklyQuest{},
		&WeeklyQuestProgress{},
		&DailyChallenge{},
		&DailyChallengeProgress{},
	)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed.")

	// Seed chapters
	SeedChapters()
}

// hashFlag helper to compute SHA-256 hash of a string
func hashFlag(flag string) string {
	h := sha256.New()
	h.Write([]byte(flag))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SeedChapters seeds and synchronizes chapters from JSON files inside the campaigns folder.
// If the campaigns folder is empty or not found, it falls back to hardcoded seeding.
func SeedChapters() {
	campaignsDirs := []string{"campaigns", "../campaigns"}
	var jsonFiles []string

	for _, dir := range campaignsDirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.json"))
		if err == nil && len(files) > 0 {
			jsonFiles = files
			break
		}
	}

	if len(jsonFiles) > 0 {
		log.Printf("Found %d campaign JSON file(s) for auto-sync.", len(jsonFiles))
		for _, file := range jsonFiles {
			baseName := filepath.Base(file)
			slug := strings.TrimSuffix(baseName, filepath.Ext(baseName))

			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Failed to read campaign file %s: %v", file, err)
				continue
			}

			var jsonChapters []struct {
				Number     int    `json:"number"`
				Title      string `json:"title"`
				City       string `json:"city"`
				Difficulty string `json:"difficulty"`
				Commands   string `json:"commands"`
				Flag       string `json:"flag"`
				StoryText  string `json:"story_text"`
			}

			if err := json.Unmarshal(data, &jsonChapters); err != nil {
				log.Printf("Failed to unmarshal campaign JSON %s: %v", file, err)
				continue
			}

			for _, jc := range jsonChapters {
				var existing Chapter
				err := DB.Where("campaign_slug = ? AND number = ?", slug, jc.Number).First(&existing).Error
				if err == nil {
					// Update existing
					existing.Title = jc.Title
					existing.City = jc.City
					existing.Difficulty = jc.Difficulty
					existing.Commands = jc.Commands
					existing.FlagHash = hashFlag(jc.Flag)
					existing.StoryText = jc.StoryText
					if err := DB.Save(&existing).Error; err != nil {
						log.Printf("Failed to update chapter %d for campaign %s: %v", jc.Number, slug, err)
					}
				} else {
					// Create new
					newCh := Chapter{
						ID:           uuid.New(),
						CampaignSlug: slug,
						Number:       jc.Number,
						Title:        jc.Title,
						City:         jc.City,
						Difficulty:   jc.Difficulty,
						Commands:     jc.Commands,
						FlagHash:     hashFlag(jc.Flag),
						StoryText:    jc.StoryText,
					}
					if err := DB.Create(&newCh).Error; err != nil {
						log.Printf("Failed to insert chapter %d for campaign %s: %v", jc.Number, slug, err)
					}
				}
			}
			log.Printf("Campaign '%s' synchronized successfully (%d chapters).", slug, len(jsonChapters))
		}
		return
	}

	var count int64
	DB.Model(&Chapter{}).Count(&count)
	if count > 0 {
		log.Println("Chapters table already seeded. Skipping...")
		return
	}

	log.Println("Seeding chapters...")

	chapters := []Chapter{
		{
			CampaignSlug: "antariksha",
			Number:       0,
			Title:        "Bootcamp",
			City:         "Bangalore",
			Difficulty:   "Newcomer",
			Commands:     "echo,pwd,ls,cd,touch,cat,|,>",
			FlagHash:     hashFlag("ANTARIKSHA{TERMINAL_WAKES_2026}"),
			StoryText:    `Arjun, you've been assigned to the Cyber IR unit. Find orientation notes hint_*.txt in the filesystem to retrieve the flag.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       1,
			Title:        "The Lab",
			City:         "Bangalore",
			Difficulty:   "Beginner",
			Commands:     "ls,cd,cat,mkdir,touch",
			FlagHash:     hashFlag("ANTARIKSHA{shiva_init.sh}"),
			StoryText:    `Find the payload filename compromised by SHIVA. Payload is a SUID file dropped in home guest.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       2,
			Title:        "The Signal",
			City:         "Chennai",
			Difficulty:   "Beginner+",
			Commands:     "grep,sort,uniq,wc,|",
			FlagHash:     hashFlag("ANTARIKSHA{10.48.7.219}"),
			StoryText:    `Identify the destination IP of the SHIVA_PING packet stream in the log telemetry file.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       3,
			Title:        "The Hunt",
			City:         "Mumbai",
			Difficulty:   "Intermediate",
			Commands:     "ps,kill,pkill,top,lsof",
			FlagHash:     hashFlag("ANTARIKSHA{3847:kworker/u8}"),
			StoryText:    `A rogue processes is sending password credentials. Find the correct PID and name and kill it.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       4,
			Title:        "Cronjob of Doom",
			City:         "Delhi",
			Difficulty:   "Intermediate",
			Commands:     "crontab,systemctl,chmod +x,at",
			FlagHash:     hashFlag("ANTARIKSHA{*/15 * * * *}"),
			StoryText:    `SHIVA has planted a persistence backdoor in the cron jobs. Disable it and fetch schedule.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       5,
			Title:        "Permissions",
			City:         "Hyderabad",
			Difficulty:   "Intermediate-Advanced",
			Commands:     "chmod,chown,sudo,su,visudo",
			FlagHash:     hashFlag("ANTARIKSHA{cdac_stat:cdac_monitor}"),
			StoryText:    `A modified SUID binary allows root escalation. Find and modify permissions, lock compromised user.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       6,
			Title:        "The Archive",
			City:         "Pune",
			Difficulty:   "Advanced",
			Commands:     "tar,find,diff,sha256sum",
			FlagHash:     hashFlag("ANTARIKSHA{e9c0f83d7a8b5e2c}"),
			StoryText:    `Find the compressed file staging archive, verify its integrity and compare hashes with baseline.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       7,
			Title:        "Text Surgeon",
			City:         "Kolkata",
			Difficulty:   "Advanced",
			Commands:     "sed,awk,cut,jq",
			FlagHash:     hashFlag("ANTARIKSHA{12.9716,77.5946}"),
			StoryText:    `Extract characters from proxy logs and decode targeting coordinates.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       8,
			Title:        "The Shell Wars",
			City:         "Ahmedabad",
			Difficulty:   "Advanced",
			Commands:     "bash,ssh,trap,getopts",
			FlagHash:     hashFlag("ANTARIKSHA{11:2}"),
			StoryText:    `Write a robust shell script to clean cron entries across a multi-node cluster securely.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       9,
			Title:        "Ghost Signal",
			City:         "Chennai",
			Difficulty:   "Advanced-Pro",
			Commands:     "tcpdump,dig,nmap,nc",
			FlagHash:     hashFlag("ANTARIKSHA{10.0.0.4:10.0.0.15}"),
			StoryText:    `Listen on network interfaces and identify internal and external C2 IP addresses from raw DNS traffic.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       10,
			Title:        "SSH Tunnels",
			City:         "Remote",
			Difficulty:   "Pro",
			Commands:     "ssh,ssh-keygen,scp,rsync",
			FlagHash:     hashFlag("ANTARIKSHA{d3b07384d113edec}"),
			StoryText:    `Tunnel through three secure jump hosts using private keys to download memory dump database.`,
		},
		{
			CampaignSlug: "antariksha",
			Number:       11,
			Title:        "Final Shutdown",
			City:         "Sriharikota",
			Difficulty:   "Pro",
			Commands:     "systemctl,journalctl,vmstat,iostat,dmesg",
			FlagHash:     hashFlag("ANTARIKSHA{shiva_watchdog:f8a7e2b1}"),
			StoryText:    `Inspect systemd service dependencies, disable watchdog timers, and stop services in precise sequence.`,
		},
	}

	for _, ch := range chapters {
		ch.ID = uuid.New()
		if err := DB.Create(&ch).Error; err != nil {
			log.Fatalf("Failed to seed chapter %d: %v", ch.Number, err)
		}
	}
	log.Println("Chapter seeding complete.")
}
