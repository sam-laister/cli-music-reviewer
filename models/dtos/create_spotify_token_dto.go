package dtos

import "time"

type CreateSpotifyTokenDTO struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
