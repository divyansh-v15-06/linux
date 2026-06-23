package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var oauthConfig *oauth2.Config
var jwtSecretKey []byte

// Claims defines custom JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// InitOAuth configures the Google OAuth and JWT settings
func InitOAuth() {
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecretKey) == 0 {
		jwtSecretKey = []byte("linuxquest_fallback_secret_key_1337_secret")
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

// GenerateTokens creates an access token (15 mins) and refresh token (7 days)
func GenerateTokens(user *User) (string, string, error) {
	// Access Token
	accessClaims := Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := Claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

// ValidateToken parses and validates a JWT token string
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HandleGoogleLogin redirects to Google consent screen
func HandleGoogleLogin(c *gin.Context) {
	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth is not configured on the server."})
		return
	}
	// Generate random state string (in production, use session store state to prevent CSRF)
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback processes Google OAuth response and logs user in
func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to exchange code: %v", err)})
		return
	}

	// Fetch user details from google endpoint
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get userinfo: %v", err)})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode google profile data"})
		return
	}

	if googleUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email details missing from Google profile"})
		return
	}

	// Find or create user
	var user User
	err = DB.Where("google_id = ? OR email = ?", googleUser.ID, googleUser.Email).First(&user).Error
	isNewUser := false

	if err == gorm.ErrRecordNotFound {
		isNewUser = true
		username := strings.Split(googleUser.Email, "@")[0]
		// Clean up username to exclude special chars
		username = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
				return r
			}
			return -1
		}, username)

		// Resolve duplicates
		var checkUser User
		baseUsername := username
		count := 1
		for DB.Where("username = ?", username).First(&checkUser).Error != gorm.ErrRecordNotFound {
			username = fmt.Sprintf("%s%d", baseUsername, count)
			count++
		}

		user = User{
			ID:       uuid.New(),
			Email:    googleUser.Email,
			Username: username,
			GoogleID: &googleUser.ID,
			Elo:      800,
			XP:       0,
			Level:    1,
			Streak:   1,
		}

		now := time.Now()
		user.LastActive = &now

		if err := DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user profile: %v", err)})
			return
		}

		// Create chapter 0 progress as active for new user
		var ch0 Chapter
		if err := DB.Where("number = ?", 0).First(&ch0).Error; err == nil {
			progress := UserProgress{
				UserID:    user.ID,
				ChapterID: ch0.ID,
				Status:    "active",
			}
			DB.Create(&progress)
		}

		// Send async welcome email
		go SendWelcomeEmail(user.Email, user.Username)
	} else if err == nil {
		// Update last active
		now := time.Now()
		user.LastActive = &now
		DB.Save(&user)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database query error: %v", err)})
		return
	}

	// Generate access/refresh JWT tokens
	accessStr, refreshStr, err := GenerateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate authorization tokens: %v", err)})
		return
	}

	// Return popup script response to postMessage back to parent and self close
	c.Header("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>ISRO CIRT — Auth Successful</title>
</head>
<body>
    <div style="font-family: monospace; background: #0d1117; color: #c9d1d9; height: 100vh; display: flex; align-items: center; justify-content: center; flex-direction: column;">
        <h3>Authentication Successful!</h3>
        <p>Transferring credentials to secure shell...</p>
    </div>
    <script>
        const authData = {
            accessToken: %q,
            refreshToken: %q,
            username: %q,
            email: %q,
            isNewUser: %t
        };
        window.opener.postMessage({ type: 'AUTH_SUCCESS', data: authData }, '*');
        window.close();
    </script>
</body>
</html>`, accessStr, refreshStr, user.Username, user.Email, isNewUser)

	c.String(http.StatusOK, html)
}

// HandleTokenRefresh exchanges refresh token for new access token
func HandleTokenRefresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	claims, err := ValidateToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	var user User
	err = DB.Where("id = ?", claims.UserID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User profile not found"})
		return
	}

	// Generate new tokens
	accessStr, refreshStr, err := GenerateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authorization tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessStr,
		"refresh_token": refreshStr,
	})
}
