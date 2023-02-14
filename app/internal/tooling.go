package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) HashPassword(Password string) string {

	hash := sha256.New()
	_, err := hash.Write([]byte(Password))
	if err != nil {
		return ""
	}
	HashPassword := hex.EncodeToString(hash.Sum([]byte("")))

	return HashPassword
}

func (h *Hasher) GenerateAccessCode(AccessToken, RefreshToken string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(AccessToken + RefreshToken))
	if err != nil {
		return "", err
	}
	AccessCode := hex.EncodeToString(hash.Sum([]byte("")))
	return AccessCode[19:45], nil
}
