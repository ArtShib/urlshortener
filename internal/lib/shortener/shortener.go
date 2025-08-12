package shortener

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateShortCode(text string) string{
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])[:8]
}

func GenerateShortURL(url string, code string) string {
	return "http://" + url + "/" + code
}
