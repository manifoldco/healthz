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

// MiddlewareFunc represents a function that acts as middleware.
type MiddlewareFunc func(http.Handler) http.Handler

// TestFunc represents a function which will be executed when we run the health
// check endpoint.
type TestFunc func(context.Context) (Status, error)

// Error represents a health check error
type Error string

// Error returns the error message of our error type.
func (e Error) Error() string {
	return string(e)
}

// Status represents the state of a TestFunc
type Status string

var (
	// Available represents the success result state
	Available Status = "available"

	// Degraded represents a degraded result state
	Degraded Status = "degraded"

	// Unavailable represents the failure result state
	Unavailable Status = "unavailable"
)

// HealthCheck represents the overal health check status of the health check
// request.
type HealthCheck struct {
	CheckedAt  time.Time       `json:"checked_at"`
	DurationMs time.Duration   `json:"duration_ms"`
	Status     Status          `json:"status"`
	Tests      map[string]Test `json:"tests"`
}

// Test represents a single health check test. All the tests combined
// form the actual HealthCheck.
type Test struct {
	DurationMs time.Duration `json:"duration_ms"`
	Status     Status        `json:"status"`
	Error      Error         `json:"error,omitempty"`
}

// NewHandler wraps the given http handler with a /_healthz endpoint.
func NewHandler(dh http.Handler) http.Handler {
	return NewHandlerWithMiddleware(dh)
}

// NewHandlerWithMiddleware wraps the given handler with a new health endpoint.
// This health endpoint will be wrapped in the provided middleware.
func NewHandlerWithMiddleware(dh http.Handler, mw ...MiddlewareFunc) http.Handler {
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
		Status:    Available,
	}

	statuses := []Status{}
	ctx := r.Context()
	for name, test := range healthCheckTests {
		hct := Test{
			Status: Available,
		}

		tStart := time.Now()
		testStatus, err := test(ctx)
		if err != nil {
			hct.Status = testStatus
			hct.Error = Error(err.Error())
		}

		statuses = append(statuses, testStatus)
		hct.DurationMs = time.Since(tStart) / time.Millisecond
		hc.Tests[name] = hct
	}

	hc.Status = getOverallStatus(statuses)
	switch hc.Status {
	case Unavailable:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusOK)
	}

	hc.DurationMs = time.Since(start) / time.Millisecond
	if err := json.NewEncoder(w).Encode(hc); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getOverallStatus(statuses []Status) Status {
	status := Available
	for _, s := range statuses {
		if s == Unavailable {
			return s
		}

		if s == Degraded {
			status = Degraded
		}
	}

	return status
}

func defaultCheck(ctx context.Context) (Status, error) {
	return Available, nil
}

func init() {
	RegisterTest("default", defaultCheck)
}
