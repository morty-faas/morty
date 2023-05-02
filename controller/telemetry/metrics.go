package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels = []string{"function", "isColdStart"}

	FunctionInvocationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "morty",
			Name:      "function_invocation",
			Help:      "The number of invocations per function",
		}, labels)

	FunctionInvocationDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "morty",
			Name:      "function_invocation_duration",
			Help:      "The invocation duration in seconds for a given function",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2},
		}, labels)
)

// init() will get called on main() function initialization, so
// we don't need to call it manually and we can setup everything
// required here to configure our metrics.
func init() {
	prometheus.MustRegister(FunctionInvocationCounter)
	prometheus.MustRegister(FunctionInvocationDurationHistogram)
}
