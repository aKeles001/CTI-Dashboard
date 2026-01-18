package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"CTI-Dashboard/models"
	"CTI-Dashboard/scraper/config"
	"CTI-Dashboard/scraper/extractor"
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
	db     *sql.DB
}

func NewApp(cfg config.Config, client *http.Client, writer *output.Writer, db *sql.DB) *App {
	return &App{
		cfg:    cfg,
		client: client,
		writer: writer,
		db:     db,
	}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Add forum
func (a *App) CreateForum(forumData models.Forum) (string, error) {
	if forumData.ForumName == "" || forumData.ForumURL == "" {
		return "Error: Forum name and URL cannot be empty.", errors.New("forum name and URL cannot be empty")
	}

	forum_id := uuid.New().String()
	statement, err := a.db.Prepare(`INSERT INTO forums (forum_id, forum_name, forum_url, forum_description, last_scaned) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		logger.Error("Could not prepare the database statement", "error", err)
		return "Error: Could not prepare the database statement", err
	}
	defer statement.Close()

	_, err = statement.Exec(forum_id, forumData.ForumName, forumData.ForumURL, forumData.ForumDescription, "NULL")
	if err != nil {
		logger.Error("Could not insert forum into the database", "error", err)
		return "Error: Could not insert forum into the database", err
	}
	logger.Info("Successfully added forum", "name", forumData.ForumName)
	return fmt.Sprintf("Successfully added forum: %s", forumData.ForumName), nil
}

// Get Forum
func (a *App) GetForums() []models.Forum {
	rows, err := a.db.Query("SELECT forum_id, forum_url, forum_description, forum_name, last_scaned FROM forums")
	if err != nil {
		logger.Error("Could not prepare the database statement", "error", err)
		return nil
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var f models.Forum
		err := rows.Scan(&f.ForumID, &f.ForumURL, &f.ForumDescription, &f.ForumName, &f.LastScaned)
		if err != nil {
			logger.Error("Could not scan the database rows", "error", err)
			continue
		}
		forums = append(forums, f)
	}
	return forums

}

// Singular Forum Scrape
func (a *App) SingularScrape(forum models.Forum) error {
	err := scanner.Run(scanner.Options{
		Targets:    []string{forum.ForumURL},
		Client:     a.client,
		Writer:     a.writer,
		Timeout:    a.cfg.Timeout,
		Retries:    a.cfg.MaxRetries,
		TargetName: forum.ForumName,
		Proxy:      a.cfg.TorProxy,
		DB:         a.db,
	})
	if err != nil {
		logger.Error("Could not scrape forum", "error", err)
		return err
	}
	logger.Info("Successfully scraped forum", "name", forum.ForumName)
	return nil
}

// Multiple Forum Scrape
func (a *App) MultipleScrape(forums []models.Forum) []models.Forum {
	var errforums []models.Forum
	for _, forum := range forums {
		err := scanner.Run(scanner.Options{
			Targets:    []string{forum.ForumURL},
			Client:     a.client,
			Writer:     a.writer,
			Timeout:    a.cfg.Timeout,
			Retries:    a.cfg.MaxRetries,
			TargetName: forum.ForumName,
			Proxy:      a.cfg.TorProxy,
			DB:         a.db,
		})
		if err != nil {
			logger.Error("Could not scrape forum", "error", err)
			errforums = append(errforums, forum)
			continue
		}
	}
	if errforums != nil {
		logger.Error("Could not scrape all forums", "error", errforums)
		return errforums
	}
	logger.Info("All forums scraped")
	return nil
}

// Delete Forum
func (a *App) DeleteForum(forumID string) error {
	statement, err := a.db.Prepare(`DELETE FROM forums WHERE forum_id = ?`)
	if err != nil {
		logger.Error("Could not prepare the database statement", "error", err)
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(forumID)
	if err != nil {
		logger.Error("Could not delete forum from the database", err)
		return err
	}
	logger.Info("Successfully deleted forum", "id", forumID)
	return err
}

func (a *App) Extract_posts(forum_id string) (int, error) {
	result, err := extractor.PostExtract(forum_id, a.db)
	return result, err
}

func (a *App) GetPosts(forumID string) ([]models.Post, error) {
	statement, err := a.db.Prepare(`SELECT post_id, forum_id, thread_url, author, content, date FROM posts WHERE forum_id = ?`)
	if err != nil {
		logger.Error("Could not prepare statement", "error", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(forumID)
	if err != nil {
		logger.Error("Could not query posts from the database", "error", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var threadUrl sql.NullString
		err := rows.Scan(&post.PostID, &post.ForumID, &threadUrl, &post.PostAuthor, &post.PostContent, &post.PostDate)
		if err != nil {
			logger.Error("Could not scan post row", "error", err)
			continue
		}
		if threadUrl.Valid {
			post.ThreadURL = threadUrl.String
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error during rows iteration", "error", err)
		return nil, err
	}

	return posts, nil
}

func (a *App) ScanPosts(forumID string) error {
	statement, err := a.db.Prepare(`SELECT post_id, thread_url FROM posts WHERE forum_id = ? and status = 'pending'`)
	if err != nil {
		logger.Error("Could not prepare statement", "error", err)
		return err
	}
	fmt.Print(statement)

	return nil
}
