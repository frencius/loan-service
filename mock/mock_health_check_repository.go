package mock

import (
	"github.com/frencius/loan-service/model"
)

type MockHealthCheckRepository struct {
	PingFunc func() (*model.Ping, error)
}

func (m *MockHealthCheckRepository) Ping() (*model.Ping, error) {
	return m.PingFunc()
}
