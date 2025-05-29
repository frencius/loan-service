# Loan Management System

A system that handles:
 * loan creation
 * loan approval
 * investment
 * loan disbursement

## Run Locally
``` 
    make run-local
```

## Test
```
    make test
```

## Test Coverage
```
    make test-coverage
```

## Notes
* Create Environment variable in `configuration/{env}.env`

## DB Migration
```sh
$ migrate -database "postgres://username:password@localhost:5432/db_name?sslmode=disable" -path db/migrations up

```

```sh
# create migration file
$ migrate create -ext sql -dir ./db/migrations -seq {your_migration_file_name}
```







## Analysis and Design
### Requirement Analysis
1. One way Loan state machine: proposed -> approved -> invested -> disbursed
2. each state has required data and validations
3. Investor and loan has many to many relationship

### Requirement Assumption
1. loan state can be seen from borrower and staff POV
2. added new canceled state for handling cancellations
3. added new rejected state if loan could not be approved
4. added new published state for when the loan is offered to intestors/ lenders
5. when published to investors/ lenders and got no investment, loan is canceled
6. invested state is when total_invested_amount == principal_amount

### state diagram:
[Loan State Machine](docs/state-diagram.png)

### Data Model Design
Entity/ object:
1. Loan
    properties:
        - id
        - borrower_id
        - principal_amount
        - total_invested_amount
        - interest_rate
        - roi_rate
        - state
        - visit_proof_url
        - validated_at 
        - validated_by 
        - loan_agreement_letter_url
        - is_loan_aggrement_signed
        - loan_aggrement_signed_at
        --audit--
        - created_at
        - created_by
        - approved_at
        - approved_by        
        - rejected_at
        - rejected_by
        - rejected_reason
        - canceled_at
        - canceled_by
        - canceled_reason
        - published_at
        - published_by
        - invested_at
        - disbursed_at
        - disbursed_by
    methods:
        - createLoan()
        - updateLoan()
        - getLoan()
        - rejectLoan()
        - cancelLoan()
        - approveLoan()
        - publishLoan()
        - disburseLoan()

2. Investment
    properties:
        - id
        - investor_id
        - loan_id
        - invested_amount
        - investment_agreement_letter_url
        - is_investment_aggrement_signed
        - investment_aggrement_signed_at
        - total_profit
    methods:
        - investLoan()
        - generateAgreement()
        - sendAgreement()

3. Investor
    properties:
        - id
        - name
        - nik
        - npwp
        - email
        - phone_number

4. Employee
    properties:
        - id
        - name
        - employee_number

5. Borrower
    properties
        - id
        - name
        - address
        - occupation
        - nik
        - dob

### API Design
    API:
        POST /v1/loans
            requestBody:
                - borrower_id
                - principal_amount
                - interest_rate
                - roi_rate
                - visit_proof_url
                - validated_at 
                - validated_by 
            response:
                - 201 Created:
                    - loan_id
                    - state: proposed
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - borrower_id is exist
                - principal_amount is not empty
                - interest_rate is not empty
                - roi_rate is not empty
        PUT /v1/loans/{id}
            requestBody:
                - borrower_id
                - principal_amount
                - interest_rate
                - roi_rate
                - visit_proof_url
                - validated_at 
                - validated_by 
                - loan_agreement_letter_url
                - is_loan_aggrement_signed
                - loan_aggrement_signed_at
            response:
                - 200 Success
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - loan id is exist
                - borrower_id is exist
                - basic validation (empty, number, string)
        
        GET /v1/loans/{id}
            response:
                - 200 Success:
                    - [all loan properties]
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - loan id is exist
        
        PATCH /v1/loans/{id}
            requestBody:
                - state: canceled | rejected | approved | publihsed | invested | disbursed 
            response:
                - 200 Success
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - loan id id exist
                - current state: proposed, eligible state [rejected, canceled, approved]
                - current state: approved, eligible state [canceled, published]
                - current state published, eligible state [canceled, invested]
                - current state invested, eligible state [canceled, disbursed]
                - current state disbursed, eligible state [canceled]
                - approved:
                    - current state is proposed
                    - visit_proof_url is not empty
                    - validated_at is not empty
                    - validated_by is not empty
                - published:
                    - current state is approved
                - disbursed:
                    - current state is invested
                    - loan_agreement_letter_url is not empty
                    - is_loan_aggrement_signed is not empty
                    - loan_aggrement_signed_at is not empty
                    - disbursed_at is not empty
                    - disbursed_by is not empty

        POST /v1/loans/{id}/investments
            requestBody:
                - investor_id
                - invested_amount
            response:
                - 201 Created:
                    - investment_id
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - investor_id is exist
                - loan_id is exist
                - invested_amount is not empty
            logic:
                - increment total_invested_amount in Loan for every investment creation
                - if (total_invested_amount == principal_amount && current state == published) then change loan state to invested
                - once invested, generate agreement letter:
                    - this could be built internally or using 3rd party integration like Privy, DocuSign, MekariSign, etc
                    - investor: include roi_rate
                    - borrower: include interest_rate
                - send aggreement letter url to investors email

        POST /v1/files
            - requestBody:
                - byte file
            response:
                - id
                - url