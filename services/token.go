package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(email, userKey string, isAdmin bool) (string, error) {
	var signingKey = []byte(os.Getenv("SECRET"))

	claims := jwt.MapClaims{
		"email":   email,
		"userKey": userKey,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		"iss":     "GreenLibrary",
		"has":     isAdmin,
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
	// remove "Bearer " from token
	tokenstring = tokenstring[7:]
	token, err := jwt.ParseWithClaims(tokenstring, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok && time.Until(time.Unix(int64(exp), 0)) < 10*time.Minute {
			email := claims["email"].(string)
			userKey := claims["userKey"].(string)
			isAdmin := claims["has"].(bool)
			return GenerateToken(email, userKey, isAdmin)
		} else {
			return "", fmt.Errorf("token not ready to be refreshed yet")
		}
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
