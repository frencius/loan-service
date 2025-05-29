package controller

import (
	"net/http"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/model"
	"github.com/frencius/loan-service/service"
)

type IHealthCheckController interface {
	Ping(w http.ResponseWriter, r *http.Request)
}
type HealthCheckController struct {
	HealthCheckService service.IHealthCheckService
}

func NewHealthCheckController(app *application.App) IHealthCheckController {
	return &HealthCheckController{
		HealthCheckService: service.NewHealthCheckService(app),
	}
}

func (hcc *HealthCheckController) Ping(w http.ResponseWriter, r *http.Request) {
	resp, err := hcc.HealthCheckService.Ping()
	if err != nil {
		errMsg, respCode := getErrorResponse(err)
		result := model.ComposeErrorResponse(respCode, err.Error(), errMsg)
		WriteHTTPResponse(w, respCode, result)
		return
	}

	respCode := http.StatusOK
	result := model.ComposeResponse(resp, respCode)
	WriteHTTPResponse(w, respCode, result)
}
