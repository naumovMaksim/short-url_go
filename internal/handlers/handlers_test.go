package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/naumovMaksim/short-url_go/internal/config"
	"github.com/naumovMaksim/short-url_go/internal/handlers"
	"github.com/naumovMaksim/short-url_go/internal/service"
	"github.com/naumovMaksim/short-url_go/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_AddHandler(t *testing.T) {
	conf := &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	service := service.NewService(store, conf)
	h := handlers.NewHandler(service)

	r := chi.NewRouter()
	r.Post("/", h.AddHandler)

	srv := httptest.NewServer(r)
	defer srv.Close()

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
				statusCode: http.StatusNotFound,
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
			req := resty.New().R()
			req.Method = http.MethodPost
			req.SetBody(tt.longUrl)
			req.URL = srv.URL + tt.request

			resp, err := req.Send()
			assert.NoError(t, err, "error while creating http method")

			assert.Equal(t, tt.want.statusCode, resp.StatusCode(), "Wrong response code")

			if tt.want.statusCode == http.StatusCreated {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"), "Wrong content type")
				assert.NotEmpty(t, string(resp.Body()), "Body is empty")
				assert.Contains(t, string(resp.Body()), tt.want.shortUrlPrefix, "Body dose not contain http://localhost:8080/")
			}
		})
	}
}

func TestHandler_GetHandler(t *testing.T) {
	conf := &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	store.Set("Fpfrew35gbniufmh", "https://www.google.com/")
	service := service.NewService(store, conf)
	h := handlers.NewHandler(service)

	r := chi.NewRouter()
	r.Get("/{id}", h.GetHandler)
	r.Get("/", h.GetHandler)

	srv := httptest.NewServer(r)
	defer srv.Close()

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
			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			req := client.R()
			req.Method = http.MethodGet
			req.URL = srv.URL + tt.request

			resp, _ := req.Send()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.statusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.longUrl, resp.Header().Get("Location"))
			}
		})
	}
}
