# Loan Management System

## Problem Statement
### Example 3: Loan Service (system design and abstraction)

we are building a loan engine. A loan can have a multiple states: proposed , approved, invested, disbursed. the rule of state:
1. proposed is the initial state (when loan created it will has proposed state):
2. approved is once it approved by our staff.
    a. a approval must contains several information:
        i. the picture proof of the a field validator has visited the borrower 
        ii. the employee id of field validator
        iii. date of approval
    b. once approved it can not go back to proposed state
    c. once approved loan is ready to be offered to investors/lender
3. invested is once total amount of invested is equal the loan principal
    a. loan can have multiple investors, each with each their own amount
    b. total of invested amount can not be bigger than the loan principal amount
    c. once invested all investors will receive an email containing link to agreement letter (pdf)
4. disbursed is when is loan is given to borrower.
    a. a disbursement must contains several information:
        i. the loan agreement letter signed by borrower (pdf/jpeg)
        ii. the employee id of the field officer that hands the money and/or collect the agreement letter
        iii. date of disbursement

movement between state can only move forward, and a loan only need following information:
- borrower id number
- principal amount
- rate, will define total interest that borrower will pay
- ROI return of investment, will define total profit received by investors link to the generated agreement letter

design a RESTFful api that satisfy above requirement.

A system that handles:
 * loan creation
 * loan approval
 * investment
 * loan disbursement


## Analysis and Design
### Requirement Analysis
1. One way Loan state machine: proposed -> approved -> invested -> disbursed
2. each state has required data and validations
3. Investor and loan has many to many relationship

### Requirement Assumption
1. If the proposed loan is not eligible, reject Loan -> rejected
2. Once Loan approved, user/ Amartha cancels Loan -> canceled
3. Once Loan published to investor -> offered
4. Once Loan published but did not get invested -> canceled
5. positif flow of state proposed -> approved -> offered -> invested -> disbursed

state diagram:

### Data Model
Entity/ object:
1. Loan
    properties:
        - id (uuid, system generated)
        - borrower_id (uuid, input)
        - principal_amount (numeric, input)
        - interest_rate (numeric, input)
        - state (state enum: proposed , approved, offered, invested, disbursed, rejected, canceled - system generated)
        - created_at (timestamptz, system generated)
        - created_by (uuid, logged in user_id)

        - visit_proof_url (text, blob store url)
        - validator_id (uuid, input)
        - validated_at (timestamptz, input)

        - approved_at (timestamptz, system generated)
        - approved_by (uuid, logged in user_id)
        - rejected_at (timestamptz, system generated)
        - rejected_by (uuid, logged in user_id)
        - rejected_reason (text, input)

        - invested_amount
        - loan_agreement_letter_url
        - is_loan_aggrement_signed

        - disburse_by
        - disburse_at

        - publish_at
        - publish_by
    methods:
        - createLoan()
        - updateLoan()
        - getLoan()
        - approveLoan()
        - publishLoan()
        - disburseLoan()
    apis:
        POST /v1/loans
            requestBody:
                - borrower_id
                - principal_amount
                - interest_rate
                - visit_proof_url
                - validator_id
                - validated_at
            response:
                - 201 Created
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
            validations:
                - basic validation (empty, etc)
        PUT /v1/loans/{id}
            requestBody:
                - borrower_id
                - principal_amount
                - interest_rate
                - visit_proof_url
                - validator_id
                - validated_at
            response:
                - 200 Success
                - 404 Not Found
                - 400 Bad Request
                - 401 Unauthorized
                - 500 Internal Server Error
        GET /v1/loans/{id}
            response:
                - [all loan properties + investment]
        POST /v1/files
            - requestBody:
                - byte file
            response:
                - id
                - url

        POST /v1/loans/{id}/approvals
            requestBody:
                - rejected_reason
                - state: approve | reject
            response:
                -
            validation:
                - state must be proposed state
                - visit_proof_url, validator_id, and validated_at are not empty
        
        PATCH /v1/loan/{id}/publish
            requestBody:
                - state: offered
            response:
                -
            validation:
                -

2. Investment
    properties:
        - id
        - investor_id
        - loan_id
        - invested_amount
        - investment_agreement_letter_url
        - is_investment_aggrement_signed
        - total_profit
        - roi_percentage
        - signed_at
    methods:
        - createInvestment()
        - generateAgreement()
        - sendAgreement()
        - signAgreement()

3. Employee
    properties:
        - id
        - name
        - employee_number

4. Borrower
    properties
        - id
        - name
        - address
        - occupation
        - nik
        - dob

5. Investor
    properties:
        - id
        - name
        - nik
        - npwp


assumptions:
1. once approved it can not go back to proposed state. - Where can go if loan cancelled by borrower? I assume to cancelled state
2. the employee id of the field officer that hands the money and/or collect the agreement letter and the employee id of field validator - is one employee

to think:
1. once approved loan is ready to be offered to investors/lender
2. create state machine?
