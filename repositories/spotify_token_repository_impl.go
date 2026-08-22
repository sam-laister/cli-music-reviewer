package repositories

import (
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/models/entities"
	"database/sql"
	"errors"
	"time"
)

type SpotifyTokenRepositoryImpl struct {
	db *sql.DB
}

func (r *SpotifyTokenRepositoryImpl) Create(request dtos.CreateSpotifyTokenDTO) (int64, error) {
	result, err := r.db.Exec("INSERT INTO spotify_tokens (access_token, refresh_token, expires_at, updated_at) VALUES (?, ?, ?, ?)", request.AccessToken, request.RefreshToken, request.ExpiresAt, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SpotifyTokenRepositoryImpl) GetLatestOrNull() (*entities.SpotifyToken, error) {
	var token entities.SpotifyToken

	row := r.db.QueryRow("SELECT id, access_token, refresh_token, expires_at, updated_at FROM spotify_tokens ORDER BY id DESC LIMIT 1")
	if err := row.Scan(&token.ID, &token.AccessToken, &token.RefreshToken, &token.ExpiresAt, &token.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func NewSpotifyTokenRepository(db *sql.DB) *SpotifyTokenRepositoryImpl {
	return &SpotifyTokenRepositoryImpl{
		db: db,
	}
}
