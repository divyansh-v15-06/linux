package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts requests to validate JWT token in Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authorization token"})
			c.Abort()
			return
		}

		// Inject details into context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

var (
	sandboxStarts     = make(map[string][]time.Time)
	sandboxStartsLock sync.Mutex
)

// SandboxRateLimitMiddleware limits sandbox starts to max 3 requests per minute per user
func SandboxRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.Next()
			return
		}

		sandboxStartsLock.Lock()
		now := time.Now()
		cutoff := now.Add(-1 * time.Minute)

		starts := sandboxStarts[userID]
		var active []time.Time
		for _, t := range starts {
			if t.After(cutoff) {
				active = append(active, t)
			}
		}

		if len(active) >= 3 {
			sandboxStartsLock.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Max 3 sandbox starts per minute."})
			c.Abort()
			return
		}

		active = append(active, now)
		sandboxStarts[userID] = active
		sandboxStartsLock.Unlock()

		c.Next()
	}
}
