package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
	"time"

	"github.com/jmoiron/sqlx"
)

type SpotifyTokenRepositoryImpl struct {
	*EntityRepositoryImpl[*entities.SpotifyToken]
}

func NewSpotifyTokenRepository(db *sqlx.DB) *SpotifyTokenRepositoryImpl {
	return &SpotifyTokenRepositoryImpl{
		EntityRepositoryImpl: NewEntityRepositoryImpl[*entities.SpotifyToken](db),
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
