package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"goback/domain/auth"
	"goback/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(hashPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func HashPassword(password string, bcryptCostSize int) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcryptCostSize,
	)
	return string(encryptedPassword), err

}
func CreateAccessToken(user *entity.User, secret string, expiryHours int64) (string, error) {
	claims := auth.NewAccessTokenUserData(*user.Username, user.ID.String(), expiryHours)
	return createJwtToken(claims, secret)
}

func createJwtToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return jwtString, err
}

func CreateRefreshToken(user *entity.User, secret string, expiryHours int64) (refreshToken string, err error) {
	claims := auth.NewRefreshTokenUserData(user.ID.String(), expiryHours)
	return createJwtToken(claims, secret)
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, success := token.Claims.(jwt.MapClaims)

	if !success && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	return claims["id"].(string), nil
}
