package severity

import (
	"CTI-Dashboard/scraper/logger"
	"database/sql"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

type SeverityLevel string

const (
	High    SeverityLevel = "high"
	Medium  SeverityLevel = "medium"
	Low     SeverityLevel = "low"
	Unknown SeverityLevel = "unknown"
)

var keywordSets = map[SeverityLevel][]string{
	High: {
		"turkey",
		"turkish",
		"TR",
		"turknet",
		"turkcell",
		"kablonet",
		"hacking",
		"hijacking",
		"breaching",
		"cracking",
		"vulnerability",
		"exploit",
		"attack",
		"malware",
		"botnet",
		"source code",
	},
	Medium: {
		"phishing",
		"scam",
		"otp code",
		"qr code scanner",
	},
	Low: {
		"card",
		"cc",
		"pirated",
	},
}

func AssessSeverity(postBody io.Reader, db *sql.DB, thread_url string) error {
	doc, err := goquery.NewDocumentFromReader(postBody)
	if err != nil {
		return err
	}
	severity := doc.Find("div.bbWrapper").Each(func(i int, s *goquery.Selection) {
		content := strings.ToLower(s.Text())
		for severity, keywords := range keywordSets {
			for _, keyword := range keywords {
				if strings.Contains(content, keyword) {
					_, err := db.Exec(`UPDATE posts SET severity_level = ? WHERE thread_url = ?`, severity, thread_url)
					if err != nil {
						logger.Error("Could not insert severity level to the database", "error", err)
						return
					}
					return
				}
			}
		}
	})
	logger.Info("Severity: ", &severity)
	return nil
}
