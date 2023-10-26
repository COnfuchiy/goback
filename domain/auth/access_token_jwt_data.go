package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type AccessTokenUserData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.RegisteredClaims
}

func NewAccessTokenUserData(name string, ID string, expiryHours int64) *AccessTokenUserData {
	expiryAtTime := time.Now().Add(time.Hour * time.Duration(expiryHours))
	return &AccessTokenUserData{
		Name: name,
		ID:   ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAtTime),
		},
	}
}
