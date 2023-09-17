package errors

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func NotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func BadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func UnexpectedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
