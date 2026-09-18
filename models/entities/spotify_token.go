package entities

import (
	"time"
)

type SpotifyToken struct {
	GenericEntity
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (t *SpotifyToken) TableName() string {
	return "spotify_tokens"
}

var _ EntityInterface = (*SpotifyToken)(nil)
