package healthz

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Prefix represents the prefix for the health check endpoint.
var Prefix = ""

// Endpoint represents the endpoint we'll run the health check endpoint on
var Endpoint = "/_healthz"

var healthCheckTests = map[string]TestFunc{}

type middlewareFunc func(http.Handler) http.Handler

// TestFunc represents a function which will be executed when we run the health
// check endpoint.
type TestFunc func(context.Context) error

// Error represents a health check error
type Error string

// Error returns the error message of our error type.
func (e Error) Error() string {
	return string(e)
}

// Result represent the state of a TestFunc
type Result string

var (
	// Success represents the success result state
	Success Result = "success"

	// Failure represents the failure result state
	Failure Result = "failure"
)

// HealthCheck represents the overal health check status of the health check
// request.
type HealthCheck struct {
	CheckedAt  time.Time       `json:"checked_at"`
	DurationMs time.Duration   `json:"duration_ms"`
	Result     Result          `json:"result"`
	Tests      map[string]Test `json:"tests"`
}

// Test represents a single health check test. All the tests combined
// form the actual HealthCheck.
type Test struct {
	DurationMs time.Duration `json:"duration_ms"`
	Result     Result        `json:"result"`
	Error      Error         `json:"error,omitempty"`
}

// NewHandler wraps the given http handler with a /_healthz endpoint.
func NewHandler(dh http.Handler) http.Handler {
	return NewHandlerWithMiddleware(dh)
}

// NewHandlerWithMiddleware wraps the given handler with a new health endpoint.
// This health endpoint will be wrapped in the provided middleware.
func NewHandlerWithMiddleware(dh http.Handler, mw ...middlewareFunc) http.Handler {
	var handler http.Handler
	h := http.NewServeMux()

	handler = http.HandlerFunc(healthHandler)
	for _, mwh := range mw {
		handler = mwh(handler)
	}

	h.Handle(Prefix+Endpoint, handler)
	h.Handle("/", dh)

	return h
}

// RegisterTest adds a test to the HealthCheck handler. If a tests with the
// given name is already registered, this will panic.
func RegisterTest(name string, test TestFunc) {
	if _, ok := healthCheckTests[name]; ok {
		panic("Test already registered")
	}

	healthCheckTests[name] = test
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()

	hc := HealthCheck{
		CheckedAt: time.Now(),
		Tests:     map[string]Test{},
		Result:    Success,
	}

	success := true
	ctx := r.Context()
	for name, test := range healthCheckTests {
		hct := Test{
			Result: Success,
		}

		tStart := time.Now()

		if err := test(ctx); err != nil {
			success = false
			hct.Result = Failure
			hct.Error = Error(err.Error())
		}

		hct.DurationMs = time.Since(tStart) / time.Millisecond
		hc.Tests[name] = hct
	}

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		hc.Result = Failure
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	hc.DurationMs = time.Since(start) / time.Millisecond
	if err := json.NewEncoder(w).Encode(hc); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func defaultCheck(ctx context.Context) error {
	return nil
}

func init() {
	RegisterTest("default", defaultCheck)
}
