package model

import "time"

type LoanState string

const (
	LoanStateProposed  LoanState = "proposed"
	LoanStateRejected  LoanState = "rejected"
	LoanStateCanceled  LoanState = "canceled"
	LoanStateApproved  LoanState = "approved"
	LoanStatePublished LoanState = "published"
	LoanStateInvested  LoanState = "invested"
	LoanStateDisbursed LoanState = "disbursed"
)

var (
	ValidLoanState = map[LoanState]bool{
		LoanStateCanceled:  true,
		LoanStateRejected:  true,
		LoanStateProposed:  true,
		LoanStateApproved:  true,
		LoanStatePublished: true,
		LoanStateInvested:  true,
		LoanStateDisbursed: true,
	}
)

// data model
type (
	Loan struct {
		ID                     string
		BorrowerID             string
		PrincipalAmount        float64
		TotalInvestedAmount    float64
		InterestRate           float64
		ROIRate                float64
		State                  LoanState
		VisitProofURL          string
		ValidatedAt            *time.Time
		ValidatedBy            string
		LoanAgreementLetterURL string
		IsLoanAggrementSigned  bool
		LoanAggrementSignedAt  *time.Time
		CreatedAt              *time.Time
		CreatedBy              string
		ApprovedAt             *time.Time
		ApprovedBy             string
		RejectedAt             *time.Time
		RejectedBy             string
		RejectedReason         string
		CanceledAt             *time.Time
		CanceledBy             string
		CanceledReason         string
		PublishedAt            *time.Time
		PublishedBy            string
		InvestedAt             *time.Time
		DisbursedAt            *time.Time
		DisbursedBy            string
		UpdatedAt              *time.Time
	}
)

// request response
type (
	CreateLoanRequest struct {
		BorrowerID      string  `json:"borrower_id" validate:"required"`
		PrincipalAmount float64 `json:"principal_amount" validate:"required"`
		InterestRate    float64 `json:"interest_rate" validate:"required"`
		ROIRate         float64 `json:"roi_rate" validate:"required"`
	}

	CreateLoanResponse struct {
		LoanID string `json:"loan_id"`
		State  string `json:"state"`
	}
)
