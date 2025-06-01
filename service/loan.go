package service

import (
	"context"

	"slices"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/repository"
)

type ILoanService interface {
	CreateLoan(ctx context.Context, createLoanRequest *model.CreateLoanRequest) (createLoanResponse *model.CreateLoanResponse, err error)
	UpdateLoanState(ctx context.Context, updateLoanStateRequest *model.UpdateLoanStateRequest) (updateLoanStateResponse *model.UpdateLoanStateResponse, err error)
	CreateLoanInvestment(ctx context.Context, createLoanInvestmentRequest *model.CreateLoanInvestmentRequest) (createLoanInvestmentResponse *model.CreateLoanInvestmentResponse, err error)
}

type LoanService struct {
	LoanRepository       repository.ILoanRepository
	BorrowerRepository   repository.IBorrowerRepository
	InvestorRepository   repository.IInvestorRepository
	InvestmentRepository repository.IInvestmentRepository
}

func NewLoanService(app *application.App) ILoanService {
	return &LoanService{
		LoanRepository:       repository.NewLoanRepository(app),
		BorrowerRepository:   repository.NewBorrowerRepository(app),
		InvestorRepository:   repository.NewInvestorRepository(app),
		InvestmentRepository: repository.NewInvestmentRepository(app),
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
		CreatedBy:       ctx.Value("userID").(string),
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

func (ls *LoanService) UpdateLoanState(ctx context.Context, updateLoanStateRequest *model.UpdateLoanStateRequest) (updateLoanStateResponse *model.UpdateLoanStateResponse, err error) {
	loanID := updateLoanStateRequest.LoanID
	newLoanState := model.LoanState(updateLoanStateRequest.State)

	// validate loan id
	loan, err := ls.LoanRepository.GetLoanByID(ctx, loanID)
	if err != nil {
		return
	}

	// validate loan state
	if !model.ValidLoanState[newLoanState] {
		err = model.ErrorLoanStateInvalid
		return
	}

	// validate state transition
	if !ls.canStateTransition(loan.State, newLoanState) {
		err = model.ErrorLoanStateTransitionNotAllowed
		return
	}

	// state validation
	if !ls.isStateRequirementValid(loan, newLoanState) {
		err = model.ErrorStateTransitionRequirementNotFulfilled
		return
	}

	// update state
	err = ls.LoanRepository.UpdateLoanState(ctx, loan, newLoanState)
	if err != nil {
		return
	}

	return
}

func (ls *LoanService) canStateTransition(from, to model.LoanState) bool {
	validNextStates, ok := model.ValidStateTransitions[from]
	if !ok {
		return false
	}

	return slices.Contains(validNextStates, to)
}

func (ls *LoanService) isStateRequirementValid(loan *model.Loan, loanState model.LoanState) (valid bool) {
	switch loanState {
	case model.LoanStateApproved:
		if loan.State == model.LoanStateProposed && loan.VisitProofURL != "" && loan.ValidatedAt != nil && loan.ValidatedBy != "" {
			return true
		}
	case model.LoanStatePublished:
		if loan.State == model.LoanStateApproved {
			return true
		}
	case model.LoanStateDisbursed:
		if loan.State == model.LoanStateInvested && loan.LoanAgreementLetterURL != "" &&
			loan.IsLoanAggrementSigned && loan.LoanAggrementSignedAt != nil && loan.DisbursedAt != nil &&
			loan.DisbursedBy != "" {
			return true
		}
	case model.LoanStateInvested:
		if loan.State == model.LoanStatePublished && loan.TotalInvestedAmount >= loan.PrincipalAmount {
			return true
		}
	case model.LoanStateCanceled:
		if loan.CanceledReason != "" {
			return true
		}
	case model.LoanStateRejected:
		if loan.RejectedReason != "" {
			return true
		}
	}

	return
}

func (ls *LoanService) CreateLoanInvestment(ctx context.Context, createLoanInvestmentRequest *model.CreateLoanInvestmentRequest) (createLoanInvestmentResponse *model.CreateLoanInvestmentResponse, err error) {
	loanID := createLoanInvestmentRequest.LoanID
	investorID := createLoanInvestmentRequest.InvestorID
	investedAmount := createLoanInvestmentRequest.InvestmentAmount

	// validate loan
	loan, err := ls.LoanRepository.GetLoanByID(ctx, loanID)
	if err != nil {
		return
	}

	if loan.State != model.LoanStatePublished {
		err = model.ErrorStateMustBePublished
		return
	}

	// validate investor
	investor, err := ls.InvestorRepository.GetInvestorByID(ctx, investorID)
	if err != nil {
		return
	}

	// validate existing investment
	investment, err := ls.InvestmentRepository.GetInvestmentByInvestorID(ctx, investorID)
	if err != nil && err != model.ErrorInvestmentNotFound {
		return
	}

	if investment != nil && err == nil {
		err = model.ErrorInvestmentExist
		return
	}

	// create investment
	newInvestment := &model.Investment{
		LoanID:         loan.ID,
		InvestorID:     investor.ID,
		InvestedAmount: investedAmount,
	}

	investmentID, err := ls.InvestmentRepository.CreateInvestment(ctx, newInvestment)
	if err != nil {
		return
	}

	// update loan total invested amount
	loan.TotalInvestedAmount += investedAmount
	err = ls.LoanRepository.UpdateLoanTotalInvestedAmount(ctx, loan)
	if err != nil {
		return
	}

	// update state if eligible
	oldLoanState := loan.State
	if loan.TotalInvestedAmount >= loan.PrincipalAmount {
		loan.State = model.LoanStateInvested
		// TODO:
		// 1. generate aggrement letter
		// 2. send aggreement letter url to investors email
	}

	if oldLoanState != loan.State {
		updateLoanStateRequest := model.UpdateLoanStateRequest{
			LoanID: loan.ID,
			State:  string(loan.State),
		}

		_, err = ls.UpdateLoanState(ctx, &updateLoanStateRequest)
		if err != nil {
			return
		}
	}

	createLoanInvestmentResponse = &model.CreateLoanInvestmentResponse{
		InvestmentID: investmentID,
	}

	return
}
