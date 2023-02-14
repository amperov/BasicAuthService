package tools

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
}

func Error(err error) {
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type TokenClaims struct {
	*jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

func (t *TokenManager) GenerateToken(ctx context.Context, UserID int) (string, string) {
	if UserID == 0 {
		return "", ""
	}

	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAccess := jwt.NewNumericDate(time.Now().Add(1 * time.Minute))
	expiresRefresh := jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAccess,
		},
		UserID,
	})
	//Gen Access Token
	accessToken, err := accessClaims.SignedString([]byte(singingKey))
	Error(err)

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresRefresh,
		},
		UserID,
	})
	//Gen Refresh Token
	refreshToken, err := refreshClaims.SignedString([]byte(singingKey))
	Error(err)

	return accessToken, refreshToken
}

func (t *TokenManager) ValidateToken(ctx context.Context, token string) (int, error) {

	aToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}

		return []byte(singingKey), nil
	})

	claims, ok := aToken.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("invalid token") //TODO
	}

	return claims.UserId, err
}
