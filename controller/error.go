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
		errMsg = model.ErrorBorrowerNotFound.Error()
		respCode = http.StatusNotFound
	case model.ErrorLoanNotFound:
		errMsg = model.ErrorLoanNotFound.Error()
		respCode = http.StatusNotFound
	case model.ErrorLoanStateInvalid:
		errMsg = model.ErrorLoanStateInvalid.Error()
		respCode = http.StatusBadRequest
	case model.ErrorLoanStateTransitionNotAllowed:
		errMsg = model.ErrorLoanStateTransitionNotAllowed.Error()
		respCode = http.StatusBadRequest
	case model.ErrorStateTransitionRequirementNotFulfilled:
		errMsg = model.ErrorStateTransitionRequirementNotFulfilled.Error()
		respCode = http.StatusBadRequest
	case model.ErrorTransitionToTheSameState:
		errMsg = model.ErrorTransitionToTheSameState.Error()
		respCode = http.StatusBadRequest
	case model.ErrorInvestorNotFound:
		errMsg = model.ErrorInvestorNotFound.Error()
		respCode = http.StatusNotFound
	case model.ErrorStateMustBePublished:
		errMsg = model.ErrorStateMustBePublished.Error()
		respCode = http.StatusBadRequest
	case model.ErrorInvestmentNotFound:
		errMsg = model.ErrorInvestmentNotFound.Error()
		respCode = http.StatusNotFound
	case model.ErrorInvestmentExist:
		errMsg = model.ErrorInvestmentExist.Error()
		respCode = http.StatusBadRequest
	default:
		errMsg = "Something wrong in the system!"
		respCode = http.StatusInternalServerError
	}
	return errMsg, respCode
}
