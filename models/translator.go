package models

import "time"

type Translator struct {
	ID             int    `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	HashedPassword string
	Suspension     bool `json:"suspension,omitempty"`
	RefreshToken   string
	Dialects       []string `json:"dialects,omitempty"`
	CreatedAt      time.Time
}
