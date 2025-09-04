package response

import (
	"errors"
	"net/http"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

// Data is a generic type for the data field in the response.
type Data any

// Response represents the response structure for HTTP responses.
type Response[T Data] struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    Data           `json:"data,omitempty"`
}

type ErrorResponse struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

// Success returns status code 200 and success response with data.
func Success(data Data) (int, *Response[Data]) {
	return http.StatusOK, &Response[Data]{
		Success: true,
		Data:    data,
	}
}

// Error returns error status code and error.
func Error(err error) (int, *Response[Data]) {
	var e *pkgerror.Error
	if errors.As(err, &e) {
		return responseCode[e.Code], &Response[Data]{
			Error: &ErrorResponse{
				Name:    responseName[e.Code],
				Message: err.Error(),
			},
		}
	}
	return InternalServerError(err)
}

// BadRequest returns status code 400 and error response.
func BadRequest(err error) (int, *Response[Data]) {
	return http.StatusBadRequest, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:    "BadRequest",
			Message: err.Error(),
		},
	}
}

// Unauthorized returns status code 401 and error response.
func Unauthorized(err error) (int, *Response[Data]) {
	return http.StatusUnauthorized, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:    "Unauthorized",
			Message: err.Error(),
		},
	}
}

// Forbidden returns status code 403 and error response.
func Forbidden(err error) (int, *Response[Data]) {
	return http.StatusForbidden, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:    "Forbidden",
			Message: err.Error(),
		},
	}
}

// InternalServerError returns status code 500 and error response.
func InternalServerError(err error) (int, *Response[Data]) {
	return http.StatusInternalServerError, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:    "InternalServerError",
			Message: err.Error(),
		},
	}
}

// responseCode is a slice of integer HTTP status codes used for error response mapping.
var responseCode = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusConflict,
	http.StatusInternalServerError,
}

// responseName is a list of string representations for HTTP response status codes.
var responseName = []string{
	"Ok",
	"BadRequest",
	"Unauthorized",
	"Forbidden",
	"NotFound",
	"Conflict",
	"InternalServerError",
}
