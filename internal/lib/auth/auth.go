package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type AuthService struct {
	secret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{secret: []byte(secret)}
}

func (a AuthService) GenerateUserID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (a *AuthService) signData(data string) []byte {
	h := hmac.New(sha256.New, a.secret)
	DataBytes, _ := hex.DecodeString(data)
	h.Write(DataBytes)
	sign := h.Sum(nil)
	return sign // hex.EncodeToString(sign)
}

func (a *AuthService) verifySignature(data string, sign []byte) bool {
	expected := a.signData(data)
	return hmac.Equal(sign, expected)
}

func (a *AuthService) CreateToken(userID string) string {
	UserIDBytes, _ := hex.DecodeString(userID)
	sign := a.signData(userID)
	token := append(UserIDBytes, sign...)
	return hex.EncodeToString(token)
}

func (a *AuthService) ValidateToken(token string) bool {
	tokenBytes, _ := hex.DecodeString(token)
	userID := tokenBytes[:16]
	sign := tokenBytes[16:]
	expected := hex.EncodeToString(userID)

	return a.verifySignature(expected, sign)
}

//func (a *AuthService) splitToken(token string) []string {
//	return strings.Split(token, ":")
//}

func (a AuthService) GetUserID(token string) string {
	tokenBytes, _ := hex.DecodeString(token)
	userID := tokenBytes[:16]

	return hex.EncodeToString(userID)
}
