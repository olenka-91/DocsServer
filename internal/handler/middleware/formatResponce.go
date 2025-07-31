package middleware

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Error    *ErrorResponse `json:"error,omitempty"`
	Response any            `json:"response,omitempty"`
	Data     any            `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Success(c *gin.Context, status int, data, resp any) {
	c.JSON(status, ApiResponse{Response: resp, Data: data})
}

func Fail(c *gin.Context, status, code int, msg string) {
	c.JSON(status, ApiResponse{Error: &ErrorResponse{Code: code, Text: msg}})
}

type HandlerFunc func(*gin.Context) (data any, rsp any, apiErr *ErrorResponse, httpStatus int)

func Wrap(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, rsp, apiErr, httpStatus := h(c)
		if apiErr != nil {
			Fail(c, httpStatus, apiErr.Code, apiErr.Text)
			return
		}
		Success(c, httpStatus, data, rsp)
	}
}
