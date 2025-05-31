package controller

import (
	"encoding/json"
	"net/http"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/service"
	"github.com/google/uuid"
)

type ILoanController interface {
	CreateLoan(w http.ResponseWriter, r *http.Request)
}

type LoanController struct {
	LoanService service.ILoanService
}

func NewLoanController(app *application.App) ILoanController {
	return &LoanController{
		LoanService: service.NewLoanService(app),
	}
}

func (acc *LoanController) CreateLoan(w http.ResponseWriter, r *http.Request) {
	// decode body request
	createLoanRequest := model.CreateLoanRequest{}
	err := json.NewDecoder(r.Body).Decode(&createLoanRequest)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	_, err = uuid.Parse(createLoanRequest.BorrowerID)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Borrower ID invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// validate request
	valid, err := model.IsValid(createLoanRequest)
	if !valid {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// call business logic
	resp, err := acc.LoanService.CreateLoan(r.Context(), &createLoanRequest)
	if err != nil {
		errMsg, respCode := getErrorResponse(err)
		result := model.ComposeErrorResponse(respCode, err.Error(), errMsg)
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// return response
	respCode := http.StatusOK
	result := model.ComposeResponse(resp, respCode)
	WriteHTTPResponse(w, respCode, result)
}
