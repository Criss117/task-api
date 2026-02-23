package commonreponse

import "net/http"

type CommonResponse struct {
	Status           int                 `json:"status"`
	Message          string              `json:"message"`
	Data             any                 `json:"data"`
	Error            *string             `json:"error"`
	ValidationErrors map[string][]string `json:"validation_errors"`
}

func InvalidBody(message string, errors map[string][]string) *CommonResponse {
	return &CommonResponse{
		Status:           http.StatusBadRequest,
		Message:          message,
		Error:            &message,
		ValidationErrors: errors,
	}
}

func Ok(message string, data any) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func NoContent(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusNoContent,
		Message: message,
	}
}

func Deleted(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusNoContent,
		Message: message,
	}
}

func Created(message string, data any) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusCreated,
		Message: message,
		Data:    data,
	}
}

func InternalServerError(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   &message,
	}
}

func NotFound(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusNotFound,
		Message: message,
		Error:   &message,
	}
}

func BadRequest(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   &message,
	}
}

func Unauthorized(message string) *CommonResponse {
	return &CommonResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
		Error:   &message,
	}
}
