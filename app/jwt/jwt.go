package jwt

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Sub string `json:"sub"`
	jwt.StandardClaims
}

type Service interface {
	GetIdFromToken(string) (string, error)
}

type Impl struct {
	PublicKey *rsa.PublicKey
}

func NewJwt() (*Impl, error) {
	publicPemPath := os.Getenv("JWT_PUBLIC_PATH")

	keyBytes, err := os.ReadFile(publicPemPath)
	if err != nil {
		return nil, err
	}

	PublicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}

	return &Impl{PublicKey: PublicKey}, nil
}

func (j *Impl) GetIdFromToken(token string) (string, error) {
	unwrapToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.PublicKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := unwrapToken.Claims.(*Claims)
	if !ok {
		return "", errors.New("id not found")
	}

	userId := claims.Sub

	return userId, nil
}
