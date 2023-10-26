package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type RefreshTokenUserData struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func NewRefreshTokenUserData(ID string, expiryHours int64) *RefreshTokenUserData {
	expiryAtTime := time.Now().Add(time.Hour * time.Duration(expiryHours))
	return &RefreshTokenUserData{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAtTime),
		},
	}
}
