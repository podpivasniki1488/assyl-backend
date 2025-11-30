package model

import "net/http"

type AppError struct {
	HttpStatusCode int
	Message        string
	InternalError  error
}

func NewAppError(httpStatusCode int, err error) *AppError {
	return &AppError{
		HttpStatusCode: httpStatusCode,
		Message:        err.Error(),
		InternalError:  err,
	}
}

func (e AppError) Error() string {
	msg := e.Message
	if e.InternalError != nil {
		msg += ": " + e.InternalError.Error()
	}

	return msg
}

func (e AppError) WithErr(err error) AppError {
	e.InternalError = err
	return e
}

var (
	ErrDBUnexpected      = AppError{HttpStatusCode: http.StatusInternalServerError, Message: "unexpected db error"}
	ErrPasswordMatch     = AppError{HttpStatusCode: http.StatusUnauthorized, Message: "invalid username or password"}
	ErrUserNotFound      = AppError{HttpStatusCode: http.StatusNotFound, Message: "user not found"}
	ErrUserAlreadyExists = AppError{HttpStatusCode: http.StatusBadRequest, Message: "user already exists"}
	ErrRecordNotFound    = AppError{HttpStatusCode: http.StatusNotFound, Message: "record not found"}
	ErrUserNotApproved   = AppError{HttpStatusCode: http.StatusUnauthorized, Message: "user not approved"}
)
