package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
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

	flag.StringVar(
		&config.FileStoragePath,
		"f",
		"storage.json",
		"file storage path",
	)

	flag.Parse()

	if env := os.Getenv("SERVER_ADDRESS"); env != "" {
		config.ServerAddress = env
	}

	if env := os.Getenv("BASE_URL"); env != "" {
		config.BaseURL = env
	}

	if env := os.Getenv("FILE_STORAGE_PATH"); env != "" {
		config.FileStoragePath = env
	}

	return config
}
