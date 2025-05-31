package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type ILoanRepository interface {
	CreateLoan(ctx context.Context, loan *model.Loan) (ID string, err error)
}

type LoanRepository struct {
	DB *sql.DB
}

func NewLoanRepository(app *application.App) ILoanRepository {
	return &LoanRepository{
		DB: app.DB,
	}
}

func (acr *LoanRepository) CreateLoan(ctx context.Context, loan *model.Loan) (ID string, err error) {
	query := `
		INSERT INTO
			loans (
				borrower_id,
				principal_amount,
				interest_rate,
				roi_rate,
				state
			)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING
			id
		`

	err = acr.DB.QueryRowContext(ctx, query,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.ROIRate,
		loan.State,
	).Scan(&ID)

	if err != nil {
		log.Println("CreateLoan error ", err)
		return
	}

	return
}
