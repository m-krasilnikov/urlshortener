package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func New() *Config {
	config := &Config{}

	flag.StringVar(
		&config.ServerAddress,
		"a",
		"localhost:8080",
		"HTTP server address",
	)

	flag.StringVar(
		&config.BaseURL,
		"b",
		"http://localhost:8080",
		"base address for shortened URLs",
	)

	flag.Parse()

	if env := os.Getenv("SERVER_ADDRESS"); env != "" {
		config.ServerAddress = env
	}

	if env := os.Getenv("BASE_URL"); env != "" {
		config.BaseURL = env
	}

	return config
}
