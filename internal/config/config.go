package config

import (
	"flag"
	"os"
)

const (
	baseServerAdress = "localhost:8080"
	baseUrl          = "http://localhost:8080"
	baseLogLevel     = "info"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	LogLevel      string
}

func ParseFlags() *Config {
	conf := &Config{}

	flag.StringVar(&conf.ServerAddress, "a", baseServerAdress, "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", baseUrl, "base address for shortened URL")
	flag.StringVar(&conf.LogLevel, "l", baseLogLevel, "log level")

	flag.Parse()

	serverAddressFromEnv := os.Getenv("SERVER_ADDRESS")
	baseURLFromEnv := os.Getenv("BASE_URL")
	baseLogLevelFromEnv := os.Getenv("LOG_LEVEL")

	if serverAddressFromEnv != "" {
		conf.ServerAddress = serverAddressFromEnv
	}

	if baseURLFromEnv != "" {
		conf.BaseURL = baseURLFromEnv
	}

	if baseLogLevelFromEnv != "" {
		conf.LogLevel = baseLogLevelFromEnv
	}

	if conf.ServerAddress != baseServerAdress &&
		conf.BaseURL == baseUrl &&
		baseURLFromEnv == "" {
		conf.BaseURL = "http://" + conf.ServerAddress
	}

	return conf
}
