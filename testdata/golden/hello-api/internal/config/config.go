package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort int
}

const defaultHTTPPort = 8080

func Load() (Config, error) {
	port := defaultHTTPPort
	if rawPort := os.Getenv("HTTP_PORT"); rawPort != "" {
		parsedPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return Config{}, fmt.Errorf("invalid HTTP_PORT %q: must be an integer", rawPort)
		}
		if parsedPort <= 0 || parsedPort > 65535 {
			return Config{}, fmt.Errorf("invalid HTTP_PORT %q: must be between 1 and 65535", rawPort)
		}
		port = parsedPort
	}

	return Config{HTTPPort: port}, nil
}
