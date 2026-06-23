package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for unit testing
func setupTestDB(t *testing.T) {
	var err error
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open sqlite database: %v", err)
	}

	err = DB.AutoMigrate(
		&User{},
		&Chapter{},
		&UserProgress{},
		&Submission{},
		&SandboxSession{},
		&CommandHistory{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	SeedChapters()
}

func TestJWTGenerationAndValidation(t *testing.T) {
	InitOAuth() // Sets default fallback secret

	user := User{
		ID:       uuid.New(),
		Email:    "test@isro.gov.in",
		Username: "test_agent",
	}

	accessToken, refreshToken, err := GenerateTokens(&user)
	if err != nil {
		t.Fatalf("GenerateTokens failed: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatal("AccessToken or RefreshToken is empty")
	}

	// Validate access token
	claims, err := ValidateToken(accessToken)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.Email != user.Email || claims.Username != user.Username || claims.UserID != user.ID.String() {
		t.Errorf("Claims mismatch: got %+v, want %+v", claims, user)
	}
}

func TestGetCampaigns(t *testing.T) {
	setupTestDB(t)
	InitOAuth()

	user := User{
		ID:       uuid.New(),
		Email:    "test@isro.gov.in",
		Username: "test_agent",
		Elo:      800,
	}
	DB.Create(&user)

	// Create test token
	token, _, _ := GenerateTokens(&user)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/campaigns", AuthMiddleware(), GetCampaigns)

	req, _ := http.NewRequest("GET", "/api/campaigns", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	campaigns, ok := resp["campaigns"].([]interface{})
	if !ok || len(campaigns) == 0 {
		t.Fatal("Expected campaigns list in response")
	}

	campaign := campaigns[0].(map[string]interface{})
	if campaign["slug"] != "antariksha" {
		t.Errorf("Expected campaign slug 'antariksha', got %v", campaign["slug"])
	}
}

func TestSubmitFlag(t *testing.T) {
	setupTestDB(t)
	InitOAuth()

	user := User{
		ID:       uuid.New(),
		Email:    "agent_arjun@isro.gov.in",
		Username: "arjun",
		Elo:      800,
		XP:       0,
	}
	DB.Create(&user)

	// Fetch Chapter 0 ID
	var ch0 Chapter
	DB.Where("number = ?", 0).First(&ch0)

	// Set progress for Ch 0 as active
	progress := UserProgress{
		UserID:    user.ID,
		ChapterID: ch0.ID,
		Status:    "active",
	}
	DB.Create(&progress)

	token, _, _ := GenerateTokens(&user)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/campaigns/:slug/chapters/:num/submit", AuthMiddleware(), SubmitFlag)

	// 1. Submit incorrect flag
	bodyIncorrect := map[string]string{"flag": "WRONG_FLAG"}
	bodyBytes, _ := json.Marshal(bodyIncorrect)
	req, _ := http.NewRequest("POST", "/api/campaigns/antariksha/chapters/0/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}

	var respIncorrect map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &respIncorrect)
	if respIncorrect["correct"].(bool) != false {
		t.Error("Expected incorrect validation result")
	}

	// 2. Submit correct flag: ANTARIKSHA{TERMINAL_WAKES_2026}
	bodyCorrect := map[string]string{"flag": "ANTARIKSHA{TERMINAL_WAKES_2026}"}
	bodyBytes, _ = json.Marshal(bodyCorrect)
	req, _ = http.NewRequest("POST", "/api/campaigns/antariksha/chapters/0/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}

	var respCorrect map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &respCorrect)
	if respCorrect["correct"].(bool) != true {
		t.Errorf("Expected correct validation, response: %v", w.Body.String())
	}

	// Verify User has completed chapter, XP increased, ELO updated
	var updatedUser User
	DB.First(&updatedUser, user.ID)
	if updatedUser.XP == 0 {
		t.Error("Expected XP to increase after correct flag submission")
	}
	if updatedUser.Elo <= 800 {
		t.Errorf("Expected ELO to increase, got %d", updatedUser.Elo)
	}

	var updatedProgress UserProgress
	DB.Where("user_id = ? AND chapter_id = ?", user.ID, ch0.ID).First(&updatedProgress)
	if updatedProgress.Status != "complete" {
		t.Errorf("Expected status to be complete, got %s", updatedProgress.Status)
	}
}

func TestProgressionGating(t *testing.T) {
	setupTestDB(t)
	InitOAuth()

	user := User{
		ID:       uuid.New(),
		Email:    "gatekeeper@isro.gov.in",
		Username: "gatekeeper",
		Elo:      800,
	}
	DB.Create(&user)

	token, _, _ := GenerateTokens(&user)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/campaigns/:slug/chapters/:num", AuthMiddleware(), GetChapterDetail)

	// Try to fetch Chapter 1 (locked since Chapter 0 is not complete)
	req, _ := http.NewRequest("GET", "/api/campaigns/antariksha/chapters/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden for locked chapter, got %d", w.Code)
	}
}

func TestMockGmailSMTPFallback(t *testing.T) {
	os.Setenv("GMAIL_USER", "")
	os.Setenv("GMAIL_APP_PASSWORD", "")

	err := SendEmail("arjun@isro.gov.in", "Dossier Access", "Workstation initialized.")
	if err != nil {
		t.Errorf("SendEmail should succeed with mock console logger fallback: %v", err)
	}
}
