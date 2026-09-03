package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
)

type SpotifyTokenRepositoryInterface interface {
	EntityRepositoryInterface[*entities.SpotifyToken]
	CreateFromDTO(request *dtos.CreateSpotifyTokenDTO) (*entities.SpotifyToken, error)
}
