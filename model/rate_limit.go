package model

import "time"

const (
	RateLimitMaximum        = 60
	RateLimitExpireDuration = 60 * time.Second
)

// RateLimit RateLimit
type RateLimit struct {
	Count  int `json:"count" redis:"count"`
	Expire int `json:"expire" redis:"expire"`
}
