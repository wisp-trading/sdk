package batch

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spot Batch Ingestor Suite")
}
