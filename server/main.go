package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if present
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	// Set gin mode based on environment
	env := os.Getenv("ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database (connects, migrates, and seeds chapters)
	InitDB()

	// Initialize Google OAuth config
	InitOAuth()

	// Initialize competition systems
	StartCompetitionScheduler()
	InitWebSocket()

	r := gin.Default()

	// CORS middleware supporting cookies/authorization headers
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Range")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, HEAD")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Range, Accept-Ranges, Content-Length")
		c.Writer.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")


		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.Static("/images", "./images")

	// Public Routes
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.GET("/api/ws/leaderboard", ServeLeaderboardWS)
	r.GET("/api/users/:username/badge.svg", GetUserBadgeSVG)

	// Auth Routes
	r.GET("/api/auth/google", HandleGoogleLogin)
	r.GET("/api/auth/google/callback", HandleGoogleCallback)
	r.POST("/api/auth/refresh", HandleTokenRefresh)
	r.GET("/api/auth/bypass", func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			username = "divyanshxanshu"
		}
		var user User
		err := DB.Where("username = ?", username).First(&user).Error
		if err != nil {
			user = User{
				ID:       uuid.New(),
				Email:    username + "@example.com",
				Username: username,
				Elo:      800,
				XP:       0,
				Level:    1,
				Streak:   1,
			}
			DB.Create(&user)
		}
		accessStr, refreshStr, _ := GenerateTokens(&user)
		c.JSON(http.StatusOK, gin.H{
			"accessToken":  accessStr,
			"refreshToken": refreshStr,
			"username":     user.Username,
			"email":        user.Email,
			"isNewUser":    false,
		})
	})

	// Protected Routes (Required JWT Token)
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		// Campaigns & Progress API
		protected.GET("/campaigns", GetCampaigns)
		protected.GET("/campaigns/:slug", GetCampaignBySlug)
		protected.GET("/campaigns/:slug/chapters/:num", GetChapterDetail)
		protected.POST("/campaigns/:slug/chapters/:num/submit", SubmitFlag)

		// Sandbox
		protected.POST("/sandbox/start", SandboxRateLimitMiddleware(), StartSandbox)

		// Profiles & Leaderboard
		protected.GET("/users/:username/profile", GetUserProfile)
		protected.GET("/users/:username/progress", GetUserProgress)
		protected.GET("/leaderboard", GetLeaderboard)

		// Competition System
		protected.GET("/competition/weekly", GetWeeklyQuestDetail)
		protected.POST("/competition/weekly/submit", SubmitWeeklyQuestFlag)
		protected.GET("/competition/daily", GetDailyChallengeDetail)
		protected.POST("/competition/daily/submit", SubmitDailyChallengeFlag)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
