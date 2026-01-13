package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Service структура сервиса авторизации
type Service struct {
	secret []byte
}

// NewAuthService конструктор Service
func NewAuthService(secret string) *Service {
	return &Service{secret: []byte(secret)}
}

// GenerateUserID создание userID
func (a Service) GenerateUserID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a Service) signData(data string) []byte {
	h := hmac.New(sha256.New, a.secret)
	DataBytes, _ := hex.DecodeString(data)
	h.Write(DataBytes)
	sign := h.Sum(nil)
	return sign // hex.EncodeToString(sign)
}

func (a Service) verifySignature(data string, sign []byte) bool {
	expected := a.signData(data)
	return hmac.Equal(sign, expected)
}

// CreateToken создание токена
func (a Service) CreateToken(userID string) string {
	UserIDBytes, _ := hex.DecodeString(userID)
	sign := a.signData(userID)
	token := append(UserIDBytes, sign...)
	return hex.EncodeToString(token)
}

// ValidateToken валидация токена
func (a Service) ValidateToken(token string) bool {
	tokenBytes, _ := hex.DecodeString(token)
	userID := tokenBytes[:16]
	sign := tokenBytes[16:]
	expected := hex.EncodeToString(userID)

	return a.verifySignature(expected, sign)
}

// GetUserID получеие userID
func (a Service) GetUserID(token string) string {
	tokenBytes, _ := hex.DecodeString(token)
	userID := tokenBytes[:16]

	return hex.EncodeToString(userID)
}
