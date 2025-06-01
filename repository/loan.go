package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type ILoanRepository interface {
	CreateLoan(ctx context.Context, loan *model.Loan) (ID string, err error)
	GetLoanByID(ctx context.Context, id string) (loan *model.Loan, err error)
	UpdateLoanState(ctx context.Context, loan *model.Loan, newLoanState model.LoanState) (err error)
	UpdateLoanTotalInvestedAmount(ctx context.Context, loan *model.Loan) (err error)
}

type LoanRepository struct {
	DB *sql.DB
}

func NewLoanRepository(app *application.App) ILoanRepository {
	return &LoanRepository{
		DB: app.DB,
	}
}

func (lr *LoanRepository) CreateLoan(ctx context.Context, loan *model.Loan) (ID string, err error) {
	query := `
		INSERT INTO
			loans (
				borrower_id,
				principal_amount,
				interest_rate,
				roi_rate,
				state,
				created_by
			)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id
		`

	err = lr.DB.QueryRowContext(ctx, query,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.ROIRate,
		loan.State,
		loan.CreatedBy,
	).Scan(&ID)

	if err != nil {
		log.Println("CreateLoan error ", err)
		return
	}

	return
}

func (lr *LoanRepository) GetLoanByID(ctx context.Context, id string) (loan *model.Loan, err error) {
	query := `
		SELECT
			id,
			borrower_id,
			principal_amount,
			interest_rate,
			roi_rate,
			state,
			COALESCE(total_invested_amount, 0),
			COALESCE(visit_proof_url, ''),
			COALESCE(validated_at, null),
			COALESCE(validated_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(loan_agreement_letter_url, ''),
			COALESCE(is_loan_aggrement_signed, false),
			COALESCE(loan_aggrement_signed_at, null),
			COALESCE(created_at, null),
			COALESCE(created_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(approved_at, null),
			COALESCE(approved_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(rejected_at, null),
			COALESCE(rejected_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(rejected_reason, ''),
			COALESCE(canceled_at, null),
			COALESCE(canceled_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(canceled_reason, ''),
			COALESCE(published_at, null),
			COALESCE(published_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(invested_at, null),
			COALESCE(disbursed_at, null),
			COALESCE(disbursed_by, '00000000-0000-0000-0000-000000000000'),
			COALESCE(updated_at, null)
		FROM
			loans 
		WHERE
			id = $1
		`

	loan = &model.Loan{}
	err = lr.DB.QueryRowContext(ctx, query, id).Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.InterestRate,
		&loan.ROIRate,
		&loan.State,
		&loan.TotalInvestedAmount,
		&loan.VisitProofURL,
		&loan.ValidatedAt,
		&loan.ValidatedBy,
		&loan.LoanAgreementLetterURL,
		&loan.IsLoanAggrementSigned,
		&loan.LoanAggrementSignedAt,
		&loan.CreatedAt,
		&loan.CreatedBy,
		&loan.ApprovedAt,
		&loan.ApprovedBy,
		&loan.RejectedAt,
		&loan.RejectedBy,
		&loan.RejectedReason,
		&loan.CanceledAt,
		&loan.CanceledBy,
		&loan.CanceledReason,
		&loan.PublishedAt,
		&loan.PublishedBy,
		&loan.InvestedAt,
		&loan.DisbursedAt,
		&loan.DisbursedBy,
		&loan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			err = model.ErrorLoanNotFound
			log.Println("GetLoanByID ", err)
			return
		}
		log.Println("GetLoanByID ", err)
		return
	}

	return
}

func (lr *LoanRepository) UpdateLoanState(ctx context.Context, loan *model.Loan, newLoanState model.LoanState) (err error) {

	query, args, err := lr.buildLoanUpdateQuery(loan, newLoanState, ctx.Value("userID").(string))
	if err != nil {
		log.Println("CreateLoan buildLoanUpdateQuery error ", err)
		return
	}

	rows, err := lr.DB.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("CreateLoan ExecContext error ", err)
		return
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		log.Println("CreateLoan RowsAffected error ", err)
		return
	}

	if affected < 1 {
		err = model.ErrorLoanNotFound
		log.Println("CreateLoan affected < 1 error ", err)
		return
	}

	return
}

func (lr *LoanRepository) buildLoanUpdateQuery(loan *model.Loan, newLoanState model.LoanState, userID string) (query string, args []any, err error) {
	setParts := []string{}
	idx := 1

	fields, ok := model.StateUpdates[newLoanState]
	if !ok {
		return
	}

	for _, field := range fields {
		switch field {
		case "state":
			setParts = append(setParts, fmt.Sprintf("state = $%d", idx))
			args = append(args, newLoanState)
		case "approved_by", "rejected_by", "canceled_by", "published_by", "disbursed_by":
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, idx))
			args = append(args, userID)
		default:
			// timestamp field
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, idx))
			args = append(args, time.Now())
		}
		idx++
	}

	// WHERE clause
	query = fmt.Sprintf("UPDATE loans SET %s WHERE id = $%d", strings.Join(setParts, ", "), idx)
	args = append(args, loan.ID)

	return
}

func (lr *LoanRepository) UpdateLoanTotalInvestedAmount(ctx context.Context, loan *model.Loan) (err error) {
	query := `
		UPDATE
			loans
		SET
			total_invested_amount = $2
		WHERE
			id = $1
	`
	rows, err := lr.DB.ExecContext(ctx, query, loan.ID, loan.TotalInvestedAmount)
	if err != nil {
		log.Println("UpdateLoanTotalInvestedAmount ExecContext error ", err)
		return
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		log.Println("UpdateLoanTotalInvestedAmount RowsAffected error ", err)
		return
	}

	if affected < 1 {
		err = model.ErrorLoanNotFound
		log.Println("UpdateLoanTotalInvestedAmount affected < 1 error ", err)
		return
	}

	return
}
