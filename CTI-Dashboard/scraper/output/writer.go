package output

import (
	"os"
	"path/filepath"
)

type Writer struct {
	outputDir string
}

func NewWriter(outputDir string) (*Writer, error) {
	if outputDir == "" {
		outputDir = "output/"
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}
	return &Writer{outputDir: outputDir}, nil
}

func (w *Writer) WriteResult(target string, body []byte, screenshot []byte) (paths []string, err error) {
	safeTarget := filepath.Base(target)

	htmlDir := filepath.Join(w.outputDir, "html")
	htmlPath := filepath.Join(htmlDir, safeTarget+".html")
	if err := os.WriteFile(htmlPath, body, 0644); err != nil {
		return nil, err
	}
	screenshotDir := filepath.Join(w.outputDir, "screenshots")
	screenshotPath := filepath.Join(screenshotDir, safeTarget+".png")
	if err := os.WriteFile(screenshotPath, screenshot, 0644); err != nil {
		return nil, err
	}

	paths = []string{htmlPath, screenshotPath}
	return paths, nil
}
