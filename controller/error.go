package controller

import (
	"net/http"
)

func getErrorResponse(err error) (string, int) {
	var errMsg string

	respCode := http.StatusBadRequest
	switch err {
	default:
		errMsg = "Something wrong in the system!"
		respCode = http.StatusInternalServerError
	}
	return errMsg, respCode
}
