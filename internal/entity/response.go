package entity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Error    *ErrorResponse `json:"error,omitempty"`
	Response interface{}    `json:"response,omitempty"`
	Data     interface{}    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func SendDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Data: data,
	})
}

func SendSuccessResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Response: response,
	})
}

func SendErrorResponse(c *gin.Context, status int, errorCode int, errorText string) {
	c.JSON(status, ApiResponse{
		Error: &ErrorResponse{
			Code: errorCode,
			Text: errorText,
		},
	})
}
