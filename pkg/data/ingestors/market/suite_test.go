package market

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMarketCoordinator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Market Data Coordinator Suite")
}
