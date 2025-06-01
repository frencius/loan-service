mockgen -source=./repository/loan.go -destination=./mock/mock_loan_repository.go -package=mock
mockgen -source=./repository/borrower.go -destination=./mock/mock_borrower_repository.go -package=mock