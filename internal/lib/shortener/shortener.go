package shortener

import (
	"crypto/rand"
	"encoding/base64"
)

// Shortener структура сервиса Shortener
type Shortener struct{}

// NewShortener конструктор Shortener
func NewShortener() *Shortener {
	return &Shortener{}
}

// GenerateUUID генерация uuid
func (s Shortener) GenerateUUID() (string, error) {
	lenUUID := 8
	bytes := make([]byte, lenUUID)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:lenUUID], nil
}

// GenerateShortURL генерация ShortURL
func (s Shortener) GenerateShortURL(url string, uuid string) string {
	return url + "/" + uuid
}
