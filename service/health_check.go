package service

import (
	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/repository"
)

type IHealthCheckService interface {
	Ping() (healthCheckResponse *model.HealthCheckResponse, err error)
}

type HealthCheckService struct {
	HealthCheckRepository repository.IHealthCheckRepository
}

func NewHealthCheckService(app *application.App) IHealthCheckService {
	return &HealthCheckService{
		HealthCheckRepository: repository.NewHealthCheckRepository(app),
	}
}

func (hcs *HealthCheckService) Ping() (healthCheckResponse *model.HealthCheckResponse, err error) {
	ping, err := hcs.HealthCheckRepository.Ping()
	if err != nil {
		return
	}

	healthCheckResponse = &model.HealthCheckResponse{
		DatabaseStatus: ping.DatabaseStatus,
		ServiceStatus:  ping.ServiceStatus,
	}

	return
}
