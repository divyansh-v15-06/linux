package main

import (
	"testing"
	"time"
)

func TestCalculateQuestScore(t *testing.T) {
	// Base time check: 100s duration, 0 hints, 0 wrong attempts -> Score = 100
	t1 := time.Now()
	t2 := t1.Add(100 * time.Second)
	score := CalculateQuestScore(t1, t2, 0, 0)
	if score != 100.0 {
		t.Errorf("Expected score 100.0, got %f", score)
	}

	// 100s duration, 2 hints, 3 wrong attempts
	// Score = 100 * (1 + 0.1*2) * (1 + 0.05*3) = 100 * 1.2 * 1.15 = 138.0
	score2 := CalculateQuestScore(t1, t2, 2, 3)
	expected := 100.0 * 1.2 * 1.15
	if score2 != expected {
		t.Errorf("Expected score %f, got %f", expected, score2)
	}
}

func TestSeedActiveCompetitions(t *testing.T) {
	// Initialize in-memory SQLite for isolated tests
	InitDB()

	// Seed competitions
	SeedActiveCompetitions()

	// Verify Daily Challenge seeded
	now := time.Now().In(IST)
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, IST)
	var daily DailyChallenge
	if err := DB.Where("date = ?", todayMidnight).First(&daily).Error; err != nil {
		t.Errorf("Failed to find seeded Daily Challenge: %v", err)
	}

	// Verify Weekly Quest seeded
	var weekly WeeklyQuest
	if err := DB.Where("starts_at <= ? AND ends_at >= ?", now, now).First(&weekly).Error; err != nil {
		t.Errorf("Failed to find seeded Weekly Quest: %v", err)
	}
}
