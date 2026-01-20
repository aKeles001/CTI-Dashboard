package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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
	wg     sync.WaitGroup
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
	statement, err := a.db.Prepare(`SELECT p.post_id, p.thread_url, f.forum_name FROM posts p JOIN forums f ON p.forum_id = f.forum_id WHERE p.forum_id = ?`)
	if err != nil {
		logger.Error("Could not prepare statement", "error", err)
		return err
	}
	defer statement.Close()

	rows, err := statement.Query(forumID)
	if err != nil {
		logger.Error("Could not fetch the rows", "error", err)
		return err
	}
	defer rows.Close()

	var allJobs []models.Job
	for rows.Next() {
		var job models.Job
		var threadUrl sql.NullString
		err := rows.Scan(&job.JobID, &threadUrl, &job.ForumName)
		if err != nil {
			logger.Error("Could not scan job row", "error", err)
			continue
		}
		if threadUrl.Valid {
			job.ThreadURL = threadUrl.String
		}
		allJobs = append(allJobs, job)
	}
	if err = rows.Err(); err != nil {
		logger.Error("Error during rows iteration", "error", err)
		return err
	}

	batchSize := a.cfg.Workers
	if batchSize <= 0 {
		batchSize = 10
	}

	numJobs := len(allJobs)
	for i := 0; i < numJobs; i += batchSize {
		end := i + batchSize
		if end > numJobs {
			end = numJobs
		}
		batch := allJobs[i:end]
		logger.Info("Processing batch", "start", i, "end", end-1, "total_jobs", numJobs)

		var batchWg sync.WaitGroup
		batchWg.Add(len(batch))

		for _, job := range batch {
			go func(j models.Job) {
				defer batchWg.Done()

				logger.Info("Processing job", "JobID", j.JobID)
				err := scanner.RunPost(scanner.Options{
					Targets:    []string{j.ThreadURL},
					Client:     a.client,
					Writer:     a.writer,
					Timeout:    a.cfg.Timeout,
					Retries:    a.cfg.MaxRetries,
					TargetName: j.ForumName,
					Proxy:      a.cfg.TorProxy,
					DB:         a.db,
				})

				if err != nil {
					logger.Error("Job failed", "id", j.JobID, "error", err, "Post URL", j.ThreadURL)
					a.db.Exec("UPDATE posts SET status = 'failed' WHERE post_id = ?", j.JobID)
				} else {
					a.db.Exec("UPDATE posts SET status = 'scraped' WHERE post_id = ?", j.JobID)
				}
				time.Sleep(50 * time.Second)
			}(job)
		}

		batchWg.Wait()
		logger.Info("Batch finished.", "size", len(batch))
	}
	logger.Info("All posts are scanned")
	return nil
}
