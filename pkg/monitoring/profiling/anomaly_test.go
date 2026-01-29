package profiling_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/monitoring/profiling"
	profiling2 "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
)

var _ = Describe("AnomalyDetector", func() {
	var (
		detector          profiling2.AnomalyDetector
		strategyName      string
		warningThreshold  float64
		criticalThreshold float64
		windowSize        int
	)

	BeforeEach(func() {
		strategyName = "test-strategy"
		warningThreshold = 1.5  // 150% of baseline
		criticalThreshold = 2.0 // 200% of baseline
		windowSize = 100
		detector = profiling.NewAnomalyDetector(warningThreshold, criticalThreshold, windowSize)
	})

	Describe("CheckExecution", func() {
		Context("when no baseline exists", func() {
			It("should return OK status", func() {
				alert := detector.CheckExecution(strategyName, 100*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.OK))
				Expect(alert.Message).To(ContainSubstring("No baseline"))
			})
		})

		Context("when execution is within normal range", func() {
			It("should return OK status", func() {
				// Establish baseline
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				// Check execution within range (110ms is 1.1x baseline)
				alert := detector.CheckExecution(strategyName, 110*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.OK))
				Expect(alert.Message).To(ContainSubstring("normal range"))
			})
		})

		Context("when execution exceeds warning threshold", func() {
			It("should return Warning status", func() {
				// Establish baseline of 100ms
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				// Check execution at 1.6x baseline (160ms) - above warning (1.5x) but below critical (2.0x)
				alert := detector.CheckExecution(strategyName, 160*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.Warning))
				Expect(alert.Message).To(ContainSubstring("WARNING"))
				Expect(alert.Message).To(ContainSubstring("1.6x"))
			})
		})

		Context("when execution exceeds critical threshold", func() {
			It("should return Critical status", func() {
				// Establish baseline of 100ms
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				// Check execution at 2.5x baseline (250ms) - above critical (2.0x)
				alert := detector.CheckExecution(strategyName, 250*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.Critical))
				Expect(alert.Message).To(ContainSubstring("CRITICAL"))
				Expect(alert.Message).To(ContainSubstring("2.5x"))
			})
		})

		Context("when execution is exactly at warning threshold", func() {
			It("should return Warning status", func() {
				// Establish baseline of 100ms
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				// Check execution at exactly 1.5x baseline (150ms)
				alert := detector.CheckExecution(strategyName, 150*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.Warning))
			})
		})

		Context("when execution is exactly at critical threshold", func() {
			It("should return Critical status", func() {
				// Establish baseline of 100ms
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				// Check execution at exactly 2.0x baseline (200ms)
				alert := detector.CheckExecution(strategyName, 200*time.Millisecond)

				Expect(alert.Severity).To(Equal(profiling2.Critical))
			})
		})
	})

	Describe("UpdateBaseline", func() {
		Context("when creating initial baseline", func() {
			It("should set baseline to first measurement", func() {
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)

				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(100 * time.Millisecond))
			})
		})

		Context("when building up to window size", func() {
			It("should use simple average", func() {
				// Add measurements: 100ms, 200ms, 300ms
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)
				detector.UpdateBaseline(strategyName, 200*time.Millisecond)
				detector.UpdateBaseline(strategyName, 300*time.Millisecond)

				// Average should be (100 + 200 + 300) / 3 = 200ms
				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(200 * time.Millisecond))
			})
		})

		Context("when window is full", func() {
			It("should use exponential moving average", func() {
				// Fill window with 100ms measurements
				for i := 0; i < 100; i++ {
					detector.UpdateBaseline(strategyName, 100*time.Millisecond)
				}

				baselineBefore := detector.GetBaseline(strategyName)
				Expect(baselineBefore).To(Equal(100 * time.Millisecond))

				// Add a new measurement of 200ms
				detector.UpdateBaseline(strategyName, 200*time.Millisecond)

				// Baseline should have moved slightly toward 200ms (EMA)
				baselineAfter := detector.GetBaseline(strategyName)
				Expect(baselineAfter).To(BeNumerically(">", 100*time.Millisecond))
				Expect(baselineAfter).To(BeNumerically("<", 200*time.Millisecond))
			})
		})

		Context("when updating multiple strategies", func() {
			It("should maintain separate baselines", func() {
				strategy1 := "strategy-1"
				strategy2 := "strategy-2"

				detector.UpdateBaseline(strategy1, 100*time.Millisecond)
				detector.UpdateBaseline(strategy2, 200*time.Millisecond)

				baseline1 := detector.GetBaseline(strategy1)
				baseline2 := detector.GetBaseline(strategy2)

				Expect(baseline1).To(Equal(100 * time.Millisecond))
				Expect(baseline2).To(Equal(200 * time.Millisecond))
			})
		})
	})

	Describe("GetBaseline", func() {
		Context("when baseline exists", func() {
			It("should return the baseline duration", func() {
				detector.UpdateBaseline(strategyName, 150*time.Millisecond)

				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(150 * time.Millisecond))
			})
		})

		Context("when baseline does not exist", func() {
			It("should return zero duration", func() {
				baseline := detector.GetBaseline("nonexistent-strategy")
				Expect(baseline).To(Equal(time.Duration(0)))
			})
		})
	})

	Describe("Reset", func() {
		Context("when resetting existing baseline", func() {
			It("should clear the baseline", func() {
				// Create baseline
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)
				Expect(detector.GetBaseline(strategyName)).To(Equal(100 * time.Millisecond))

				// Reset
				detector.Reset(strategyName)

				// Baseline should be gone
				Expect(detector.GetBaseline(strategyName)).To(Equal(time.Duration(0)))
			})
		})

		Context("when resetting non-existent baseline", func() {
			It("should not panic", func() {
				Expect(func() {
					detector.Reset("nonexistent-strategy")
				}).NotTo(Panic())
			})
		})

		Context("after reset", func() {
			It("should allow creating new baseline", func() {
				// Create and reset
				detector.UpdateBaseline(strategyName, 100*time.Millisecond)
				detector.Reset(strategyName)

				// Create new baseline
				detector.UpdateBaseline(strategyName, 200*time.Millisecond)

				// New baseline should be established
				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(200 * time.Millisecond))
			})
		})
	})

	Describe("Concurrency", func() {
		It("should handle concurrent updates safely", func() {
			done := make(chan bool)

			// Multiple goroutines updating baseline
			for i := 0; i < 10; i++ {
				go func(duration time.Duration) {
					defer GinkgoRecover()
					for j := 0; j < 100; j++ {
						detector.UpdateBaseline(strategyName, duration)
					}
					done <- true
				}(time.Duration(i+1) * 10 * time.Millisecond)
			}

			// Wait for all goroutines
			for i := 0; i < 10; i++ {
				<-done
			}

			// Should not panic and should have a baseline
			baseline := detector.GetBaseline(strategyName)
			Expect(baseline).To(BeNumerically(">", 0))
		})

		It("should handle concurrent reads and writes", func() {
			done := make(chan bool)

			// Writer goroutines
			for i := 0; i < 5; i++ {
				go func(duration time.Duration) {
					defer GinkgoRecover()
					for j := 0; j < 50; j++ {
						detector.UpdateBaseline(strategyName, duration)
					}
					done <- true
				}(100 * time.Millisecond)
			}

			// Reader goroutines
			for i := 0; i < 5; i++ {
				go func() {
					defer GinkgoRecover()
					for j := 0; j < 50; j++ {
						detector.GetBaseline(strategyName)
						detector.CheckExecution(strategyName, 100*time.Millisecond)
					}
					done <- true
				}()
			}

			// Wait for all goroutines
			for i := 0; i < 10; i++ {
				<-done
			}

			// Should not panic
			Expect(detector.GetBaseline(strategyName)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Edge Cases", func() {
		Context("with very small durations", func() {
			It("should handle microsecond precision", func() {
				detector.UpdateBaseline(strategyName, 100*time.Microsecond)

				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(100 * time.Microsecond))

				// 110μs is 1.1x baseline - should be OK
				alert := detector.CheckExecution(strategyName, 110*time.Microsecond)
				Expect(alert.Severity).To(Equal(profiling2.OK))
			})
		})

		Context("with very large durations", func() {
			It("should handle second-level durations", func() {
				detector.UpdateBaseline(strategyName, 10*time.Second)

				baseline := detector.GetBaseline(strategyName)
				Expect(baseline).To(Equal(10 * time.Second))

				// 25 seconds is 2.5x baseline - should be critical
				alert := detector.CheckExecution(strategyName, 25*time.Second)
				Expect(alert.Severity).To(Equal(profiling2.Critical))
			})
		})

		Context("with zero duration baseline", func() {
			It("should handle division by zero", func() {
				detector.UpdateBaseline(strategyName, 0)

				// This should not panic (division by zero protection)
				Expect(func() {
					detector.CheckExecution(strategyName, 100*time.Millisecond)
				}).NotTo(Panic())
			})
		})
	})
})
