package service

import (
	"context"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/repository"
)

type ILoanService interface {
	CreateLoan(ctx context.Context, createLoanRequest *model.CreateLoanRequest) (createLoanResponse *model.CreateLoanResponse, err error)
}

type LoanService struct {
	LoanRepository     repository.ILoanRepository
	BorrowerRepository repository.IBorrowerRepository
}

func NewLoanService(app *application.App) ILoanService {
	return &LoanService{
		LoanRepository:     repository.NewLoanRepository(app),
		BorrowerRepository: repository.NewBorrowerRepository(app),
	}
}

func (ls *LoanService) CreateLoan(ctx context.Context, createLoanRequest *model.CreateLoanRequest) (createLoanResponse *model.CreateLoanResponse, err error) {
	// validate borrower_id is exist
	_, err = ls.BorrowerRepository.GetBorrowerByID(ctx, createLoanRequest.BorrowerID)
	if err != nil {
		return
	}

	loan := &model.Loan{
		BorrowerID:      createLoanRequest.BorrowerID,
		PrincipalAmount: createLoanRequest.PrincipalAmount,
		InterestRate:    createLoanRequest.InterestRate,
		ROIRate:         createLoanRequest.ROIRate,
		State:           model.LoanStateProposed,
	}

	loanID, err := ls.LoanRepository.CreateLoan(ctx, loan)
	if err != nil {
		return
	}

	createLoanResponse = &model.CreateLoanResponse{
		LoanID: loanID,
		State:  string(loan.State),
	}

	return
}
