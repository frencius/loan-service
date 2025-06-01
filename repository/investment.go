package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type IInvestmentRepository interface {
	CreateInvestment(ctx context.Context, investment *model.Investment) (ID string, err error)
	GetInvestmentByInvestorID(ctx context.Context, id string) (investment *model.Investment, err error)
}

type InvestmentRepository struct {
	DB *sql.DB
}

func NewInvestmentRepository(app *application.App) IInvestmentRepository {
	return &InvestmentRepository{
		DB: app.DB,
	}
}

func (ir *InvestmentRepository) CreateInvestment(ctx context.Context, investment *model.Investment) (ID string, err error) {
	query := `
		INSERT INTO
			investments (
				loan_id,
				investor_id,
				invested_amount
			)
		VALUES 
			($1, $2, $3)
		RETURNING
			id
		`

	err = ir.DB.QueryRowContext(ctx, query,
		investment.LoanID,
		investment.InvestorID,
		investment.InvestedAmount,
	).Scan(&ID)

	if err != nil {
		log.Println("CreateInvestment error ", err)
		return
	}

	return
}

func (ir *InvestmentRepository) GetInvestmentByInvestorID(ctx context.Context, id string) (investment *model.Investment, err error) {
	query := `
		SELECT
			id,
			investor_id,
			loan_id,
			invested_amount,
			COALESCE(investment_agreement_letter_url, ''),
			COALESCE(is_investment_aggrement_signed, false),
			COALESCE(investment_aggrement_signed_at, null),
			COALESCE(total_profit, 0)
		FROM
			investments
		WHERE
			investor_id = $1
	`

	investment = &model.Investment{}
	err = ir.DB.QueryRowContext(ctx, query, id).Scan(
		&investment.ID,
		&investment.InvestorID,
		&investment.LoanID,
		&investment.InvestedAmount,
		&investment.InvestmentAgreementLetterURL,
		&investment.IsInvestmentAggrementSigned,
		&investment.InvestmentAggrementSignedAt,
		&investment.TotalProfit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("GetInvestmentByID ", err)
			err = model.ErrorInvestmentNotFound
			return
		}

		log.Println("GetInvestmentByID ", err)
		return
	}

	return
}
