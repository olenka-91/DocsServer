package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
)

func (h *Handler) signUp(c *gin.Context) {

	var req entity.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	tokens, err := h.services.Authorization.SignUp(req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Couldnt create user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "User created successfully",
		Data:    tokens,
	})
	return
}

func (h *Handler) signIn(c *gin.Context) {
	var req entity.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	tokens, err := h.services.Authorization.SignIn(req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Authentication failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Login successfully",
		Data:    tokens,
	})
	return
}

func (h *Handler) logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid request",
			Error:   "Parameter userID not found",
		})
		return
	}

	err := h.services.Authorization.Logout(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Failed to logout",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Logout successfully",
	})
	return
}

func (h *Handler) refreshToken(c *gin.Context) {
	var req entity.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	tokens, err := h.services.Authorization.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Failed to refresh token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Tokens refreshed successfully",
		Data:    tokens,
	})
	return
}
