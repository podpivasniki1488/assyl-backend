package model

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError AppError
		want     string
	}{
		{
			name: "error with only message",
			appError: AppError{
				HttpStatusCode: http.StatusBadRequest,
				Message:        "test error",
				InternalError:  nil,
			},
			want: "test error",
		},
		{
			name: "error with message and internal error",
			appError: AppError{
				HttpStatusCode: http.StatusInternalServerError,
				Message:        "database error",
				InternalError:  errors.New("connection failed"),
			},
			want: "database error: connection failed",
		},
		{
			name: "error with empty message and internal error",
			appError: AppError{
				HttpStatusCode: http.StatusInternalServerError,
				Message:        "",
				InternalError:  errors.New("internal issue"),
			},
			want: ": internal issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appError.Error()
			if got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name           string
		httpStatusCode int
		err            error
		wantMessage    string
		wantStatusCode int
	}{
		{
			name:           "create app error from standard error",
			httpStatusCode: http.StatusNotFound,
			err:            errors.New("not found"),
			wantMessage:    "not found",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "create app error with 500 status",
			httpStatusCode: http.StatusInternalServerError,
			err:            errors.New("server error"),
			wantMessage:    "server error",
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAppError(tt.httpStatusCode, tt.err)
			
			if got.HttpStatusCode != tt.wantStatusCode {
				t.Errorf("NewAppError().HttpStatusCode = %v, want %v", got.HttpStatusCode, tt.wantStatusCode)
			}
			
			if got.Message != tt.wantMessage {
				t.Errorf("NewAppError().Message = %v, want %v", got.Message, tt.wantMessage)
			}
			
			if got.InternalError != tt.err {
				t.Errorf("NewAppError().InternalError = %v, want %v", got.InternalError, tt.err)
			}
		})
	}
}

func TestAppError_WithErr(t *testing.T) {
	tests := []struct {
		name          string
		baseError     AppError
		internalError error
		wantMessage   string
	}{
		{
			name: "add internal error to base error",
			baseError: AppError{
				HttpStatusCode: http.StatusBadRequest,
				Message:        "validation failed",
				InternalError:  nil,
			},
			internalError: errors.New("field 'username' is required"),
			wantMessage:   "validation failed: field 'username' is required",
		},
		{
			name: "replace internal error",
			baseError: AppError{
				HttpStatusCode: http.StatusInternalServerError,
				Message:        "db error",
				InternalError:  errors.New("old error"),
			},
			internalError: errors.New("new error"),
			wantMessage:   "db error: new error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.baseError.WithErr(tt.internalError)
			
			if got.InternalError != tt.internalError {
				t.Errorf("WithErr().InternalError = %v, want %v", got.InternalError, tt.internalError)
			}
			
			if got.Error() != tt.wantMessage {
				t.Errorf("WithErr().Error() = %v, want %v", got.Error(), tt.wantMessage)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        AppError
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "ErrDBUnexpected",
			err:        ErrDBUnexpected,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "unexpected db error",
		},
		{
			name:       "ErrPasswordMatch",
			err:        ErrPasswordMatch,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "invalid username or password",
		},
		{
			name:       "ErrUserNotFound",
			err:        ErrUserNotFound,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "user not found",
		},
		{
			name:       "ErrUserAlreadyExists",
			err:        ErrUserAlreadyExists,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "user already exists",
		},
		{
			name:       "ErrRecordNotFound",
			err:        ErrRecordNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "record not found",
		},
		{
			name:       "ErrUserNotApproved",
			err:        ErrUserNotApproved,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "user not approved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.HttpStatusCode != tt.wantStatus {
				t.Errorf("%s.HttpStatusCode = %v, want %v", tt.name, tt.err.HttpStatusCode, tt.wantStatus)
			}
			
			if tt.err.Message != tt.wantMsg {
				t.Errorf("%s.Message = %v, want %v", tt.name, tt.err.Message, tt.wantMsg)
			}
			
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("%s.Error() = %v, want %v", tt.name, tt.err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestAppError_ErrorsAs(t *testing.T) {
	// Test that AppError can be used with errors.As
	baseErr := ErrUserNotFound.WithErr(errors.New("database connection lost"))
	
	var appErr AppError
	if !errors.As(baseErr, &appErr) {
		t.Error("errors.As should recognize AppError")
	}
	
	if appErr.HttpStatusCode != http.StatusBadRequest {
		t.Errorf("errors.As extracted wrong status code: got %v, want %v", appErr.HttpStatusCode, http.StatusBadRequest)
	}
}

func TestAppError_WithErr_Immutability(t *testing.T) {
	// Test that WithErr doesn't mutate the original error
	original := ErrDBUnexpected
	originalInternal := original.InternalError
	
	modified := original.WithErr(errors.New("new internal error"))
	
	if original.InternalError != originalInternal {
		t.Error("WithErr mutated the original error")
	}
	
	if modified.InternalError == nil {
		t.Error("WithErr did not set internal error")
	}
	
	if original.InternalError == modified.InternalError {
		t.Error("WithErr returned the same internal error reference")
	}
}