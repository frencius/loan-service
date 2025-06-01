package model

import "time"

type Investment struct {
	ID                           string
	LoanID                       string
	InvestorID                   string
	InvestedAmount               float64
	InvestmentAgreementLetterURL string
	IsInvestmentAggrementSigned  bool
	InvestmentAggrementSignedAt  *time.Time
	TotalProfit                  float64
}
