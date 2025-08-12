package model

type URL struct{
	LongUrl string
	ShortUrl string
	ShortCode string
}


type HttpServerConfig struct {
	Host string
	Port string
}

type ShortServiceConfig struct {
	ShortUrl string
}
