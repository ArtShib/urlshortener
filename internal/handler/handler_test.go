package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


type URLServiceTest struct {
	Shortenfunc func(url string) (string, error)
	GetIDfunc func(shortCode string) (string, error)
}


func(s *URLServiceTest) Shorten(url string) (string, error) {
	return s.Shortenfunc(url)
} 

func(s * URLServiceTest) GetID(shortCode string) (string, error) {
	return s.GetIDfunc(shortCode)
}

func TestUrlHandler_Shorten(t *testing.T) {
	urlServiceTest := &URLServiceTest{
		Shortenfunc: func(url string) (string, error) {
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
		GetIDfunc: func(url string) (string, error) {
			return "sedczfrH", nil
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
