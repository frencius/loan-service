package controller

import (
	"encoding/json"
	"net/http"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/service"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type ILoanController interface {
	CreateLoan(w http.ResponseWriter, r *http.Request)
	UpdateLoanState(w http.ResponseWriter, r *http.Request)
	CreateLoanInvestment(w http.ResponseWriter, r *http.Request)
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

func (acc *LoanController) UpdateLoanState(w http.ResponseWriter, r *http.Request) {
	// decode body request
	updateLoanStateRequest := model.UpdateLoanStateRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateLoanStateRequest)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// validate request
	valid, err := model.IsValid(updateLoanStateRequest)
	if !valid {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// get loan id path param
	loanID := chi.URLParam(r, "id")
	_, err = uuid.Parse(loanID)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Loan ID invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	updateLoanStateRequest.LoanID = loanID

	// call business logic
	resp, err := acc.LoanService.UpdateLoanState(r.Context(), &updateLoanStateRequest)
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

func (acc *LoanController) CreateLoanInvestment(w http.ResponseWriter, r *http.Request) {
	// decode body request
	createLoanInvestmentRequest := model.CreateLoanInvestmentRequest{}
	err := json.NewDecoder(r.Body).Decode(&createLoanInvestmentRequest)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// validate request
	valid, err := model.IsValid(createLoanInvestmentRequest)
	if !valid {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Request body invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	_, err = uuid.Parse(createLoanInvestmentRequest.InvestorID)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Investor ID invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	// get loan id path param
	loanID := chi.URLParam(r, "id")
	_, err = uuid.Parse(loanID)
	if err != nil {
		respCode := http.StatusBadRequest
		result := model.ComposeErrorResponse(respCode, err.Error(), "Loan ID invalid")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	createLoanInvestmentRequest.LoanID = loanID

	// call business logic
	resp, err := acc.LoanService.CreateLoanInvestment(r.Context(), &createLoanInvestmentRequest)
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
