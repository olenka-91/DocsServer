package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"

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

func (h *Handler) getDoc(ctx *gin.Context) (any, any, *middleware.ErrorResponse, int) {
	logrus.Debug("Entering getDoc handler")

	id, _ := uuid.Parse(ctx.Param("id"))
	//login := ctx.GetHeader("X-User")
	//	login := ctx.Query("login")
	token := ctx.Query("token")

	doc, err := h.services.GetDoc(ctx, id, token)
	switch err {
	case nil:
		if doc == nil { //сервис уже передал файл
			return nil, nil, nil, 0
		}
		return doc, nil, nil, http.StatusOK
	case os.ErrNotExist:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusNotFound, Text: err.Error()}, http.StatusNotFound

	case service.ErrNotFound:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusNotFound, Text: err.Error()}, http.StatusNotFound
	case service.ErrForbidden:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusForbidden, Text: err.Error()}, http.StatusForbidden
	default:
		return nil, nil, &middleware.ErrorResponse{Code: http.StatusInternalServerError, Text: err.Error()}, http.StatusInternalServerError
	}
}

type UploadResponse struct {
	JSON interface{} `json:"json"`
	File string      `json:"file"`
}

func (h *Handler) postDoc(ctx *gin.Context) (any, any, *middleware.ErrorResponse, int) {
	logrus.Debug("Entering uploadDocument handler")

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, nil, &middleware.ErrorResponse{
			Code: http.StatusBadRequest,
			Text: "failed to parse multipart form",
		}, http.StatusBadRequest
	}
	defer form.RemoveAll()

	// Обработка meta-части
	metaValues := form.Value["meta"]
	if len(metaValues) == 0 {
		return nil, nil, &middleware.ErrorResponse{
			Code: http.StatusBadRequest,
			Text: "missing meta field",
		}, http.StatusBadRequest
	}

	var meta entity.UploadMeta
	if err := json.Unmarshal([]byte(metaValues[0]), &meta); err != nil {
		return nil, nil, &middleware.ErrorResponse{
			Code: http.StatusBadRequest,
			Text: "invalid meta format",
		}, http.StatusBadRequest
	}

	var jsonData entity.JSONB
	if jsonValues, exists := form.Value["json"]; exists && len(jsonValues) > 0 {
		//if err := json.Unmarshal([]byte(jsonValues[0]), &jsonData); err != nil {
		if err := jsonData.Scan(jsonValues); err != nil {
			return nil, nil, &middleware.ErrorResponse{
				Code: http.StatusBadRequest,
				Text: "invalid json format",
			}, http.StatusBadRequest
		}
	}

	var fileHeader *multipart.FileHeader
	if fileHeaders, exists := form.File["file"]; exists && len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	// currentUser, err := getUserFromContext(ctx)
	// if err != nil {
	// 	return nil, nil, &middleware.ErrorResponse{
	// 		Code: http.StatusUnauthorized,
	// 		Text: "authentication required",
	// 	}, http.StatusUnauthorized
	// }

	login := ctx.Query("login")
	token := ctx.Query("token")

	// Вызываем сервисный слой
	_, err = h.services.PostDoc(
		//ctx.Request.Context(),
		ctx,
		login,
		token,
		meta,
		jsonData,
		fileHeader,
	)

	if err != nil {
		logrus.Errorf("Upload document error: %v", err)
		return nil, nil, &middleware.ErrorResponse{
			Code: http.StatusInternalServerError,
			Text: "failed to upload document",
		}, http.StatusInternalServerError
	}

	data := UploadResponse{}
	data.JSON = jsonData
	if meta.File {
		data.File = meta.Name
	}

	return data, nil, nil, http.StatusCreated

}
