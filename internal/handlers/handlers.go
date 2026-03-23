package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naumovMaksim/short-url_go/internal/service"
)

const (
	idLength = 16
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) != 1 {
		http.Error(w, "Wrong path", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	shortUrl := h.service.CreateShortUrl(string(body))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if len(id) < idLength {
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
