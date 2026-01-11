package main

import (
	"embed"
	"time"

	"CTI-Dashboard/scraper/config"
	"CTI-Dashboard/scraper/logger"
	"CTI-Dashboard/scraper/output"
	"CTI-Dashboard/scraper/proxy"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.Config{
		Timeout:    time.Duration(30) * time.Second,
		MaxRetries: 3,
		OutputDir:  "output/",
		TorProxy:   "127.0.0.1:9050",
		TargetFile: "targets.yaml",
	}

	if err := logger.Init(cfg.OutputDir); err != nil {
		println("Error initializing logger:", err.Error())
		return
	}
	defer logger.Close()

	client, err := proxy.TorClient(cfg)
	if err != nil {
		logger.Error("Error initializing Tor client: %v", err)
		return
	}

	writer, err := output.NewWriter(cfg.OutputDir)
	if err != nil {
		logger.Error("Error initializing writer: %v", err)
		return
	}

	app := NewApp(cfg, client, writer)

	err = wails.Run(&options.App{
		Title:  "CTI-Dashboard",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		logger.Error("Error starting application: %v", err)
	}
}
