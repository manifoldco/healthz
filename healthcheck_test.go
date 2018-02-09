package healthz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
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
	t.Run("with multiple tests", func(t *testing.T) {
		defer resetTests()

		executedTests := 0
		var l sync.Mutex
		tstFuncFail := func(_ context.Context) (Status, error) {
			l.Lock()
			defer l.Unlock()
			executedTests++
			return Unavailable, errors.New("unavailable")
		}
		tstFuncSuccess := func(_ context.Context) (Status, error) {
			l.Lock()
			defer l.Unlock()
			executedTests++
			return Available, nil
		}
		RegisterTest("custom-failure1", tstFuncFail)
		RegisterTest("custom-success1", tstFuncSuccess)
		RegisterTest("custom-failure2", tstFuncFail)
		RegisterTest("custom-success2", tstFuncSuccess)

		hc, _, err := getHealth()
		if err != nil {
			t.Fatalf("Expected no error, got '%s'", err.Error())
		}

		// account for default test
		if ln := len(hc.Tests); ln != 5 {
			t.Fatalf("Expected '%d' tests, got '%d'", 5, ln)
		}

		if executedTests != 4 {
			t.Fatalf("Expected '%d' tests to be executed, got '%d'", 4, executedTests)
		}
	})

	t.Run("with an exceeded deadline", func(t *testing.T) {
		defer resetTests()

		Timeout = 50 * time.Millisecond
		tstFunc := func(_ context.Context) (Status, error) {
			time.Sleep(time.Second)
			return Available, nil
		}
		RegisterTest("success", tstFunc)

		hc, _, err := getHealth()
		if err != nil {
			t.Fatalf("Expected no error, got '%s'", err.Error())
		}

		// our registered test shouldn't be skipped
		if ln := len(hc.Tests); ln != 2 {
			t.Fatalf("Expected '%d' tests, got '%d'", 1, ln)
		}

		// our registered test should be marked unavailable despite what the
		// outcome is as the response will be too late
		if hc.Tests["success"].Status != Unavailable {
			t.Fatalf("Expected 'success' test to be Unavailable, got '%s'", hc.Tests["success"].Status)
		}

		if hc.Tests["success"].Error != ErrTimeout {
			t.Fatalf("Expected 'success' test to be timeout")
		}
	})

	t.Run("with a failing test", func(t *testing.T) {
		defer resetTests()

		tstFunc := func(_ context.Context) (Status, error) {
			return Unavailable, errors.New("unavailable")
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
		if hc.Status != Unavailable {
			t.Fatalf("Expected result to equal '%s', got '%s'", Unavailable, hc.Status)
		}
		if d := string(hc.Tests["custom-failure"].Error); d != string(Unavailable) {
			t.Fatalf("Expected custom-failure to equal '%s', got '%s'", string(Unavailable), d)
		}
	})

	t.Run("with a degraded test", func(t *testing.T) {
		defer resetTests()

		tstFunc := func(_ context.Context) (Status, error) {
			return Degraded, errors.New("degraded")
		}
		RegisterTest("custom-failure", tstFunc)

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
		if hc.Status != Degraded {
			t.Fatalf("Expected result to equal '%s', got '%s'", Degraded, hc.Status)
		}
		if d := string(hc.Tests["custom-failure"].Error); d != string(Degraded) {
			t.Fatalf("Expected custom-failure to equal '%s', got '%s'", string(Degraded), d)
		}
	})

	t.Run("with a degraded and unavailable test", func(t *testing.T) {
		defer resetTests()

		RegisterTest("degraded", func(_ context.Context) (Status, error) {
			return Degraded, errors.New("degraded")
		})
		RegisterTest("unavailable", func(_ context.Context) (Status, error) {
			return Unavailable, errors.New("unavailable")
		})

		hc, sc, err := getHealth()

		if err != nil {
			t.Fatalf("Expected no error, got '%s'", err.Error())
		}
		if sc != http.StatusServiceUnavailable {
			t.Fatalf("Expected status code to equal '%d', got '%d'", http.StatusServiceUnavailable, sc)
		}
		if ln := len(hc.Tests); ln != 3 {
			t.Fatalf("Expected '%d' tests, got '%d'", 3, ln)
		}
		if hc.Status != Unavailable {
			t.Fatalf("Expected result to equal '%s', got '%s'", Unavailable, hc.Status)
		}
		if d := string(hc.Tests["degraded"].Error); d != string(Degraded) {
			t.Fatalf("Expected degraded test to equal '%s', got '%s'", string(Degraded), d)
		}
		if d := string(hc.Tests["unavailable"].Error); d != string(Unavailable) {
			t.Fatalf("Expected unavailable test to equal '%s', got '%s'", string(Unavailable), d)
		}
	})

	t.Run("with a passing test", func(t *testing.T) {
		defer resetTests()

		tstFunc := func(_ context.Context) (Status, error) {
			return Available, nil
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
		if hc.Status != Available {
			t.Fatalf("Expected result to equal '%s', got '%s'", Available, hc.Status)
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
	Timeout = 5 * time.Second
	RegisterTest("default", defaultCheck)
}
