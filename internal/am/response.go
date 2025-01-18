package am

import (
	"encoding/json"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

const (
	ErrorCodeInternalError = "INTERNAL_ERROR"
	ErrorCodeBadRequest    = "BAD_REQUEST"
	ErrorCodeNotFound      = "NOT_FOUND"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

func NewSuccessResponse(message string, data interface{}) Response {
	return Response{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, code string, details string) Response {
	return Response{
		Status:  StatusError,
		Message: message,
		Error: &APIError{
			Code:    code,
			Details: details,
		},
	}
}

func Respond(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
