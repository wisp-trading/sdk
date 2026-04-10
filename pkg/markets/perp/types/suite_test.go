package types_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPerpTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Perp Types Suite")
}
