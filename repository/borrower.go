package repository

import (
	"context"
	"database/sql"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type IBorrowerRepository interface {
	GetBorrowerByID(ctx context.Context, id string) (borrower *model.Borrower, err error)
}

type BorrowerRepository struct {
	DB *sql.DB
}

func NewBorrowerRepository(app *application.App) IBorrowerRepository {
	return &BorrowerRepository{
		DB: app.DB,
	}
}

func (acr *BorrowerRepository) GetBorrowerByID(ctx context.Context, id string) (borrower *model.Borrower, err error) {
	query := `
		SELECT
			id,
			name,
			address,
			occupation,
			nik,
			dob
		FROM
			borrowers
		WHERE
			id = $1
	`

	borrower = &model.Borrower{}
	err = acr.DB.QueryRowContext(ctx, query, id).Scan(
		&borrower.ID,
		&borrower.Name,
		&borrower.Address,
		&borrower.Occupation,
		&borrower.NIK,
		&borrower.DOB,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = model.ErrorBorrowerNotFound
			return
		}
		return
	}

	return
}
