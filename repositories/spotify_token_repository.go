package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
)

type SpotifyTokenRepository interface {
	Create(request dtos.CreateSpotifyTokenDTO) (int64, error)
	GetLatestOrNull() (*entities.SpotifyToken, error)
}
