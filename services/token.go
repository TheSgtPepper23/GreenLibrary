package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	email string
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {
	var signingKey = []byte(os.Getenv("SECRET"))

	claims := customClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Issuer:    "GreenLibrary",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RefreshToken(tokenstring string) (string, error) {
	var signingKey = []byte(os.Getenv("SECRET"))
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		if time.Until(claims.ExpiresAt.Time) < 10*time.Minute {
			return GenerateToken(claims.email)
		} else {
			return "", fmt.Errorf("token not ready to be refreshed yet")
		}
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
