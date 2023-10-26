package services

import (
	"goback/domain/entity"
	"goback/utils"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	CreateAccessAndRefreshTokens(user *entity.User) (string, string, error)
	ComparePassword(hashPassword string, password string) error
	HashPassword(password string) (string, error)
	ExtractIDFromAccessToken(requestToken string) (string, error)
	ExtractIDFromRefreshToken(requestToken string) (string, error)
}

type AuthService struct {
	accessTokenSecret       string
	refreshTokenSecret      string
	accessTokenExpiryHours  int64
	refreshTokenExpiryHours int64
}

func NewAuthService(accessTokenSecret string, refreshTokenSecret string, accessTokenExpiryHours int64, refreshTokenExpiryHours int64) *AuthService {
	return &AuthService{
		accessTokenSecret:       accessTokenSecret,
		refreshTokenSecret:      refreshTokenSecret,
		accessTokenExpiryHours:  accessTokenExpiryHours,
		refreshTokenExpiryHours: refreshTokenExpiryHours,
	}
}

func (s AuthService) CreateAccessAndRefreshTokens(user *entity.User) (string, string, error) {
	accessToken, err := utils.CreateAccessToken(user, s.accessTokenSecret, s.accessTokenExpiryHours)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := utils.CreateRefreshToken(user, s.refreshTokenSecret, s.refreshTokenExpiryHours)
	return accessToken, refreshToken, err
}

func (s AuthService) ComparePassword(hashPassword string, password string) error {
	return utils.ComparePassword(hashPassword, password)
}

func (s AuthService) HashPassword(password string) (string, error) {
	return utils.HashPassword(password, bcrypt.DefaultCost)
}

func (s AuthService) ExtractIDFromAccessToken(requestToken string) (string, error) {
	return utils.ExtractIDFromToken(requestToken, s.accessTokenSecret)
}

func (s AuthService) ExtractIDFromRefreshToken(requestToken string) (string, error) {
	return utils.ExtractIDFromToken(requestToken, s.refreshTokenSecret)
}
