package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"

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
	token := &entities.SpotifyToken{
		AccessToken:  request.AccessToken,
		RefreshToken: request.RefreshToken,
		ExpiresAt:    request.ExpiresAt,
	}
	token.MarkCreated()
	token.MarkUpdated()

	return r.Create(token)
}

var _ SpotifyTokenRepositoryInterface = (*SpotifyTokenRepositoryImpl)(nil)
