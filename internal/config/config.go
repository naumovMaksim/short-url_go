package config

import (
	"flag"
	"fmt"
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
	conf := Config{}
	flag.StringVar(&conf.ServerAddress, "a", baseServerAdress, "http server adress")
	flag.StringVar(&conf.BaseURL, "b", baseUrl, "URL adress to return")

	flag.Parse()
	if conf.BaseURL != baseUrl || conf.ServerAddress != baseServerAdress {
		if conf.ServerAddress != baseServerAdress && conf.BaseURL == baseUrl {
			newBaseUrl := "http://" + conf.ServerAddress
			conf.BaseURL = newBaseUrl
		}
		fmt.Printf("Server started on %v, and will return %v", conf.ServerAddress, conf.BaseURL)
	}
	return &conf
}
