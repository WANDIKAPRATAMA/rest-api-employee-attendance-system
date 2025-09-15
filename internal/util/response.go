package utils

import (
	"github.com/gofiber/fiber/v2"
)

// StatusResponse untuk menentukan status response
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

// ErrorDetail untuk menyimpan detail error validasi
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIResponse adalah struktur response standar
type APIResponse struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	Payload    Payload `json:"payload"`
}

// Payload berisi data, pagination, dan errors
type Payload struct {
	Data       interface{}   `json:"data"`
	Pagination interface{}   `json:"pagination"`
	Errors     []ErrorDetail `json:"errors"`
}

// Pagination untuk informasi paginasi
type Pagination struct {
	HasNextPage bool `json:"has_next_page"`
	NextPage    *int `json:"next_page"`
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
}

// SuccessResponse untuk response sukses
func SuccessResponse(statusCode int, message string, data interface{}, pagination interface{}) APIResponse {
	return APIResponse{
		Status:     StatusSuccess,
		StatusCode: statusCode,
		Message:    message,
		Payload: Payload{
			Data:       data,
			Pagination: pagination,
			Errors:     []ErrorDetail{},
		},
	}
}

// ErrorResponse untuk response error
func ErrorResponse(statusCode int, message string, errors []ErrorDetail) APIResponse {
	status := StatusError
	if statusCode >= fiber.StatusInternalServerError {
		status = StatusFail
	}
	return APIResponse{
		Status:     status,
		StatusCode: statusCode,
		Message:    message,
		Payload: Payload{
			Data:       nil,
			Pagination: struct{}{},
			Errors:     errors,
		},
	}
}
