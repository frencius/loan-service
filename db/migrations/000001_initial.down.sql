DROP TRIGGER IF EXISTS set_timestamp ON investments;
DROP TRIGGER IF EXISTS set_timestamp ON loans;
DROP TRIGGER IF EXISTS set_timestamp ON employees;
DROP TRIGGER IF EXISTS set_timestamp ON investors;
DROP TRIGGER IF EXISTS set_timestamp ON borrowers;

DROP INDEX IF EXISTS idx_loans_state;
DROP INDEX IF EXISTS idx_loans_borrower_id;
DROP INDEX IF EXISTS idx_loans_created_at;

DROP INDEX IF EXISTS idx_investments_loan_id;
DROP INDEX IF EXISTS idx_investments_investor_id;

DROP INDEX IF EXISTS idx_employees_employee_number;
DROP INDEX IF EXISTS idx_investors_email;
DROP INDEX IF EXISTS idx_borrowers_nik;

DROP TABLE IF EXISTS investments;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS investors;
DROP TABLE IF EXISTS borrowers;

DROP FUNCTION IF EXISTS trigger_set_timestamp;
DROP TYPE IF EXISTS loan_state;
DROP EXTENSION IF EXISTS "uuid-ossp";