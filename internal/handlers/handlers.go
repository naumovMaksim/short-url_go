package handlers

import (
	"io"
	"net/http"

	"github.com/naumovMaksim/short-url_go/short-url_go/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AddHandler(w, r)
	case http.MethodGet:
		h.GetHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Path not allowed", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(body) <= 0 {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	shortUrl := h.service.CreateShortUrl(string(body))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	if len(urlPath) < 2 {
		http.Error(w, "Path not allowed", http.StatusBadRequest)
		return
	}

	id := urlPath[1:]

	if len(id) < 1 {
		http.Error(w, "Wrong id", http.StatusBadRequest)
		return
	}

	longUrl, ok := h.service.GetLongUrl(id)

	if !ok {
		http.Error(w, "No such id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", longUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
