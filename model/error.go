package model

import "errors"

var (
	ErrorBorrowerNotFound                       = errors.New("borrower is not found")
	ErrorLoanNotFound                           = errors.New("loan is not found")
	ErrorLoanStateInvalid                       = errors.New("loan state invalid")
	ErrorLoanStateTransitionNotAllowed          = errors.New("loan state transition is not allowed")
	ErrorStateTransitionRequirementNotFulfilled = errors.New("loan state transition requirement is not fulfilled")
	ErrorTransitionToTheSameState               = errors.New("loan state transition could not be in the same sate")
	ErrorInvestorNotFound                       = errors.New("investor is not found")
	ErrorStateMustBePublished                   = errors.New("loan state must be publihsed")
	ErrorInvestmentExist                        = errors.New("investment exist")
	ErrorInvestmentNotFound                     = errors.New("investment is not found")
)
