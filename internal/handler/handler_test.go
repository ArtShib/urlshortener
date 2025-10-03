package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
)

type URLServiceTest struct {
	Shortenfunc          func(ctx context.Context, url string) (string, error)
	GetIDfunc            func(ctx context.Context, shortCode string) (*model.URL, error)
	ShortenJsonfunc      func(ctx context.Context, url string) (*model.ResponseShortener, error)
	Pingfunc             func(ctx context.Context) error
	ShortenJSONBatchfunc func(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error)
	GetJSONBatchfunc     func(w http.ResponseWriter, r *http.Request) (model.URLUserBatch, error)
	DeleteBatchfunc      func(ctx context.Context, request *model.DeleteRequest) error
}

func (s *URLServiceTest) DeleteBatch(ctx context.Context, request *model.DeleteRequest) error {
	return s.DeleteBatchfunc(ctx, request)
}

func (s *URLServiceTest) Shorten(ctx context.Context, url string) (string, error) {
	return s.Shortenfunc(ctx, url)
}

func (s *URLServiceTest) GetID(ctx context.Context, shortCode string) (*model.URL, error) {
	return s.GetIDfunc(ctx, shortCode)
}

func (s *URLServiceTest) ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error) {
	return s.ShortenJsonfunc(ctx, url)
}

func (s *URLServiceTest) Ping(ctx context.Context) error {
	return nil
}

func (s *URLServiceTest) ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error) {
	return nil, nil
}

func (s *URLServiceTest) GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error) {
	return nil, nil
}

func TestUrlHandler_Shorten(t *testing.T) {
	urlServiceTest := &URLServiceTest{
		Shortenfunc: func(ctx context.Context, url string) (string, error) {
			return "sedczfrH", nil
		},
	}

	handler := NewURLHandler(urlServiceTest)

	bodyReq := strings.NewReader(`https://yandex.ru`)
	req := httptest.NewRequest("POST", "/", bodyReq)

	w := httptest.NewRecorder()
	handler.Shorten(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status %d expected, status %d received", http.StatusCreated, resp.StatusCode)
	}
}

func TestURLHandler_GetID(t *testing.T) {
	urlServiceTest := &URLServiceTest{
		GetIDfunc: func(ctx context.Context, url string) (*model.URL, error) {
			return &model.URL{
				ShortURL: "sedczfrH",
			}, nil
		},
	}

	handler := NewURLHandler(urlServiceTest)

	req := httptest.NewRequest("GET", "/e9db20b2", nil)

	w := httptest.NewRecorder()
	handler.GetID(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Status %d expected, status %d received", http.StatusTemporaryRedirect, resp.StatusCode)
	}
}
