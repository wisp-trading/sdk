package indicators_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIndicators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Indicators Suite")
}
