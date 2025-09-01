package shortener

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateUUID() (string, error) {
	len := 8
	bytes := make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:len], nil	
}

func GenerateShortURL(url string, uuid string) string {
	return url + "/" + uuid
}
