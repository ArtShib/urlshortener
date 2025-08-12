package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// type UrlTest struct {
// 	LongUrl string
// 	ShortUrl string
// 	ShortCode string
// }

type UrlServiceTest struct {
	Shortenfunc func(url string) (string, error)
	GetIDfunc func(shortCode string) (string, error)
}


func(s *UrlServiceTest) Shorten(url string) (string, error) {
	return s.Shortenfunc(url)
} 

func(s * UrlServiceTest) GetID(shortCode string) (string, error) {
	return s.GetIDfunc(shortCode)
}

func TestUrlHandler_Shorten(t *testing.T) {
	urlServiceTest := &UrlServiceTest{
		Shortenfunc: func(url string) (string, error) {
			return "sedczfrH", nil
		},
	}
	
	handler := NewUrlHandler(urlServiceTest)

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

func TestUrlHandler_GetID(t *testing.T) {
	urlServiceTest := &UrlServiceTest{
		GetIDfunc: func(url string) (string, error) {
			return "sedczfrH", nil
		},
	}
	
	handler := NewUrlHandler(urlServiceTest)

	req := httptest.NewRequest("GET", "/e9db20b2", nil)

	w := httptest.NewRecorder()
	handler.GetID(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Status %d expected, status %d received", http.StatusTemporaryRedirect, resp.StatusCode)
	}	
}
