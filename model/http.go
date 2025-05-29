package model

import "net/http"

type (
	GenericResponse struct {
		Code    int            `json:"code,omitempty"`
		Message string         `json:"message,omitempty"`
		Result  interface{}    `json:"result,omitempty"`
		Error   *ErrorResponse `json:"error,omitempty"`
	}

	ErrorResponse struct {
		Detail  string `json:"detail,omitempty"`
		Message string `json:"message,omitempty"`
	}
)

func ComposeResponse(result interface{}, code int) GenericResponse {
	return GenericResponse{
		Code:    code,
		Message: http.StatusText(code),
		Result:  result,
	}
}

func ComposeErrorResponse(code int, errDetail string, errMessage string) GenericResponse {
	return GenericResponse{
		Code:    code,
		Message: http.StatusText(code),
		Error: &ErrorResponse{
			Detail:  errDetail,
			Message: errMessage,
		},
	}
}
