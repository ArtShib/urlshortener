package model

import (
	"errors"
)

// URL
type URL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// URLArray список URL
type URLArray []URL

// RequestShortener
type RequestShortener struct {
	URL    string `json:"url"`
	UserID string
}

// ResponseShortener структура для ответа в json
type ResponseShortener struct {
	Result string `json:"result"`
}

// RequestShortenerBatchArray список RequestShortenerBatch
type RequestShortenerBatchArray []RequestShortenerBatch

// ResponseShortenerBatch структура для ответа в json
type RequestShortenerBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResponseShortenerBatchArray список ResponseShortenerBatch
type ResponseShortenerBatchArray []ResponseShortenerBatch

// ResponseShortenerBatch структура для ответа в json
type ResponseShortenerBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ErrURLConflict кастомная ошибка "URL already exists"
var ErrURLConflict = errors.New("URL already exists")

// URLUser структура для ответа в json
type URLUser struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// URLUserBatch список URLUser
type URLUserBatch []URLUser

type contextKey string

// contextKey
const (
	UserIDKey   contextKey = "userID"
	OriginalURL contextKey = "originalURL"
)

// URLUserRequest структура для запроса url по userid
type URLUserRequest struct {
	UUID   string
	UserID string
}

// URLUserRequestArray список URLUserRequest
type URLUserRequestArray []URLUserRequest

// DeleteRequest структура запроса на удаление
type DeleteRequest struct {
	UUIDs  []string `json:"uuids"`
	UserID string   `json:"user_id"`
}
