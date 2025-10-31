package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/service"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getDocsList(ctx *gin.Context) {
	logrus.Debug("Entering getDocsList handler")

	input := entity.LimitedDocsListInput{
		//	Login: ctx.DefaultQuery("login", login),
		Key:   ctx.Query("key"),
		Value: ctx.Query("value"),
	}

	limitInt, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limitInt = 0
	}
	input.Limit = limitInt

	logrus.Infof("Fetching documents list with filters: %+v", input)
	filteredDocs, err := h.services.Docs.GetDocsList(ctx, input)

	switch err {
	case nil:
		if ctx.Request.Method == http.MethodGet {
			ctx.JSON(http.StatusOK, entity.SuccessResponse{
				Message: "Docs fetched successfully",
				Data:    filteredDocs,
			})
			return
		} else if ctx.Request.Method == http.MethodHead {
			ctx.JSON(http.StatusOK, entity.SuccessResponse{
				Message: "Docs fetched successfully",
			})
			return
		}
	case service.ErrForbidden:
		{
			ctx.JSON(http.StatusForbidden, entity.ErrorResponse{
				Message: "Forbidden",
				Error:   err.Error(),
			})
			return
		}
	case service.ErrBadRequest:
		{
			ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Message: "Bad",
				Error:   err.Error(),
			})
			return
		}
	default:
		{
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Message: "Internal Server Error",
				Error:   err.Error(),
			})
			return
		}
	}

}

func (h *Handler) getDoc(ctx *gin.Context) {
	logrus.Debug("Entering getDoc handler")

	login, exists := ctx.Get("login")
	if !exists {
		{
			ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
				Message: "Unauthorized",
			})
			return
		}
	}
	id, _ := uuid.Parse(ctx.Param("id"))

	doc, err := h.services.GetDoc(ctx, id, login.(string))
	switch err {
	case nil:
		if doc == nil { //сервис уже передал файл
			return
		}
		{
			ctx.JSON(http.StatusOK, entity.SuccessResponse{
				Message: "Doc fetched successfully",
				Data:    doc,
			})
			return
		}
	case os.ErrNotExist:
		{
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Message: "Not found",
				Error:   err.Error(),
			})
			return
		}
	case service.ErrNotFound:
		{
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Message: "Not found",
				Error:   err.Error(),
			})
			return
		}
	case service.ErrForbidden:
		{
			ctx.JSON(http.StatusForbidden, entity.ErrorResponse{
				Message: "Acess Forbidden",
				Error:   err.Error(),
			})
			return
		}
	default:
		{
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Message: "Internal Server Error",
				Error:   err.Error(),
			})
			return
		}
	}
}

type UploadResponse struct {
	JSON interface{} `json:"json"`
	File string      `json:"file"`
}

func (h *Handler) postDoc(ctx *gin.Context) {
	logrus.Debug("Entering postDoc handler")

	login, exists := ctx.Get("login")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}
	logrus.Debug("MultipartForm")

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Failed to parse multipart form",
			Error:   err.Error(),
		})
		return
	}
	defer form.RemoveAll()
	logrus.Debug("metaValues")
	// Обработка meta-части
	metaValues := form.Value["meta"]
	if len(metaValues) == 0 {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Missing meta field",
			Error:   err.Error(),
		})
		return
	}

	logrus.Debug("Unmarshal meta handler")
	var meta entity.UploadMeta
	if err := json.Unmarshal([]byte(metaValues[0]), &meta); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid meta format",
			Error:   err.Error(),
		})
		return
	}

	logrus.Debug("jsonData handler")
	var jsonData entity.JSONB
	if jsonValues, exists := form.Value["json"]; exists && len(jsonValues) > 0 {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: "Invalid JSON format",
			Error:   err.Error(),
		})
		return
	}

	var fileHeader *multipart.FileHeader
	if fileHeaders, exists := form.File["file"]; exists && len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	// Вызываем сервисный слой
	_, err = h.services.PostDoc(
		ctx,
		login.(string),
		meta,
		jsonData,
		fileHeader,
	)

	if err != nil {
		logrus.Errorf("Upload document error: %v", err)

		ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: "Failed to Upload document",
			Error:   err.Error(),
		})
		return
	}

	data := UploadResponse{}
	data.JSON = jsonData
	if meta.File {
		data.File = meta.Name
	}

	ctx.JSON(http.StatusCreated, entity.SuccessResponse{
		Message: "Doc uploaded successfully",
	})
	return
}

func (h *Handler) deleteDoc(ctx *gin.Context) {
	logrus.Debug("Entering deleteDoc handler")

	docID, _ := uuid.Parse(ctx.Param("id"))
	login, exists := ctx.Get("login")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, entity.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	resp, err := h.services.DeleteDoc(ctx, docID, login.(string))
	switch err {
	case nil:
		{
			ctx.JSON(http.StatusOK, entity.SuccessResponse{
				Message: "Doc deleted successfully",
				Data:    resp,
			})
			return
		}
	case os.ErrNotExist:
		{
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Message: "Not found",
				Error:   err.Error(),
			})
			return
		}
	case service.ErrNotFound:
		{
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Message: "Not found",
				Error:   err.Error(),
			})
			return
		}
	case service.ErrForbidden:
		{
			ctx.JSON(http.StatusForbidden, entity.ErrorResponse{
				Message: "Forbidden",
				Error:   err.Error(),
			})
			return
		}
	default:
		{
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Message: "Internal Server Error",
				Error:   err.Error(),
			})
			return
		}
	}
}
