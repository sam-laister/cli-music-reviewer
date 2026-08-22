package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type HttpHandlerImpl struct {
	spotifyHandler SpotifyHandler
	port           string
	server         *http.Server
	done           chan error
}

func (s *HttpHandlerImpl) Setup() error {
	s.done = make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.spotifyAuthCallback)
	s.server = &http.Server{Handler: mux}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return err
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.done <- err
		}
	}()

	return nil
}

func (s *HttpHandlerImpl) Wait() error {
	err := <-s.done

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.server.Shutdown(shutdownCtx)

	return err
}

func (s *HttpHandlerImpl) spotifyAuthCallback(w http.ResponseWriter, req *http.Request) {
	if authErr := req.URL.Query().Get("error"); authErr != "" {
		http.Error(w, "Spotify authorization failed: "+authErr, http.StatusBadRequest)
		s.done <- fmt.Errorf("spotify authorization denied: %s", authErr)
		return
	}

	code := req.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code parameter", http.StatusBadRequest)
		s.done <- errors.New("callback missing code parameter")
		return
	}

	if _, err := s.spotifyHandler.Token(code); err != nil {
		http.Error(w, "failed to exchange token", http.StatusInternalServerError)
		s.done <- err
		return
	}

	fmt.Fprint(w, "Spotify authorization complete — you can close this tab and return to the CLI.")
	s.done <- nil
}

func NewHttpHandler(handler SpotifyHandler, port string) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		spotifyHandler: handler,
		port:           port,
	}
}
