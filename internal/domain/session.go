package domain

import "time"

type Session struct {
	UserID       ID
	RefreshToken string
	ExpiresAt    time.Time
}
