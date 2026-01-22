package scanner

import (
	"CTI-Dashboard/scraper/logger"
	"CTI-Dashboard/scraper/output"
	"CTI-Dashboard/scraper/severity"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	DB         *sql.DB
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

func Run(opts Options) error {
	scanner := NewScanner(opts.Client, opts.Writer, opts.Timeout, opts.Proxy)
	response, err := scanner.checkTorStatus()
	if err != nil {
		logger.Error("Failed to check Tor status: ", "error", err)
		return err
	} else {
		if response.IsTor {
			logger.Info("Connected to Tor network", "IP", response.IP)
		} else {
			logger.Error("Not connected to Tor network")
			return err
		}
	}
	for _, target := range opts.Targets {
		fmt.Printf("Scanning target: %s  (Name: %s)\n", target, opts.TargetName)
		for i := 0; i < opts.Retries; i++ {
			fmt.Printf("Scraping - Attempt %d/3\n", i+1)
			response, err := opts.Client.Get(target)
			if err != nil {
				logger.Error("Request failed", "error", err, "target", target, "attempt", i+1)
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
				if i == opts.Retries-1 {
					return err
				}
				continue
			}
			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					logger.Error("Failed to read response body", "error", err, "target", target)
					fmt.Println(body)
					response.Body.Close()
					if i == opts.Retries-1 {
						return err
					}
					continue
				}

				screenShot, err := scanner.CaptureScreenshot(target, opts)
				if err != nil {
					logger.Error("Screenshot capture failed", "error", err, "target", target)
					response.Body.Close()
					if i == opts.Retries-1 {
						return err
					}
					continue
				}
				response.Body.Close()
				paths, err := opts.Writer.WriteResult(target, body, screenShot)
				if err != nil {
					logger.Error("Failed to write result", "error", err, "target", target)
					if i == opts.Retries-1 {
						return err
					}
					continue
				}
				logger.Info("Successfully scraped target", "target", target)
				UpdateLastScan(target, opts.TargetName, paths, opts.DB, body)
				break
			}
			response.Body.Close()
			return fmt.Errorf("request failed with status: %s", response.Status)
		}
	}
	return err
}

func RunPost(opts Options) error {
	scanner := NewScanner(opts.Client, opts.Writer, opts.Timeout, opts.Proxy)
	response, err := scanner.checkTorStatus()
	if err != nil {
		logger.Error("Failed to check Tor status: ", "error", err)
		return err
	} else {
		if response.IsTor {
			logger.Info("Connected to Tor network", "IP", response.IP)
		} else {
			logger.Error("Not connected to Tor network")
			return errors.New("not connected to Tor network")
		}
	}
	for _, target := range opts.Targets {
		fmt.Printf("Scanning target: %s  (Name: %s)\n", target, opts.TargetName)
		for i := 0; i < opts.Retries; i++ {
			fmt.Printf("Scraping - Attempt %d/3\n", i+1)
			response, err := opts.Client.Get(target)
			if err != nil {
				logger.Error("Request failed", "error", err, "target", target, "attempt", i+1)
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
				if i == opts.Retries-1 {
					return err
				}
				continue
			}
			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					logger.Error("Failed to read response body", "error", err, "target", target)
					fmt.Println(body)
					response.Body.Close()
					if i == opts.Retries-1 {
						return err
					}
					continue
				}
				response.Body.Close()
				logger.Info("Successfully scraped target", "target", target)
				UpdateLastScanPost(target, opts.DB, body)

				postBody := strings.NewReader(string(body))
				err = severity.AssessSeverity(postBody, opts.DB, target)
				if err != nil {
					logger.Error("Failed to assess severity", "error", err, "target", target)
				}
				break
			}
			response.Body.Close()
			return fmt.Errorf("request failed with status: %s", response.Status)
		}
	}
	return err
}

func (s *Scanner) CaptureScreenshot(targetURL string, opts Options) ([]byte, error) {

	optsScr := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(opts.Proxy),
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

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), optsScr...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(25*time.Second),
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *Scanner) checkTorStatus() (*TorStatus, error) {
	resp, err := s.Client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var status TorStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

func UpdateLastScan(target string, name string, paths []string, db *sql.DB, body []byte) {
	if len(paths) < 2 {
		logger.Error("Could not update: paths slice must contain at least 2 elements (HTML and Screenshot)")
		return
	}
	engine, err := identify_engine(string(body))
	if err != nil {
		logger.Error("Could not identify engine", err)
	}

	query := `UPDATE forums SET last_scaned = ?, forum_html = ?, forum_screenshot = ?, forum_engine = ? WHERE forum_url = ?`
	statement, err := db.Prepare(query)
	if err != nil {
		logger.Error("Could not prepare the database statement", err)
		return
	}
	defer statement.Close()
	ts := time.Now().Format("2006-01-02 15:04:05")
	_, err = statement.Exec(ts, paths[0], paths[1], engine, target)
	if err != nil {
		logger.Error("Could not update forum in the database", err)
		return
	}
	logger.Info("Successfully updated the last scan", "name", name)
}

func UpdateLastScanPost(target string, db *sql.DB, body []byte) {
	statement, err := db.Prepare(`UPDATE posts SET content = ? WHERE thread_url = ?`)
	if err != nil {
		logger.Error("Could not prepare the database statement", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(body, target)
	if err != nil {
		logger.Error("Could not update forum in the database", err)
		return
	}
	logger.Info("Successfully updated the last scan", "URL", target)
}

func identify_engine(html_body string) (string, error) {
	engines := map[string]string{
		`id="XF"`:       "XenForo",
		"my_post_key":   "MyBB",
		"wp-content":    "WordPress",
		"machina":       "Machina",
		"milligram.css": "RansomEXX-Custom",
	}
	for signature, engine := range engines {
		if strings.Contains(html_body, signature) {
			return engine, nil
		}
	}
	return "Unknown", nil
}
