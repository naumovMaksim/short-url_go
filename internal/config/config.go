package config

import (
	"flag"
	"os"
)

const (
	baseServerAdress = "localhost:8080"
	baseUrl          = "http://localhost:8080"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func ParseFlags() *Config {
	conf := &Config{}

	flag.StringVar(&conf.ServerAddress, "a", baseServerAdress, "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", baseUrl, "base address for shortened URL")

	flag.Parse()

	serverAddressFromEnv := os.Getenv("SERVER_ADDRESS")
	baseURLFromEnv := os.Getenv("BASE_URL")

	if serverAddressFromEnv != "" {
		conf.ServerAddress = serverAddressFromEnv
	}

	if baseURLFromEnv != "" {
		conf.BaseURL = baseURLFromEnv
	}

	if conf.ServerAddress != baseServerAdress &&
		conf.BaseURL == baseUrl &&
		baseURLFromEnv == "" {
		conf.BaseURL = "http://" + conf.ServerAddress
	}

	return conf
}
