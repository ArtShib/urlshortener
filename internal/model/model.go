package model

import (
	"errors"
)


type URL struct { 
	UUID string `json:"uuid"`
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type HTTPServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`	
}

type ShortServiceConfig struct {
	BaseURL string `env:"BASE_URL"`
}

type RepositoryConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

type RequestShortener struct{
	URL string `json:"url"`
}

type ResponseShortener struct {
	Result string `json:"result"`
} 

type RequestShortenerBatchArray [] RequestShortenerBatch

type RequestShortenerBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type ResponseShortenerBatchArray []ResponseShortenerBatch

type ResponseShortenerBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

var ErrURLConflict = errors.New("URL already exists")
