package config

import "time"

type Config struct {
	// Input
	TargetFile string

	// Network
	TorProxy   string
	Timeout    time.Duration
	MaxRetries int

	// Output
	OutputDir  string
	ReportFile string
}
