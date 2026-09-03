package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
	"database/sql"
	"time"
)

type SpotifyTokenRepositoryImpl struct {
	*EntityRepository[*entities.SpotifyToken]
}

func NewSpotifyTokenRepository(db *sql.DB) *SpotifyTokenRepositoryImpl {
	return &SpotifyTokenRepositoryImpl{
		EntityRepository: NewEntityRepository[*entities.SpotifyToken](db),
	}
}

func (r *SpotifyTokenRepositoryImpl) CreateFromDTO(request *dtos.CreateSpotifyTokenDTO) (*entities.SpotifyToken, error) {
	return r.Create(&entities.SpotifyToken{
		AccessToken:  request.AccessToken,
		RefreshToken: request.RefreshToken,
		ExpiresAt:    request.ExpiresAt,
		UpdatedAt:    time.Now(),
	})
}

var _ SpotifyTokenRepositoryInterface = (*SpotifyTokenRepositoryImpl)(nil)
