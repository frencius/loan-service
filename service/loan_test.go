package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/frencius/loan-service/mock"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/service"

	"github.com/golang/mock/gomock"
)

var _ = Describe("LoanService", func() {
	var (
		mockCtrl         *gomock.Controller
		mockLoanRepo     *mock.MockILoanRepository
		mockBorrowerRepo *mock.MockIBorrowerRepository
		loanSvc          service.ILoanService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockLoanRepo = mock.NewMockILoanRepository(mockCtrl)
		mockBorrowerRepo = mock.NewMockIBorrowerRepository(mockCtrl)

		loanSvc = &service.LoanService{
			LoanRepository:     mockLoanRepo,
			BorrowerRepository: mockBorrowerRepo,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("CreateLoan", func() {
		It("should create a loan and return loan id", func() {
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

			mockLoanRepo.EXPECT().
				CreateLoan(ctx, gomock.Any()).
				Return("123", nil)

			resp, err := loanSvc.CreateLoan(ctx, createReq)

			Expect(err).To(BeNil())
			Expect(resp).NotTo(BeNil())
			Expect(resp.LoanID).To(Equal("123"))
			Expect(resp.State).To(Equal("proposed"))
		})

		It("should return error if borrower not found", func() {
			ctx := context.Background()
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

			mockLoanRepo.EXPECT().
				CreateLoan(ctx, gomock.Any()).
				Return("", errors.New("insert failed"))

			resp, err := loanSvc.CreateLoan(ctx, createReq)

			Expect(err).To(MatchError("insert failed"))
			Expect(resp).To(BeNil())
		})
	})
})
