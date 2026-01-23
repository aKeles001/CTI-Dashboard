package extractor

import (
	"CTI-Dashboard/scraper/logger"
	"errors"

	"database/sql"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type ExtractFunc func(doc *goquery.Document, db *sql.DB, forum_id string) (string, error)

var engines = map[string]ExtractFunc{
	"XenForo": extractor_XF,
}

func PostExtract(forum_id string, db *sql.DB) (int, error) {
	statement, err := db.Prepare(`SELECT forum_id, forum_engine, forum_name, forum_html, forum_url FROM forums WHERE forum_id = ?`)
	if err != nil {
		logger.Error("Could not prepare the database statement", "error", err)
		return 0, err
	}
	defer statement.Close()
	var engine string
	var forum_name string
	var forum_html string
	var forum_url string // Declared, but was not being scanned.
	err = statement.QueryRow(forum_id).Scan(&forum_id, &engine, &forum_name, &forum_html, &forum_url)
	if err != nil {
		logger.Error("Could not scan the database rows", "error", err)
		return 0, err
	}

	// Add a check for empty forum_html path
	if forum_html == "" {
		logger.Error("Forum HTML path is empty for forum", "forum_url", forum_url, "forum_id", forum_id)
		return 0, errors.New("forum HTML content not found, please scrape the forum first")
	}

	file, err := os.Open(forum_html)
	if err != nil {
		logger.Error("Could not open the file", "error", err)
		return 0, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		logger.Error("Could not parse the file", "error", err)
		return 0, err
	}

	if extractor, ok := engines[engine]; ok {
		var links []string
		links = ExtractThreadLinks(doc, forum_url) // forum_id is now correctly passed
		err = ProcessExtractedLinks(forum_id, links, db)
		if err != nil {
			logger.Error("Could not process extracted links", "error", err)
			return 0, err
		}

		extractor(doc, db, forum_id)
		return len(links), nil
	}
	logger.Error("Enigne type is not supported for", "Forum ID", forum_id)
	return 0, nil
}

func extractor_XF(doc *goquery.Document, db *sql.DB, forum_id string) (string, error) {
	var posts []string
	doc.Find("div.bbWrapper").Each(func(i int, s *goquery.Selection) {
		posts = append(posts, strings.TrimSpace(s.Text()))
		logger.Info("Post", strings.TrimSpace(s.Text()))
	})

	return "", nil
}

func ExtractThreadLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		logger.Error("Could not parse base URL", "url", baseURL, "error", err)
		return links
	}
	rootURL := parsedURL.Scheme + "://" + parsedURL.Host
	doc.Find("div.structItem-title a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "/threads/") {
			threadIndex := strings.Index(href, "/threads/")
			if threadIndex != -1 {
				correctedHref := href[threadIndex:]
				fullURL := rootURL + correctedHref
				links = append(links, fullURL)
			}
		}
	})

	return links
}

func ProcessExtractedLinks(forumID string, links []string, db *sql.DB) error {
	for _, link := range links {
		post_id := uuid.New().String()
		_, err := db.Exec(`
            INSERT OR IGNORE INTO posts (post_id, forum_id, thread_url, status, title, content, author, date) 
            VALUES (?, ?, ?, 'pending', '', '', '', '')`,
			post_id, forumID, link,
		)
		if err != nil {
			logger.Error("Could not save link", "url", link, "error", err)
			return err
		}
	}
	return nil
}
