package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/handler/middleware"
)

func (h *Handler) registerUser(c *gin.Context) (any, any, *middleware.ErrorResponse, int) {

	token := c.Query("token")
	if !h.services.ValidateAdminToken(token) {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusForbidden, Text: "have no rights"}, http.StatusForbidden
	}
	login := c.Query("login")
	pswd := c.Query("pswd")
	if err := validateLogin(login); err != nil {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusBadRequest, Text: "Invalid login: " + err.Error()}, http.StatusBadRequest
	}

	if err := validatePassword(pswd); err != nil {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusBadRequest, Text: "Invalid password: " + err.Error()}, http.StatusBadRequest
	}

	l, err := h.services.Authorization.CreateUser(login, pswd)
	if err != nil {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusInternalServerError, Text: err.Error()}, http.StatusInternalServerError
	}

	return nil, entity.RegisterUserResponse{Login: l}, nil, http.StatusOK
}

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) autorizeUser(c *gin.Context) (any, any, *middleware.ErrorResponse, int) {
	login := c.Query("login")
	pswd := c.Query("pswd")

	token, err := h.services.Authorization.GenerateToken(login, pswd)
	if err != nil {
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusInternalServerError, Text: err.Error()}, http.StatusInternalServerError
	}

	return nil, entity.AuthUserResponse{Token: token}, nil, http.StatusOK
}

func validateLogin(login string) error {
	if len(login) < 8 {
		return errors.New("less then 8 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(login) {
		return errors.New("forbidden symbol")
	}

	return nil
}

// Валидация пароля
func validatePassword(pswd string) error {
	if len(pswd) < 8 {
		return errors.New("less then 8 characters")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, ch := range pswd {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case !unicode.IsLetter(ch) && !unicode.IsDigit(ch):
			hasSpecial = true
		}
	}

	errors := []string{}
	if !hasUpper || !hasLower {
		errors = append(errors, "at least one uppercase and one lowercase letter")
	}
	if !hasDigit {
		errors = append(errors, "at least one digit")
	}
	if !hasSpecial {
		errors = append(errors, "at least one special character")
	}

	if len(errors) > 0 {
		return fmt.Errorf("requirements: %s", strings.Join(errors, ", "))
	}

	return nil
}

func (h *Handler) userIdentity(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		if form, err := c.MultipartForm(); err == nil && form != nil {
			if metaRaw := form.Value["meta"]; len(metaRaw) > 0 {
				var meta entity.UploadMeta
				if err := json.Unmarshal([]byte(metaRaw[0]), &meta); err == nil {
					token = meta.Token
				}
			}
		}
	}

	if token == "" {
		return
	}

	login, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		return
	}

	c.Set("login", login)
}

func getLogin(c *gin.Context) (string, error) {
	id, ok := c.Get("login")
	if !ok {
		return "", errors.New("login not found")
	}
	loginStr, ok := id.(string)
	if !ok {
		return "", errors.New("login is of invalid type")
	}
	return loginStr, nil

}
