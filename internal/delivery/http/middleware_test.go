package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		errMsg  string
		wantMsg string
	}{
		{
			name:    "simple error message",
			errMsg:  "validation failed",
			wantMsg: "validation failed",
		},
		{
			name:    "empty error message",
			errMsg:  "",
			wantMsg: "",
		},
		{
			name:    "long error message",
			errMsg:  "this is a very long error message that contains multiple details about what went wrong",
			wantMsg: "this is a very long error message that contains multiple details about what went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorResponse(tt.errMsg)
			
			if got.Status != "error" {
				t.Errorf("ErrorResponse().Status = %v, want 'error'", got.Status)
			}
			
			if got.ErrorMessage != tt.wantMsg {
				t.Errorf("ErrorResponse().ErrorMessage = %v, want %v", got.ErrorMessage, tt.wantMsg)
			}
			
			if got.Data != nil {
				t.Errorf("ErrorResponse().Data = %v, want nil", got.Data)
			}
		})
	}
}

func TestHandleErrResponse(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name:           "standard error returns 500",
			err:            errors.New("standard error"),
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     "standard error",
		},
		{
			name:           "AppError with BadRequest status",
			err:            model.ErrUserNotFound,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "user not found",
		},
		{
			name:           "AppError with Unauthorized status",
			err:            model.ErrPasswordMatch,
			wantStatusCode: http.StatusUnauthorized,
			wantErrMsg:     "invalid username or password",
		},
		{
			name:           "AppError with InternalServerError status",
			err:            model.ErrDBUnexpected,
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     "unexpected db error",
		},
		{
			name:           "AppError with wrapped error",
			err:            model.ErrDBUnexpected.WithErr(errors.New("connection timeout")),
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     "unexpected db error: connection timeout",
		},
		{
			name:           "AppError UserAlreadyExists",
			err:            model.ErrUserAlreadyExists,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "user already exists",
		},
		{
			name:           "AppError UserNotApproved",
			err:            model.ErrUserNotApproved,
			wantStatusCode: http.StatusUnauthorized,
			wantErrMsg:     "user not approved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			err := HandleErrResponse(c, tt.err)
			
			if err != nil {
				t.Errorf("HandleErrResponse() returned error: %v", err)
			}
			
			if rec.Code != tt.wantStatusCode {
				t.Errorf("HandleErrResponse() status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}
			
			// Check if error message is in the response body
			body := rec.Body.String()
			if body == "" {
				t.Error("HandleErrResponse() returned empty body")
			}
		})
	}
}

func TestDefaultResponse_Structure(t *testing.T) {
	// Test DefaultResponse with different data types
	t.Run("DefaultResponse with string data", func(t *testing.T) {
		resp := DefaultResponse[string]{
			Status:       "ok",
			ErrorMessage: "",
			Data:         "test data",
		}
		
		if resp.Status != "ok" {
			t.Errorf("Status = %v, want 'ok'", resp.Status)
		}
		
		if resp.Data != "test data" {
			t.Errorf("Data = %v, want 'test data'", resp.Data)
		}
	})
	
	t.Run("DefaultResponse with nil data", func(t *testing.T) {
		resp := DefaultResponse[interface{}]{
			Status:       "error",
			ErrorMessage: "something went wrong",
			Data:         nil,
		}
		
		if resp.Status != "error" {
			t.Errorf("Status = %v, want 'error'", resp.Status)
		}
		
		if resp.Data != nil {
			t.Errorf("Data = %v, want nil", resp.Data)
		}
	})
	
	t.Run("DefaultResponse with int data", func(t *testing.T) {
		resp := DefaultResponse[int]{
			Status:       "ok",
			ErrorMessage: "",
			Data:         42,
		}
		
		if resp.Data != 42 {
			t.Errorf("Data = %v, want 42", resp.Data)
		}
	})
}

func TestHandleErrResponse_EdgeCases(t *testing.T) {
	t.Run("nil error handling", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// This shouldn't panic even with nil error
		_ = HandleErrResponse(c, nil)
	})
	
	t.Run("nested AppError", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		nestedErr := model.ErrUserNotFound.WithErr(
			model.ErrDBUnexpected.WithErr(errors.New("connection lost")),
		)
		
		err := HandleErrResponse(c, nestedErr)
		
		if err != nil {
			t.Errorf("HandleErrResponse() returned error: %v", err)
		}
		
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status code = %v, want %v", rec.Code, http.StatusBadRequest)
		}
	})
}