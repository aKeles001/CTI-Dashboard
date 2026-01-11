package scanner

import (
	"CTI-Dashboard/scraper/logger"
	"CTI-Dashboard/scraper/output"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	_ "github.com/mattn/go-sqlite3"
)

type Scanner struct {
	Client  *http.Client
	Writer  *output.Writer
	Timeout time.Duration
	Proxy   string
}
type Options struct {
	Targets    []string
	Client     *http.Client
	Writer     *output.Writer
	Timeout    time.Duration
	Retries    int
	TargetName string
	Proxy      string
}
type TorStatus struct {
	IP       string
	IsTor    bool
	Response string
}

func NewScanner(client *http.Client, writer *output.Writer, timeout time.Duration, proxy string) *Scanner {
	return &Scanner{
		Client:  client,
		Writer:  writer,
		Timeout: timeout,
		Proxy:   proxy,
	}
}

func Run(opts Options) {
	scanner := NewScanner(opts.Client, opts.Writer, opts.Timeout, opts.Proxy)
	response, err := scanner.checkTorStatus()
	if err != nil {
		logger.Error("Failed to check Tor status: %v", err)
	} else {
		if response.IsTor {
			logger.Info("Connected to Tor network", "IP", response.IP)
		} else {
			logger.Error("Not connected to Tor network")
		}
	}
	for _, target := range opts.Targets {
		fmt.Printf("Scanning target: %s (Name: %s)\n", target, opts.TargetName)
		for i := 0; i < opts.Retries; i++ {
			fmt.Printf("Scraping - Attempt %d/3\n", i+1)
			response, err := opts.Client.Get(target)
			if err != nil {
				logger.Error("Request failed", "error", err, "target", target, "attempt", i+1)
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
				continue
			}
			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					logger.Error("Failed to read response body", "error", err, "target", target)
					fmt.Println(body)
					response.Body.Close()
					continue
				}

				screenShot, err := scanner.CaptureScreenshot(target)
				if err != nil {
					logger.Error("Screenshot capture failed", "error", err, "target", target)
					response.Body.Close()
					continue
				}
				response.Body.Close()
				paths, err := opts.Writer.WriteResult(target, body, screenShot)
				if err != nil {
					logger.Error("Failed to write result", "error", err, "target", target)
					continue
				}
				logger.Info("Successfully scraped target", "target", target)
				UpdateLastScan(target, opts.TargetName, paths)
				break
			}
		}
	}
}

func (s *Scanner) CaptureScreenshot(targetURL string) ([]byte, error) {
	if s.Timeout == 0 {
		s.Timeout = 90 * time.Second
	}
	proxyURL := s.Proxy
	if proxyURL == "" {
		proxyURL = "127.0.0.1:9050"
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://"+proxyURL),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.NoSandbox,
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),

		chromedp.WaitVisible(`body`, chromedp.ByQuery),

		chromedp.Sleep(25*time.Second),

		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}
	return buf, nil
}

func (s *Scanner) checkTorStatus() (*TorStatus, error) {
	resp, err := s.Client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return nil, fmt.Errorf("Tor check failed: %w", err)
	}
	defer resp.Body.Close()
	var status TorStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

func updateLatsScan123(target string, name string, paths []string) {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		logger.Error("Could not connect to the database %v", err)
	}
	defer db.Close()
	statement, err := db.Prepare(`UPDATE forums SET last_scaned = ?, html_path = ?, screenshot_path = ? WHERE forum_url = ?`)
	if err != nil {
		logger.Error("Could not prepare the database statement %v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(time.Now(), paths[0], paths[1], target)
	if err != nil {
		logger.Error("Could not update forum in the database %v", err)
	}
	logger.Info("Succesfuly updated the last sacan of %v", name)
}

func UpdateLastScan(target string, name string, paths []string) {
	// 1. Check slice length before accessing indices to prevent panic
	if len(paths) < 2 {
		logger.Error("Could not update: paths slice must contain at least 2 elements")
		return
	}

	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		logger.Error("Could not connect to the database %v", err)
		return // Must return if db connection fails
	}
	defer db.Close()

	// Fixed typo: "last_scaned" to "last_scanned" (ensure this matches your SQL schema)
	query := `UPDATE forums SET last_scaned = ?, forum_html = ?, forum_screenshot = ? WHERE forum_url = ?`

	statement, err := db.Prepare(query)
	if err != nil {
		logger.Error("Could not prepare the database statement %v", err)
		return
	}
	defer statement.Close()

	// 2. Execute with parameters
	_, err = statement.Exec(time.Now(), paths[0], paths[1], target)
	if err != nil {
		logger.Error("Could not update forum in the database %v", err)
		return
	}

	// Fixed typos in log message
	logger.Info("Successfully updated the last scan of %v", name)
}
