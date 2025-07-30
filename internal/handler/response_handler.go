package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/olenka-91/DocsServer/internal/entity"
)

// Успешный ответ с данными
func SendDataResponse(w http.ResponseWriter, data interface{}) {
	sendResponse(w, http.StatusOK, nil, nil, data)
}

// Успешный ответ без данных (для подтверждения действий)
func SendSuccessResponse(w http.ResponseWriter, response interface{}) {
	sendResponse(w, http.StatusOK, nil, response, nil)
}

// Отправка ошибки
func SendError(w http.ResponseWriter, status int, errorCode int, errorText string) {
	errResp := &entity.ErrorResponse{
		Code: errorCode,
		Text: errorText,
	}
	sendResponse(w, status, errResp, nil, nil)
}

// Внутренняя функция формирования ответа
func sendResponse(
	w http.ResponseWriter,
	status int,
	error *entity.ErrorResponse,
	response interface{},
	data interface{},
) {
	resp := entity.ApiResponse{
		Error:    error,
		Response: response,
		Data:     data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
