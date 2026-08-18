package main

import (
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/services"
	"cli-music-reviewer/views"
	"database/sql"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func SetupDatabase() (*sql.DB, error) {
	dbService := services.NewDatabaseService(os.Getenv("DATABASE_URI"))

	db, err := dbService.Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupRepositories(db *sql.DB) *repositories.AppRepositories {
	return &repositories.AppRepositories{
		EntryRowRepository: repositories.NewEntryRowRepository(db),
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := SetupDatabase()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	defer db.Close()

	repos := SetupRepositories(db)

	if _, err := tea.NewProgram(views.NewHomepage(repos)).Run(); err != nil {
		log.Fatal("Uh oh:", err)
	}
}
