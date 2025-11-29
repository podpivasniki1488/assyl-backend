package http

type DefaultResponse[T any] struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	Data         T      `json:"data"`
}

func ErrorResponse(errMsg string) DefaultResponse[error] {
	return DefaultResponse[error]{
		Status:       "error",
		ErrorMessage: errMsg,
		Data:         nil,
	}
}
