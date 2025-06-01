package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type IInvestorRepository interface {
	GetInvestorByID(ctx context.Context, id string) (investor *model.Investor, err error)
}

type InvestorRepository struct {
	DB *sql.DB
}

func NewInvestorRepository(app *application.App) IInvestorRepository {
	return &InvestorRepository{
		DB: app.DB,
	}
}

func (ir *InvestorRepository) GetInvestorByID(ctx context.Context, id string) (investor *model.Investor, err error) {
	query := `
		SELECT
			id,
			name,
			nik,
			npwp,
			email,
			phone_number,
			created_at,
			updated_at
		FROM
			investors
		WHERE
			id = $1
	`

	investor = &model.Investor{}
	err = ir.DB.QueryRowContext(ctx, query, id).Scan(
		&investor.ID,
		&investor.Name,
		&investor.NIK,
		&investor.NPWP,
		&investor.Email,
		&investor.PhoneNumber,
		&investor.CreatedAt,
		&investor.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("GetInvestorByID ", err)
			err = model.ErrorInvestorNotFound
			return
		}

		log.Println("GetInvestorByID ", err)
		return
	}

	return
}
