package healthz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthChecks_Default(t *testing.T) {
	hc, sc, err := getHealth()

	if err != nil {
		t.Fatalf("Expected no error, got '%s'", err.Error())
	}

	if sc != http.StatusOK {
		t.Fatalf("Expected status code to equal '%d', got '%d'", http.StatusOK, sc)
	}

	if ln := len(hc.Tests); ln != 1 {
		t.Fatalf("Expected '%d' tests, got '%d'", 1, ln)
	}
}

func TestHealthChecks_Custom(t *testing.T) {
	t.Run("with a failing test", func(t *testing.T) {
		defer resetTests()

		tstFunc := func(_ context.Context) error {
			return errors.New("failure")
		}
		RegisterTest("custom-failure", tstFunc)

		hc, sc, err := getHealth()

		if err != nil {
			t.Fatalf("Expected no error, got '%s'", err.Error())
		}
		if sc != http.StatusServiceUnavailable {
			t.Fatalf("Expected status code to equal '%d', got '%d'", http.StatusServiceUnavailable, sc)
		}
		if ln := len(hc.Tests); ln != 2 {
			t.Fatalf("Expected '%d' tests, got '%d'", 2, ln)
		}
		if hc.Result != Failure {
			t.Fatalf("Expected result to equal '%s', got '%s'", Failure, hc.Result)
		}
		if d := string(hc.Tests["custom-failure"].Error); d != string(Failure) {
			t.Fatalf("Expected custom-failure to equal '%s', got '%s'", string(Failure), d)
		}
	})

	t.Run("with a passing test", func(t *testing.T) {
		defer resetTests()

		tstFunc := func(_ context.Context) error {
			return nil
		}
		RegisterTest("custom-success", tstFunc)

		hc, sc, err := getHealth()
		if err != nil {
			t.Fatalf("Expected no error, got '%s'", err.Error())
		}
		if sc != http.StatusOK {
			t.Fatalf("Expected status code to equal '%d', got '%d'", http.StatusOK, sc)
		}
		if ln := len(hc.Tests); ln != 2 {
			t.Fatalf("Expected '%d' tests, got '%d'", 2, ln)
		}
		if hc.Result != Success {
			t.Fatalf("Expected result to equal '%s', got '%s'", Success, hc.Result)
		}
	})
}

func getHealth() (HealthCheck, int, error) {
	hdlr := NewHandler(http.NewServeMux())
	srv := httptest.NewServer(hdlr)
	defer srv.Close()

	hc := HealthCheck{}

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/_healthz", nil)
	if err != nil {
		return hc, 0, err
	}
	cl := &http.Client{}
	rsp, err := cl.Do(req)
	if err != nil {
		return hc, 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 200 || rsp.StatusCode == 503 {
		err := json.NewDecoder(rsp.Body).Decode(&hc)
		return hc, rsp.StatusCode, err
	}

	return hc, 0, fmt.Errorf("Unexpected status code: %d", rsp.StatusCode)
}

func resetTests() {
	healthCheckTests = map[string]TestFunc{}
	RegisterTest("default", defaultCheck)
}
