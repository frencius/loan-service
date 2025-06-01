package model

import "time"

type Investor struct {
	ID          string
	Name        string
	NIK         string
	NPWP        string
	Email       string
	PhoneNumber string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
