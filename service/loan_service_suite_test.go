package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLoanService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoanService Suite")
}
