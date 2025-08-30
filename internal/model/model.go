package model

type URL struct{
	LongURL string
	ShortURL string
	ShortCode string
}


type HTTPServerConfig struct {
	Host string
	Port string
}

type ShortServiceConfig struct {
	ShortURL string
}
