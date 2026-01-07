package profiling_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProfiling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Profiling Suite")
}
