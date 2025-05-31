package model

import "time"

type Borrower struct {
	ID         string
	Name       string
	Address    string
	Occupation string
	NIK        string
	DOB        *time.Time
}
