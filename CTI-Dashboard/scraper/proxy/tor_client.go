package proxy

import (
	"net/http"

	"CTI-Dashboard/scraper/config"

	"golang.org/x/net/proxy"
)

func TorClient(cfg config.Config) (*http.Client, error) {
	if cfg.TorProxy == "" {
		cfg.TorProxy = "127.0.0.1:9050"
	}
	dialer, err := proxy.SOCKS5("tcp", cfg.TorProxy, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}
	return client, nil
}
