package httpserver

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/deleteurls"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/getid"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/getjsonbatch"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/ping"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/shorten"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/shortenjson"
	"github.com/ArtShib/urlshortener/internal/httpserver/handlers/shortenjsonbatch"
	customMiddleware "github.com/ArtShib/urlshortener/internal/httpserver/middleware"
	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService описывает интерфейс сокращения url
type URLService interface {
	Shorten(ctx context.Context, url string) (string, error)
	GetID(ctx context.Context, shortCode string) (*model.URL, error)
	ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error)
	Ping(ctx context.Context) error
	ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error)
	GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
}

// WorkerPoolDelete описывает интерфейс удаления url
type WorkerPoolDelete interface {
	AddRequest(req model.DeleteRequest)
}

// ServiceEvent описывает интерфейс сохранения аудита
type ServiceEvent interface {
	AddEventRecord(event *model.Event)
}

// NewRouter конструктор Router
func NewRouter(svc URLService, log *slog.Logger, auth *auth.Service, poolDel WorkerPoolDelete, eventSvc ServiceEvent) http.Handler {

	mux := chi.NewRouter()
	mux.Use(customMiddleware.Auth(auth, log))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(customMiddleware.New(log))
	mux.Use(customMiddleware.GzipMiddleware)

	mux.Route("/api/user", func(r chi.Router) {
		r.Get("/urls", getjsonbatch.New(log, svc))
		r.Delete("/urls", deleteurls.New(log, poolDel))
	})
	mux.Get("/ping", ping.New(log, svc))
	mux.Group(func(r chi.Router) {
		r.Use(customMiddleware.NewEvent(log, eventSvc))
		r.Post("/", shorten.New(log, svc))
		r.Post("/api/shorten", shortenjson.New(log, svc))
		r.Post("/api/shorten/batch", shortenjsonbatch.New(log, svc))
		r.Get("/{shortCode}", getid.New(log, svc))
	})

	return mux
}
