package controller

import (
	"net/http"

	"github.com/frencius/loan-service/model"
)

func getErrorResponse(err error) (string, int) {
	var errMsg string

	respCode := http.StatusBadRequest
	switch err {
	case model.ErrorBorrowerNotFound:
		errMsg = "borrower is not found"
		respCode = http.StatusNotFound
	default:
		errMsg = "Something wrong in the system!"
		respCode = http.StatusInternalServerError
	}
	return errMsg, respCode
}
