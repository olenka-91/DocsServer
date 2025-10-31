package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
				Message: "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
				Message: "Authorization header format must be Bearer!",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
				Message: "Invalid token",
				Error:   err.Error(),
			})
			ctx.Abort()
			return
		}

		if claims.Type != "access" {
			ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
				Message: "Token type must be access",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("login", claims.Login)
		ctx.Next()

	}
}
