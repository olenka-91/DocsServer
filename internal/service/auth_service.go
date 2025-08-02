package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/olenka-91/DocsServer/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	salt       = "jfls,eifnk"
	signingKey = "hkdkodjjh"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

type AuthService struct {
	repo       repository.Authorization
	adminToken string
}

func NewAuthService(r repository.Authorization, adminToken string) *AuthService {
	return &AuthService{repo: r, adminToken: adminToken}
}

func (a *AuthService) CreateUser(login, password string) (string, error) {
	password = generatePasswordHash(password)
	return a.repo.CreateUser(login, password)
}

func (a *AuthService) ValidateAdminToken(adminToken string) bool {
	return adminToken == a.adminToken
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := a.repo.GetUser(username, generatePasswordHash(password))
	logrus.Debug("pass=", generatePasswordHash(password))
	if err != nil {
		return "", ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Login,
	})
	return token.SignedString([]byte(signingKey))

}

func (a *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return "", errors.New("token claims are not of type tokenClaims")
	}

	return claims.Login, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
