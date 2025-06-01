mockgen -source=./repository/loan.go -destination=./mock/mock_loan_repository.go -package=mock
mockgen -source=./repository/borrower.go -destination=./mock/mock_borrower_repository.go -package=mock
mockgen -source=./repository/investor.go -destination=./mock/mock_investor_repository.go -package=mock
mockgen -source=./repository/investment.go -destination=./mock/mock_investment_repository.go -package=mock