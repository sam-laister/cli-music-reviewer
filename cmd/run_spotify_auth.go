package main

import (
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/services"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const port = "8888"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbService := services.NewDatabaseService(os.Getenv("DATABASE_URI"))
	db, err := dbService.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	spotifyTokenRepo := repositories.NewSpotifyTokenRepository(db)
	browserService := services.NewBrowserService()
	spotifyHandler := services.NewSpotifyHandler(browserService, spotifyTokenRepo, os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_SECRET"))
	httpHandler := services.NewHttpHandler(spotifyHandler, port)

	if err := services.Connect(spotifyHandler, httpHandler); err != nil {
		log.Fatal(err)
	}

	log.Print("Spotify connected")
}
