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
func (a *AuthService) signData(data string) string {
	h := hmac.New(sha256.New, a.secret)
	h.Write([]byte(data))
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

func (a *AuthService) verifySignature(data string, sign string) bool {
	expected := a.signData(data)
	return hmac.Equal([]byte(sign), []byte(expected))
}

func (a *AuthService) CreateToken(userID string) string {
	sign := a.signData(userID)
	token := sign + userID
	return token
}

func (a *AuthService) ValidateToken(token string) bool {
	userID := token[:sha256.Size]
	signToken := token[sha256.Size:]
	sign := a.CreateToken(userID)

	return a.verifySignature(signToken, sign)
}

//func (a *AuthService) splitToken(token string) []string {
//	return strings.Split(token, ":")
//}

func (a AuthService) GetUserID(token string) string {
	//parts := a.splitToken(token)
	//return parts[0]
	userID := token[:sha256.Size]
	return userID
}
