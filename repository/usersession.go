package repository

import "time"

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
