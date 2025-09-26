package response

import (
	"errors"
	"net/http"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/codes"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/pkgerror"
)

// Response represents the response structure for HTTP responses.
type Response struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Data   any    `json:"data,omitempty"`
	Errors any    `json:"errors,omitempty"`
}

// Success returns status code 200 and success response with data.
func Success(data any) (int, Response) {
	return http.StatusOK, Response{
		Data: data,
	}
}

// Error returns error status code and error.
func Error(err error) (int, Response) {
	var e *pkgerror.Error
	if errors.As(err, &e) {
		switch e.Code {
		case codes.NotFound:
			return NotFound(err)
		case codes.Unauthenticated:
			return Unauthorized(err)
		case codes.Forbidden:
			return Forbidden(err)
		case codes.BadRequest:
			return BadRequest(err)
		case codes.Conflict:
			return Conflict(err)
		}
	}
	return InternalServerError(err)
}

// BadRequest returns status code 400 and error response.
func BadRequest(err error) (int, Response) {
	return http.StatusBadRequest, Response{
		Title:  "Invalid Request",
		Detail: "One or more parameters to your request was invalid.",
		Errors: err,
	}
}

// Unauthorized returns status code 401 and error response.
func Unauthorized(err error) (int, Response) {
	return http.StatusUnauthorized, Response{
		Title:  "Unauthorized",
		Detail: "You are not authorized to access this resource.",
		Errors: err,
	}
}

// Forbidden returns status code 403 and error response.
func Forbidden(err error) (int, Response) {
	return http.StatusForbidden, Response{
		Title:  "Forbidden",
		Detail: "You do not have permission to access this resource.",
		Errors: err,
	}
}

// Forbidden returns status code 403 and error response.
func NotFound(err error) (int, Response) {
	return http.StatusNotFound, Response{
		Title:  "Not Found",
		Detail: "The requested resource could not be found.",
		Errors: err,
	}
}

// Conflict returns status code 409 and error response.
func Conflict(err error) (int, Response) {
	return http.StatusConflict, Response{
		Title:  "Conflict",
		Detail: "The request could not be completed due to a conflict with the current state of the resource.",
		Errors: err,
	}
}

// InternalServerError returns status code 500 and error response.
func InternalServerError(err error) (int, Response) {
	return http.StatusInternalServerError, Response{
		Title:  "Internal Server Error",
		Detail: "An unexpected error occurred on the server.",
		Errors: err,
	}
}
