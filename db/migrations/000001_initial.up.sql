-- audit support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- function to auto update updated_at field
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- create enum type for loan state
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'loan_state') THEN
    CREATE TYPE loan_state AS ENUM ('canceled', 'rejected', 'proposed', 'approved', 'published', 'invested', 'disbursed');
  END IF;
END
$$;

-- 1. Borrowers
CREATE TABLE borrowers (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  address TEXT NOT NULL,
  occupation TEXT NOT NULL,
  nik VARCHAR(16) NOT NULL UNIQUE,
  dob DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Triggers for updated_at
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON borrowers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- 2. Investors
CREATE TABLE investors (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  nik VARCHAR(16) NOT NULL UNIQUE,
  npwp VARCHAR(16) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  phone_number VARCHAR(15),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Triggers for updated_at
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON investors
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- 3. Employees
CREATE TABLE employees (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  employee_number VARCHAR(50) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Triggers for updated_at
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON employees
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- 4. Loans
CREATE TABLE loans (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  borrower_id UUID NOT NULL,
  principal_amount NUMERIC(20,2) NOT NULL,
  total_invested_amount NUMERIC(20,2) DEFAULT 0.00,
  interest_rate NUMERIC(5,2) NOT NULL,
  roi_rate NUMERIC(5,2) NOT NULL,
  state loan_state NOT NULL,
  visit_proof_url TEXT,
  validated_at TIMESTAMP,
  validated_by UUID,
  loan_agreement_letter_url TEXT,
  is_loan_aggrement_signed BOOLEAN DEFAULT FALSE,
  loan_aggrement_signed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by UUID,
  approved_at TIMESTAMP,
  approved_by UUID,
  rejected_at TIMESTAMP,
  rejected_by UUID,
  rejected_reason TEXT,
  canceled_at TIMESTAMP,
  canceled_by UUID,
  canceled_reason TEXT,
  published_at TIMESTAMP,
  published_by UUID,
  invested_at TIMESTAMP,
  disbursed_at TIMESTAMP,
  disbursed_by UUID,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_loans_borrower FOREIGN KEY (borrower_id) REFERENCES borrowers(id),
  CONSTRAINT fk_loans_validated_by FOREIGN KEY (validated_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_created_by FOREIGN KEY (created_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_rejected_by FOREIGN KEY (rejected_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_canceled_by FOREIGN KEY (canceled_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_published_by FOREIGN KEY (published_by) REFERENCES employees(id),
  CONSTRAINT fk_loans_disbursed_by FOREIGN KEY (disbursed_by) REFERENCES employees(id)
);

-- Triggers for updated_at
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON loans
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


-- 5. Investments
CREATE TABLE investments (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  investor_id UUID NOT NULL,
  loan_id UUID NOT NULL,
  invested_amount NUMERIC(20,2) NOT NULL,
  investment_agreement_letter_url TEXT,
  is_investment_aggrement_signed BOOLEAN DEFAULT FALSE,
  investment_aggrement_signed_at TIMESTAMP,
  total_profit NUMERIC(20,2) DEFAULT 0.00,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_investments_investor FOREIGN KEY (investor_id) REFERENCES investors(id),
  CONSTRAINT fk_investments_loan FOREIGN KEY (loan_id) REFERENCES loans(id)
);

-- Triggers for updated_at
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON investments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for fast lookup
CREATE INDEX idx_loans_state ON loans(state);
CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_created_at ON loans(created_at);

CREATE INDEX idx_investments_loan_id ON investments(loan_id);
CREATE INDEX idx_investments_investor_id ON investments(investor_id);

CREATE INDEX idx_employees_employee_number ON employees(employee_number);
CREATE INDEX idx_investors_email ON investors(email);
CREATE INDEX idx_borrowers_nik ON borrowers(nik);
