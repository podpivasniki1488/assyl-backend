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
	ErrDBUnexpected   = AppError{HttpStatusCode: http.StatusInternalServerError, Message: "unexpected db error"}
	ErrRecordNotFound = AppError{HttpStatusCode: http.StatusNotFound, Message: "record not found"}
	ErrInvalidInput   = AppError{HttpStatusCode: http.StatusBadRequest, Message: "invalid input"}

	ErrPasswordMatch     = AppError{HttpStatusCode: http.StatusUnauthorized, Message: "invalid username or password"}
	ErrUserNotFound      = AppError{HttpStatusCode: http.StatusNotFound, Message: "user not found"}
	ErrUserAlreadyExists = AppError{HttpStatusCode: http.StatusBadRequest, Message: "user already exists"}
	ErrUserNotApproved   = AppError{HttpStatusCode: http.StatusUnauthorized, Message: "user not approved"}

	ErrApartmentNotFound     = AppError{HttpStatusCode: http.StatusNotFound, Message: "allocation not found"}
	ErrApartmentAlreadyBound = AppError{HttpStatusCode: http.StatusConflict, Message: "apartment already bound"}
	ErrReservationNotFound   = AppError{HttpStatusCode: http.StatusNotFound, Message: "record not found"}
	ErrCinemaBusy            = AppError{HttpStatusCode: http.StatusConflict, Message: "cinema busy"}
	ErrTooManyPeople         = AppError{HttpStatusCode: http.StatusBadRequest, Message: "too many people"}
	ErrReservationImpossible = AppError{HttpStatusCode: http.StatusBadRequest, Message: "reservation impossible"}
	ErrAdminsCannotBeDeleted = AppError{HttpStatusCode: http.StatusForbidden, Message: "admins can not be deleted"}

	ErrChatNotFound         = AppError{HttpStatusCode: http.StatusNotFound, Message: "chat not found"}
	ErrUserNotAllowed       = AppError{HttpStatusCode: http.StatusForbidden, Message: "user not allowed to send message to this chat"}
	ErrSingleChatUser       = AppError{HttpStatusCode: http.StatusBadRequest, Message: "chat must include at least 2 users"}
	ErrCannotHaveDuplicates = AppError{HttpStatusCode: http.StatusBadRequest, Message: "cannot have duplicates in chat participants"}
)
