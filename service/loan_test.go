package service_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/frencius/loan-service/mock"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/service"

	"github.com/golang/mock/gomock"
)

var _ = Describe("LoanService", func() {
	var (
		mockCtrl           *gomock.Controller
		mockLoanRepo       *mock.MockILoanRepository
		mockBorrowerRepo   *mock.MockIBorrowerRepository
		mockInvestorRepo   *mock.MockIInvestorRepository
		mockInvestmentRepo *mock.MockIInvestmentRepository
		loanSvc            service.ILoanService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockLoanRepo = mock.NewMockILoanRepository(mockCtrl)
		mockBorrowerRepo = mock.NewMockIBorrowerRepository(mockCtrl)
		mockInvestorRepo = mock.NewMockIInvestorRepository(mockCtrl)
		mockInvestmentRepo = mock.NewMockIInvestmentRepository(mockCtrl)

		loanSvc = &service.LoanService{
			LoanRepository:       mockLoanRepo,
			BorrowerRepository:   mockBorrowerRepo,
			InvestorRepository:   mockInvestorRepo,
			InvestmentRepository: mockInvestmentRepo,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("CreateLoan", func() {
		It("should create a loan and return loan id", func() {
			ctx := context.WithValue(context.Background(), "userID", "user-1")
			borrowerID := "1"
			createReq := &model.CreateLoanRequest{
				BorrowerID:      borrowerID,
				PrincipalAmount: 1000000,
				InterestRate:    5.5,
				ROIRate:         2.0,
			}
			borrower := &model.Borrower{ID: borrowerID}

			mockBorrowerRepo.EXPECT().
				GetBorrowerByID(ctx, borrowerID).
				Return(borrower, nil)
			mockLoanRepo.EXPECT().
				CreateLoan(ctx, gomock.Any()).
				DoAndReturn(func(_ context.Context, loan *model.Loan) (string, error) {
					Expect(loan.BorrowerID).To(Equal(borrowerID))
					Expect(loan.State).To(Equal(model.LoanStateProposed))
					Expect(loan.CreatedBy).To(Equal("user-1"))
					return "123", nil
				})

			resp, err := loanSvc.CreateLoan(ctx, createReq)
			Expect(err).To(BeNil())
			Expect(resp).NotTo(BeNil())
			Expect(resp.LoanID).To(Equal("123"))
			Expect(resp.State).To(Equal("proposed"))
		})

		It("should return error if borrower not found", func() {
			ctx := context.WithValue(context.Background(), "userID", "user-1")
			borrowerID := "99"
			createReq := &model.CreateLoanRequest{
				BorrowerID:      borrowerID,
				PrincipalAmount: 1000000,
				InterestRate:    5.5,
				ROIRate:         2.0,
			}
			mockBorrowerRepo.EXPECT().
				GetBorrowerByID(ctx, borrowerID).
				Return(nil, errors.New("not found"))

			resp, err := loanSvc.CreateLoan(ctx, createReq)
			Expect(err).To(MatchError("not found"))
			Expect(resp).To(BeNil())
		})

		It("should return error if CreateLoan fails", func() {
			ctx := context.WithValue(context.Background(), "userID", "user-1")
			borrowerID := "1"
			createReq := &model.CreateLoanRequest{
				BorrowerID:      borrowerID,
				PrincipalAmount: 1000000,
				InterestRate:    5.5,
				ROIRate:         2.0,
			}
			borrower := &model.Borrower{ID: borrowerID}

			mockBorrowerRepo.EXPECT().
				GetBorrowerByID(ctx, borrowerID).
				Return(borrower, nil)
			mockLoanRepo.EXPECT().
				CreateLoan(ctx, gomock.Any()).
				Return("", errors.New("insert failed"))

			resp, err := loanSvc.CreateLoan(ctx, createReq)
			Expect(err).To(MatchError("insert failed"))
			Expect(resp).To(BeNil())
		})

		It("should panic if userID is missing from context", func() {
			ctx := context.Background()
			borrowerID := "1"
			createReq := &model.CreateLoanRequest{
				BorrowerID:      borrowerID,
				PrincipalAmount: 1000000,
				InterestRate:    5.5,
				ROIRate:         2.0,
			}
			borrower := &model.Borrower{ID: borrowerID}

			mockBorrowerRepo.EXPECT().
				GetBorrowerByID(ctx, borrowerID).
				Return(borrower, nil)

			Expect(func() {
				_, _ = loanSvc.CreateLoan(ctx, createReq)
			}).To(Panic())
		})

		It("should panic if userID in context is not a string", func() {
			ctx := context.WithValue(context.Background(), "userID", 12345)
			borrowerID := "1"
			createReq := &model.CreateLoanRequest{
				BorrowerID:      borrowerID,
				PrincipalAmount: 1000000,
				InterestRate:    5.5,
				ROIRate:         2.0,
			}
			borrower := &model.Borrower{ID: borrowerID}

			mockBorrowerRepo.EXPECT().
				GetBorrowerByID(ctx, borrowerID).
				Return(borrower, nil)

			Expect(func() {
				_, _ = loanSvc.CreateLoan(ctx, createReq)
			}).To(Panic())
		})
	})

	Context("UpdateLoanState", func() {
		It("should update loan state successfully", func() {
			ctx := context.Background()
			loanID := "loan-1"
			oldState := model.LoanStateProposed
			newState := model.LoanStateApproved
			now := time.Now()
			loan := &model.Loan{
				ID:            loanID,
				State:         oldState,
				VisitProofURL: "http://proof",
				ValidatedAt:   &now,
				ValidatedBy:   "admin",
			}
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  string(newState),
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockLoanRepo.EXPECT().
				UpdateLoanState(ctx, loan, newState).
				Return(nil)

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(BeNil())
			Expect(resp).To(BeNil())
		})

		It("should return error if loan not found", func() {
			ctx := context.Background()
			loanID := "loan-404"
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  string(model.LoanStateApproved),
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(nil, errors.New("not found"))

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(MatchError("not found"))
			Expect(resp).To(BeNil())
		})

		It("should return error if new state is invalid", func() {
			ctx := context.Background()
			loanID := "loan-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStateProposed,
			}
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  "invalid_state",
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(Equal(model.ErrorLoanStateInvalid))
			Expect(resp).To(BeNil())
		})

		It("should return error if state transition is not allowed", func() {
			ctx := context.Background()
			loanID := "loan-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStateProposed,
			}
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  string(model.LoanStateDisbursed),
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(Equal(model.ErrorLoanStateTransitionNotAllowed))
			Expect(resp).To(BeNil())
		})

		It("should return error if state requirements are not fulfilled", func() {
			ctx := context.Background()
			loanID := "loan-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStateProposed,
			}
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  string(model.LoanStateApproved),
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(Equal(model.ErrorStateTransitionRequirementNotFulfilled))
			Expect(resp).To(BeNil())
		})

		It("should return error if UpdateLoanState repo fails", func() {
			ctx := context.Background()
			loanID := "loan-1"
			oldState := model.LoanStateProposed
			newState := model.LoanStateApproved
			now := time.Now()
			loan := &model.Loan{
				ID:            loanID,
				State:         oldState,
				VisitProofURL: "http://proof",
				ValidatedAt:   &now,
				ValidatedBy:   "admin",
			}
			updateReq := &model.UpdateLoanStateRequest{
				LoanID: loanID,
				State:  string(newState),
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockLoanRepo.EXPECT().
				UpdateLoanState(ctx, loan, newState).
				Return(errors.New("update failed"))

			resp, err := loanSvc.UpdateLoanState(ctx, updateReq)
			Expect(err).To(MatchError("update failed"))
			Expect(resp).To(BeNil())
		})
	})

	Context("CreateLoanInvestment", func() {
		It("should return error if loan not found", func() {
			ctx := context.Background()
			loanID := "loan-404"
			investorID := "inv-1"
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(nil, errors.New("not found"))

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(MatchError("not found"))
			Expect(resp).To(BeNil())
		})

		It("should return error if loan state is not published", func() {
			ctx := context.Background()
			loanID := "loan-1"
			investorID := "inv-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStateProposed,
			}
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(Equal(model.ErrorStateMustBePublished))
			Expect(resp).To(BeNil())
		})

		It("should return error if investor not found", func() {
			ctx := context.Background()
			loanID := "loan-1"
			investorID := "inv-404"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStatePublished,
			}
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockInvestorRepo.EXPECT().
				GetInvestorByID(ctx, investorID).
				Return(nil, errors.New("not found"))

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(MatchError("not found"))
			Expect(resp).To(BeNil())
		})

		It("should return error if investment already exists", func() {
			ctx := context.Background()
			loanID := "loan-1"
			investorID := "inv-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStatePublished,
			}
			investor := &model.Investor{ID: investorID}
			existingInvestment := &model.Investment{LoanID: loanID, InvestorID: investorID}
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockInvestorRepo.EXPECT().
				GetInvestorByID(ctx, investorID).
				Return(investor, nil)
			mockInvestmentRepo.EXPECT().
				GetInvestmentByInvestorID(ctx, investorID).
				Return(existingInvestment, nil)

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(Equal(model.ErrorInvestmentExist))
			Expect(resp).To(BeNil())
		})

		It("should return error if CreateInvestment fails", func() {
			ctx := context.Background()
			loanID := "loan-1"
			investorID := "inv-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStatePublished,
			}
			investor := &model.Investor{ID: investorID}
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockInvestorRepo.EXPECT().
				GetInvestorByID(ctx, investorID).
				Return(investor, nil)
			mockInvestmentRepo.EXPECT().
				GetInvestmentByInvestorID(ctx, investorID).
				Return(nil, model.ErrorInvestmentNotFound)
			mockInvestmentRepo.EXPECT().
				CreateInvestment(ctx, gomock.Any()).
				Return("", errors.New("fail"))

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(MatchError("fail"))
			Expect(resp).To(BeNil())
		})

		It("should return error if UpdateLoanTotalInvestedAmount fails", func() {
			ctx := context.Background()
			loanID := "loan-1"
			investorID := "inv-1"
			loan := &model.Loan{
				ID:    loanID,
				State: model.LoanStatePublished,
			}
			investor := &model.Investor{ID: investorID}
			createReq := &model.CreateLoanInvestmentRequest{
				LoanID:           loanID,
				InvestorID:       investorID,
				InvestmentAmount: 100,
			}

			mockLoanRepo.EXPECT().
				GetLoanByID(ctx, loanID).
				Return(loan, nil)
			mockInvestorRepo.EXPECT().
				GetInvestorByID(ctx, investorID).
				Return(investor, nil)
			mockInvestmentRepo.EXPECT().
				GetInvestmentByInvestorID(ctx, investorID).
				Return(nil, model.ErrorInvestmentNotFound)
			mockInvestmentRepo.EXPECT().
				CreateInvestment(ctx, gomock.Any()).
				Return("invst-1", nil)
			mockLoanRepo.EXPECT().
				UpdateLoanTotalInvestedAmount(ctx, loan).
				Return(errors.New("fail"))

			resp, err := loanSvc.CreateLoanInvestment(ctx, createReq)
			Expect(err).To(MatchError("fail"))
			Expect(resp).To(BeNil())
		})
	})
})
