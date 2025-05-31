package service_test

import (
	"errors"

	"github.com/frencius/loan-service/mock"
	"github.com/frencius/loan-service/model"
	. "github.com/frencius/loan-service/service"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HealthCheckService", func() {
	var (
		mockRepo *mock.MockHealthCheckRepository
		service  *HealthCheckService
	)

	BeforeEach(func() {
		mockRepo = &mock.MockHealthCheckRepository{}
		service = &HealthCheckService{
			HealthCheckRepository: mockRepo,
		}
	})

	Describe("Ping", func() {
		It("should return health check response when repository returns successfully", func() {
			mockRepo.PingFunc = func() (*model.Ping, error) {
				return &model.Ping{
					DatabaseStatus: "OK",
					ServiceStatus:  "OK",
				}, nil
			}

			resp, err := service.Ping()
			Expect(err).To(BeNil())
			Expect(resp).NotTo(BeNil())
			Expect(resp.DatabaseStatus).To(Equal("OK"))
			Expect(resp.ServiceStatus).To(Equal("OK"))
		})

		It("should return error when repository returns error", func() {
			mockRepo.PingFunc = func() (*model.Ping, error) {
				return nil, errors.New("database down")
			}

			resp, err := service.Ping()
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
