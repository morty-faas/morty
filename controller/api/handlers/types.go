package handlers

import "time"

type CustomHttpHeader string

const (
	// Function-Response-Time
	// ref: https://github.com/morty-faas/morty/issues/25
	HttpHeaderFunctionResponseTime CustomHttpHeader = "Function-Response-Time"
	// Orchestrate-Response-Time
	// ref: https://github.com/morty-faas/morty/issues/25
	HttpHeaderOrchestrateResponseTime CustomHttpHeader = "Orchestrate-Response-Time"
	// Morty-Response-Time
	// ref: https://github.com/morty-faas/morty/issues/25
	HttpHeaderMortyResponseTime CustomHttpHeader = "Morty-Response-Time"
)

type APIError struct {
	Message string `json:"message"`
}

func makeApiError(err error) *APIError {
	return &APIError{
		Message: err.Error(),
	}
}

// InvokeAnalytics contains few information about the function invocation, it can be used for debugging purposes
// and for benchmarking
type InvokeAnalytics struct {
	start time.Time
	// Total is the total time spent in the function handler, from the moment the request is received until the
	// response is sent back to the client
	total int64
	// Orchestrate is the time spent in the orchestration layer, from the moment the request is received until the
	// function is deployed and ready to be invoked
	orchestrate int64
	// Invoke is the time spent in the function handler, from the moment the function is invoked until the response
	// is sent back to the client
	invoke int64
}

func NewInvokeAnalytics() *InvokeAnalytics {
	return &InvokeAnalytics{
		start: time.Now(),
	}
}

// UpsertOrchestrate returns the time spent in the orchestration layer, in milliseconds. If the value is not set, it will
// be calculated else it will return a cached value.
func (i *InvokeAnalytics) UpsertOrchestrate() int64 {
	if i.orchestrate != 0 {
		return i.orchestrate
	}
	i.orchestrate = time.Since(i.start).Milliseconds()
	return i.orchestrate
}

// UpsertInvoke returns the time spent in for the invoke, in milliseconds. If the value is not set, it will be calculated
// else it will return a cached value.
func (i *InvokeAnalytics) UpsertInvoke() int64 {
	if i.invoke != 0 {
		return i.invoke
	}
	i.invoke = time.Since(i.start).Milliseconds() - i.UpsertOrchestrate()
	return i.invoke
}

// UpsertTotal returns the total time spent in the function handler, in milliseconds. If it's the first call to this
// method, it will register the time spent since start and return it, otherwise it will return the cached value.
func (i *InvokeAnalytics) UpsertTotal() int64 {
	if i.total != 0 {
		return i.total
	}
	i.total = time.Since(i.start).Milliseconds()
	return i.total
}
