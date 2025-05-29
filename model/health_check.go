package model

type HealthCheckResponse struct {
	ServiceStatus  string `json:"service_status,omitempty"`
	DatabaseStatus string `json:"database_status,omitempty"`
}

type Ping struct {
	ServiceStatus  string
	DatabaseStatus string
}
