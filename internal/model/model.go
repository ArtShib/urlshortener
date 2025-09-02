package model

type URL struct { 
	UUID string `json:"uuid"`
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type HTTPServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`	
	// Host string
	// Port string `env:"SERVER_ADDRESS"`
}

type ShortServiceConfig struct {
	BaseURL string `env:"BASE_URL"`
	// ShortURL string `env:"BASE_URL"`
}

type RepositoryConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

type RequestShortener struct{
	URL string `json:"url"`
}

type ResponseShortener struct {
	Result string `json:"result"`
} 

