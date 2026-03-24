package service

import (
	"math/rand/v2"
	"strings"

	"github.com/naumovMaksim/short-url_go/internal/config"
	"github.com/naumovMaksim/short-url_go/internal/storage"
)

const (
	letters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength = 16
)

type Service struct {
	store   *storage.MemoryStorage
	baseURL string
}

func NewService(s *storage.MemoryStorage, c *config.Config) *Service {
	return &Service{
		store:   s,
		baseURL: c.BaseURL,
	}
}

func (s *Service) CreateShortUrl(url string) string {
	var key string
	for {
		key = randomString()
		_, ok := s.store.Get(key)
		if !ok {
			break
		}
	}

	s.store.Set(key, url)
	return s.baseURL + "/" + key
}

func (s *Service) GetLongUrl(key string) (longUrl string, ok bool) {
	longUrl, ok = s.store.Get(key)
	return
}

func randomString() string {
	var key strings.Builder

	for i := 0; i < idLength; i++ {
		index := rand.IntN(len(letters))
		key.WriteByte(letters[index])
	}

	return key.String()
}
