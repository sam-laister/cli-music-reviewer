package services

import (
	"cli-music-reviewer/config"
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/repositories"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	redirectUri = "http://127.0.0.1:" + config.SpotifyCallbackPort + "/callback"
)

// ErrNoStoredToken indicates no Spotify token has ever been persisted, so the
// full browser Authorize/callback flow is required instead of a refresh.
var ErrNoStoredToken = errors.New("no stored spotify token")

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyHandlerImpl struct {
	browserService   BrowserService
	spotifyTokenRepo repositories.SpotifyTokenRepositoryInterface
	clientId         string
	clientSecret     string
}

func (s *SpotifyHandlerImpl) Authorize() error {
	const (
		scope        = "user-read-private user-read-email user-library-read"
		responseType = "code"
	)

	base, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		log.Print(err)
		return err
	}

	params := url.Values{}
	params.Add("response_type", responseType)
	params.Add("client_id", s.clientId)
	params.Add("scope", scope)
	params.Add("redirect_uri", redirectUri)
	base.RawQuery = params.Encode()

	if err := s.browserService.Open(base.String()); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s *SpotifyHandlerImpl) Token(code string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("redirect_uri", redirectUri)
	form.Set("grant_type", "authorization_code")

	token, err := s.exchangeToken(form)
	if err != nil {
		return nil, err
	}

	if err := s.saveToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *SpotifyHandlerImpl) EnsureAuthorized() error {
	stored, err := s.spotifyTokenRepo.GetLatestOrNull()
	if err != nil {
		return err
	}
	if stored == nil {
		return ErrNoStoredToken
	}

	if time.Now().Before(stored.ExpiresAt) {
		return nil
	}

	return s.refresh(stored.RefreshToken)
}

func (s *SpotifyHandlerImpl) refresh(refreshToken string) error {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	token, err := s.exchangeToken(form)
	if err != nil {
		return err
	}

	// Spotify does not always return a new refresh_token on refresh — keep
	// reusing the existing one when it doesn't.
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}

	return s.saveToken(token)
}

func (s *SpotifyHandlerImpl) exchangeToken(form url.Values) (*TokenResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.clientId + ":" + s.clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		log.Print(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("spotify token exchange failed: %s", resp.Status)
		log.Print(err)
		return nil, err
	}

	return &token, nil
}

func (s *SpotifyHandlerImpl) saveToken(token *TokenResponse) error {
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	if _, err := s.spotifyTokenRepo.CreateFromDTO(&dtos.CreateSpotifyTokenDTO{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    expiresAt,
	}); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s *SpotifyHandlerImpl) GetSavedAlbums(limit, offset int, market string) (*dtos.GetSavedAlbumsResponseDTO, error) {
	stored, err := s.spotifyTokenRepo.GetLatestOrNull()
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrNoStoredToken
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	if market != "" {
		params.Set("market", market)
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/albums?"+params.Encode(), nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+stored.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("spotify get saved albums failed: %s", resp.Status)
		log.Print(err)
		return nil, err
	}

	var result dtos.GetSavedAlbumsResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Print(err)
		return nil, err
	}

	return &result, nil
}

func NewSpotifyHandler(browserService BrowserService, spotifyTokenRepo repositories.SpotifyTokenRepositoryInterface, clientId, clientSecret string) *SpotifyHandlerImpl {
	return &SpotifyHandlerImpl{
		browserService:   browserService,
		spotifyTokenRepo: spotifyTokenRepo,
		clientId:         clientId,
		clientSecret:     clientSecret,
	}
}
