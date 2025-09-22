package response

import (
	"errors"
	"net/http"
	"time"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

// Data is a generic type for the data field in the response.
type Data any

// Response represents the response structure for HTTP responses.
type Response[T Data] struct {
	Success    bool           `json:"success"`
	Error      *ErrorResponse `json:"errors,omitempty"`
	Data       Data           `json:"data,omitempty"`
	ServerTime int64          `json:"serverTime,omitempty"`
}

type ErrorResponse struct {
	Name  string `json:"name,omitempty"`
	Error error  `json:"error,omitempty"`
}

// Success returns status code 200 and success response with data.
func Success(data Data) (int, *Response[Data]) {
	return http.StatusOK, &Response[Data]{
		Success:    true,
		Data:       data,
		ServerTime: serverTime(),
	}
}

// Error returns error status code and error.
func Error(err error) (int, *Response[Data]) {
	var e *pkgerror.Error
	if errors.As(err, &e) {
		return responseCode[e.Code], &Response[Data]{
			Error: &ErrorResponse{
				Name:  responseName[e.Code],
				Error: err,
			},
			ServerTime: serverTime(),
		}
	}
	return InternalServerError(err)
}

// BadRequest returns status code 400 and error response.
func BadRequest(err error) (int, *Response[Data]) {
	return http.StatusBadRequest, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:  "BadRequest",
			Error: err,
		},
		ServerTime: serverTime(),
	}
}

// Unauthorized returns status code 401 and error response.
func Unauthorized(err error) (int, *Response[Data]) {
	return http.StatusUnauthorized, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:  "Unauthorized",
			Error: err,
		},
		ServerTime: serverTime(),
	}
}

// Forbidden returns status code 403 and error response.
func Forbidden(err error) (int, *Response[Data]) {
	return http.StatusForbidden, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:  "Forbidden",
			Error: err,
		},
		ServerTime: serverTime(),
	}
}

// InternalServerError returns status code 500 and error response.
func InternalServerError(err error) (int, *Response[Data]) {
	return http.StatusInternalServerError, &Response[Data]{
		Success: false,
		Error: &ErrorResponse{
			Name:  "InternalServerError",
			Error: err,
		},
		ServerTime: serverTime(),
	}
}

// serverTime returns the current server time in milliseconds since the Unix epoch.
func serverTime() int64 {
	return time.Now().UnixMilli()
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
