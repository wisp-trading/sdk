package signal_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSpotSignal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spot Signal Suite")
}
