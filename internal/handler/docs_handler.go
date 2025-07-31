package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/handler/middleware"
	"github.com/olenka-91/DocsServer/internal/service"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getDocsList(ctx *gin.Context) (any, any, *middleware.ErrorResponse, int) {
	logrus.Debug("Entering getDocsList handler")

	currentUser := "user"

	input := entity.LimitedDocsListInput{
		Token: ctx.Query("token"), // Обязательный параметр
		Login: ctx.DefaultQuery("login", currentUser),
		Key:   ctx.Query("key"),
		Value: ctx.Query("value"),
	}

	logrus.Infof("Fetching documents list with filters: %+v", input)
	filteredDocs, err := h.services.Docs.GetDocsList(ctx, input)

	switch err {
	case nil:
		if ctx.Request.Method == http.MethodGet {
			return filteredDocs, nil, nil, http.StatusOK
		}
		return nil, nil, nil, http.StatusOK
	case service.ErrUnauthorized:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusUnauthorized, Text: err.Error()}, http.StatusUnauthorized
	case service.ErrForbidden:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusForbidden, Text: err.Error()}, http.StatusForbidden
	case service.ErrBadRequest:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusBadRequest, Text: err.Error()}, http.StatusBadRequest
	}
	return nil, nil, &middleware.ErrorResponse{Code: http.StatusInternalServerError, Text: err.Error()}, http.StatusInternalServerError

}

// func (h *Handler) getDoc(ctx *gin.Context) {
// 	logrus.Debug("Entering getDoc handler")

// 	token := ctx.Query("token")
// 	if token == "" {
// 		entity.SendErrorResponse(ctx, http.StatusBadRequest,
// 			entity.ErrMissingParameter, "Token is required")
// 		return
// 	}

// 	docId := ctx.Param("id")
// 	if docId == "" {
// 		entity.SendErrorResponse(ctx, http.StatusBadRequest,
// 			entity.ErrInvalidParameter, "Invalid docID value")
// 		return
// 	}

// 	if ctx.Request.Method == http.MethodHead {
// 		ctx.Status(http.StatusOK)
// 		return
// 	}

// 	logrus.Infof("Fetching document with id: %+v", docId)
// 	curDoc, err := h.services.Docs.GetDoc(docId)
// 	if err != nil {
// 		logrus.Errorf("Failed to fetch doc: %v", err)

// 		status := http.StatusInternalServerError
// 		errorCode := entity.ErrDatabase

// 		switch {
// 		case errors.Is(err, service.ErrInvalidToken):
// 			status = http.StatusUnauthorized
// 			errorCode = entity.ErrInvalidToken
// 		case errors.Is(err, service.ErrAccessDenied):
// 			status = http.StatusForbidden
// 			errorCode = entity.ErrAccessDenied
// 		case errors.Is(err, service.ErrInvalidInput):
// 			status = http.StatusBadRequest
// 			errorCode = entity.ErrInvalidParameter
// 		case errors.Is(err, sql.ErrNoRows):
// 			status = http.StatusNotFound
// 			errorCode = entity.ErrDocumentNotFound
// 		}

// 		entity.SendErrorResponse(ctx, status, errorCode, err.Error())
// 		return
// 	}

// 	logrus.Info("Successfully fetched 1 doc")

// 	if curDoc.File {
// 		ctx.Header("Content-Type", curDoc.Mime)
// 		ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", curDoc.Name))
// 		ctx.File(curDoc.Path)
// 	} else {
// 		entity.SendDataResponse(ctx, curDoc.JSONData)
// 	}
// }

func (h *Handler) getDoc(ctx *gin.Context) (any, any, *middleware.ErrorResponse, int) {
	logrus.Debug("Entering getDoc handler")

	id, _ := uuid.Parse(ctx.Param("id"))
	//login := ctx.GetHeader("X-User")
	login := ctx.Query("login")
	token := ctx.Query("token")

	doc, err := h.services.GetDoc(ctx, id, login, token)
	switch err {
	case nil:
		// if doc.File && ctx.Request.Method == http.MethodGet {
		// 	//ctx.Header("Content-Type", doc.Mime)
		// 	//ctx.File(doc.Path)
		// 	return doc, nil, nil, http.StatusOK
		// }
		return doc, nil, nil, http.StatusOK
	case service.ErrNotFound:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusNotFound, Text: err.Error()}, http.StatusNotFound
	case service.ErrForbidden:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusForbidden, Text: err.Error()}, http.StatusForbidden
	default:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusInternalServerError, Text: err.Error()}, http.StatusInternalServerError
	}
}
