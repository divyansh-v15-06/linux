package main

import (
	"fmt"
	"log"
)

// SendEmail mock prints the emails to server console logs, preventing port blocks on Render.
func SendEmail(to, subject, body string) error {
	log.Printf("[MOCK-EMAIL] To: %s | Subject: %s\n----- Body -----\n%s\n----------------", to, subject, body)
	return nil
}

// SendWelcomeEmail dispatches the welcome greeting to new users
func SendWelcomeEmail(to, username string) {
	subject := "Welcome to Operation Antariksha — LinuxQuest"
	body := fmt.Sprintf(`Welcome Agent %s,

You have been successfully added to the ISRO Cyber Incident Response Unit. 

A rogue AI payload signed "SHIVA" has compromised satellite communication terminals across key centers in India. You are our lead operative on this breach.

Your secure command shell workstation has been initialized.

Instructions:
1. Boot your terminal.
2. Complete the initial Bootcamp (Chapter 0) to verify your local toolchain.
3. Advance through the cities to isolate and shut down the rogue payload.

This is a zero-cost, serverless-sandboxed CTF. Copy-paste is disabled to build your command-line memory.

Remember: Every command counts. Good luck.

— Director Mehra
ISRO Cyber Incident Response Unit`, username)

	_ = SendEmail(to, subject, body)
}

// SendBadgeEmail notifies the user when they earn a badge
func SendBadgeEmail(to, username, badgeName string) {
	subject := fmt.Sprintf("COGNITIVE BADGE UNLOCKED: %s", badgeName)
	body := fmt.Sprintf(`Agent %s,

Congratulations. You have earned a new badge for your service:

Badge: %s
Status: Issued & Saved to profile

Your updated statistics have been compiled and sent to the global leaderboard.

Keep up the outstanding work.

— Director Mehra
ISRO Cyber Incident Response Unit`, username, badgeName)

	_ = SendEmail(to, subject, body)
}

// SendStreakEmail notifies a user that their streak is at risk
func SendStreakEmail(to, username string) {
	subject := "WARNING: Daily Streak Expiry Alert"
	body := fmt.Sprintf(`Agent %s,

Our monitoring system shows zero activity on your workstation in the past 20 hours. 

Your active daily challenge streak is currently at risk of expiring at midnight.

Boot your terminal and complete today's validation quest to retain your ranking multiplier.

— Command Operations
ISRO Cyber Incident Response Unit`, username)

	_ = SendEmail(to, subject, body)
}
