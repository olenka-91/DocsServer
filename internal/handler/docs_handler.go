package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/service"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getDocsList(ctx *gin.Context) {
	logrus.Debug("Entering getDocsList handler")

	if ctx.Request.Method == http.MethodHead {
		ctx.Status(http.StatusOK)
		return
	}

	currentUser := "user"

	input := entity.LimitedDocsListInput{
		Token: ctx.Query("token"), // Обязательный параметр
		Login: ctx.DefaultQuery("login", currentUser),
		Key:   ctx.Query("key"),
		Value: ctx.Query("value"),
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			entity.SendErrorResponse(ctx, http.StatusBadRequest,
				entity.ErrInvalidParameter, "Invalid limit value")
			return
		}
		input.Limit = limit
	}

	if input.Token == "" {
		entity.SendErrorResponse(ctx, http.StatusBadRequest,
			entity.ErrMissingParameter, "Token is required")
		return
	}

	logrus.Infof("Fetching documents list with filters: %+v", input)
	filteredDocs, err := h.services.Docs.GetDocsList(input)
	if err != nil {
		logrus.Errorf("Failed to fetch docs list: %v", err)

		// Определение типа ошибки для правильного HTTP-статуса
		status := http.StatusInternalServerError
		errorCode := entity.ErrDatabase

		switch {
		case errors.Is(err, service.ErrInvalidToken):
			status = http.StatusUnauthorized
			errorCode = entity.ErrInvalidToken
		case errors.Is(err, service.ErrAccessDenied):
			status = http.StatusForbidden
			errorCode = entity.ErrAccessDenied
		case errors.Is(err, service.ErrInvalidInput):
			status = http.StatusBadRequest
			errorCode = entity.ErrInvalidParameter
		default:
			status = http.StatusInternalServerError
			errorCode = entity.ErrInvalidParameter
		}

		entity.SendErrorResponse(ctx, status, errorCode, err.Error())
		return
	}

	logrus.Infof("Successfully fetched %d docs", len(filteredDocs))
	entity.SendDataResponse(ctx, entity.DocsData{Docs: filteredDocs})
}
