package config

import "time"

// Config struct
type ClientConfig struct {
	Timeout time.Duration
}

func Get() *ClientConfig {
	return &ClientConfig{
		Timeout: 5 * time.Second,
	}
}
