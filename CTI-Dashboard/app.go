package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"CTI-Dashboard/models"
	"CTI-Dashboard/scraper/config"
	"CTI-Dashboard/scraper/logger"
	"CTI-Dashboard/scraper/output"
	"CTI-Dashboard/scraper/scanner"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// App struct
type App struct {
	ctx    context.Context
	cfg    config.Config
	client *http.Client
	writer *output.Writer
}

func NewApp(cfg config.Config, client *http.Client, writer *output.Writer) *App {
	return &App{
		cfg:    cfg,
		client: client,
		writer: writer,
	}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Add forum
func (a *App) CreateForum(forumData models.Forum) string {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		logger.Error("Could not connect to the database %v", err)
	}
	defer db.Close()
	forum_id := uuid.New().String()
	statement, err := db.Prepare(`INSERT INTO forums (forum_id, forum_name, forum_url, forum_description, last_scaned) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		logger.Error("Could not prepare the database statement %v", err)
		os.Exit(1)
	}
	defer statement.Close()

	_, err = statement.Exec(forum_id, forumData.ForumName, forumData.ForumURL, forumData.ForumDescription, "NULL")
	if err != nil {
		logger.Error("Could not insert forum into the database %v", err)
		return "Error: Could not insert forum into the database"
	}
	logger.Info("Successfully added forum: %s", forumData.ForumName)
	return fmt.Sprintf("Successfully added forum: %s", forumData.ForumName)
}

// Get Foru
func (a *App) GetForums() []models.Forum {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		logger.Error("Could not connect to the database %v", err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT forum_id, forum_url, forum_description, forum_name, last_scaned FROM forums")
	if err != nil {
		logger.Error("Could not prepare the database statement %v", err)
		os.Exit(1)
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var f models.Forum
		err := rows.Scan(&f.ForumID, &f.ForumURL, &f.ForumDescription, &f.ForumName, &f.LastScaned)
		if err != nil {
			logger.Error("Could not scan the database rows %v", err)
			continue
		}
		forums = append(forums, f)
	}

	return forums
}

func (a *App) SingularScrape(forum models.Forum) string {
	scanner.Run(scanner.Options{
		Targets:    []string{forum.ForumURL},
		Client:     a.client,
		Writer:     a.writer,
		Timeout:    a.cfg.Timeout,
		Retries:    a.cfg.MaxRetries,
		TargetName: forum.ForumName,
		Proxy:      a.cfg.TorProxy,
	})
	logger.Info("Successfully scraped forum: %s", forum.ForumName)

	return "Successfully scraped forum:"
}
