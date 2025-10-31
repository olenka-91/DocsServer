package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/repository"
	"github.com/olenka-91/DocsServer/internal/utils"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(r repository.Authorization) *AuthService {
	return &AuthService{
		repo: r,
	}
}

func (a *AuthService) SignUp(name, password string) (map[string]string, error) {
	existingUser, err := a.repo.GetUserByLogin(strings.ToLower(name))
	if (err == nil) && (existingUser != nil) {
		return nil, fmt.Errorf("user with this login already exists")
	}

	hashedPassword, err := utils.HashPaasword(password)
	if err != nil {
		return nil, err
	}

	ID, err := a.repo.CreateUser(strings.ToLower(name), hashedPassword)
	if err != nil {
		return nil, err
	}

	return a.generateAndSaveTokens(ID, strings.ToLower(name))
}

func (a *AuthService) SignIn(name, password string) (map[string]string, error) {
	existingUser, err := a.repo.GetUserByLogin(strings.ToLower(name))
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, fmt.Errorf("user with this login doesnt exist")
	}

	err = utils.CheckPasswordHash(password, existingUser.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return a.generateAndSaveTokens(existingUser.ID, strings.ToLower(name))
}

func (a *AuthService) RefreshToken(refreshToken string) (map[string]string, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	user, err := a.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Token != refreshToken {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return a.generateAndSaveTokens(claims.UserID, claims.Login)
}

func (a *AuthService) Logout(userID uuid.UUID) error {
	return a.repo.UpdateUserToken(userID, "")
}

func (a *AuthService) generateAndSaveTokens(userID uuid.UUID, login string) (map[string]string, error) {
	// Генерируем access токен
	accessToken, err := utils.GenerateToken(userID, login, "access")
	if err != nil {
		return nil, err
	}

	// Генерируем refresh токен
	refreshToken, err := utils.GenerateToken(userID, login, "refresh")
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh токен в базе
	err = a.repo.UpdateUserToken(userID, refreshToken)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
