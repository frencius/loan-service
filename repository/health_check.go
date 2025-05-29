package repository

import (
	"database/sql"
	"log"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
)

type IHealthCheckRepository interface {
	Ping() (ping *model.Ping, err error)
}

type HealthCheckRepository struct {
	DB *sql.DB
}

func NewHealthCheckRepository(app *application.App) IHealthCheckRepository {
	return &HealthCheckRepository{
		DB: app.DB,
	}
}

func (hcr *HealthCheckRepository) Ping() (ping *model.Ping, err error) {
	if err = hcr.DB.Ping(); err != nil {
		log.Println("Ping DB ERROR", err)
		return
	}

	ping = &model.Ping{}
	ping.DatabaseStatus = "OK"
	ping.ServiceStatus = "OK"

	return

}
