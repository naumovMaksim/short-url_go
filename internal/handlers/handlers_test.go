package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/naumovMaksim/short-url_go/internal/handlers"
	"github.com/naumovMaksim/short-url_go/internal/service"
	"github.com/naumovMaksim/short-url_go/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_AddHandler(t *testing.T) {
	type want struct {
		contentType    string
		statusCode     int
		shortUrlPrefix string
	}
	tests := []struct {
		name    string
		request string
		longUrl string
		want    want
	}{
		{
			name:    "Good result test",
			request: "/",
			longUrl: "https://www.google.com/",
			want: want{
				contentType:    "text/plain",
				statusCode:     http.StatusCreated,
				shortUrlPrefix: "http://localhost:8080/",
			},
		},
		{
			name:    "Wrong path error test",
			request: "/test",
			longUrl: "https://www.google.com/",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Empty body error test",
			request: "/",
			longUrl: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemoryStorage()
			service := service.NewService(store)
			h := handlers.NewHandler(service)

			r := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.longUrl))
			w := httptest.NewRecorder()

			h.AddHandler(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code, "Wrong response code")

			if tt.want.statusCode == http.StatusCreated {
				assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"), "Wrong content type")
				assert.NotEmpty(t, w.Body.String(), "Body is empty")
				assert.Contains(t, w.Body.String(), tt.want.shortUrlPrefix, "Body dose not contain http://localhost:8080/")
			}
		})
	}
}

func TestHandler_GetHandler(t *testing.T) {
	type want struct {
		statusCode int
		longUrl    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "Good result test",
			request: "/Fpfrew35gbniufmh",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				longUrl:    "https://www.google.com/",
			},
		},
		{
			name:    "Empty id in url error test",
			request: "/",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Wrong id in url error test",
			request: "/123",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "No such id in storage error test",
			request: "/1234567891234567",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemoryStorage()
			store.Set("Fpfrew35gbniufmh", "https://www.google.com/")
			service := service.NewService(store)
			h := handlers.NewHandler(service)

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			h.GetHandler(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			if tt.want.statusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.longUrl, w.Header().Get("Location"))
			}
		})
	}
}
