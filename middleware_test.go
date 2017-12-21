package healthz_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/manifoldco/healthz"
)

func BenchmarkCacheMiddleware(b *testing.B) {
	b.Run("with long cache", func(b *testing.B) {
		hdlr := healthz.NewHandlerWithMiddleware(
			http.NewServeMux(),
			healthz.CacheMiddleware(time.Minute),
		)
		srv := httptest.NewServer(hdlr)
		defer srv.Close()

		cl := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/_healthz", nil)
		if err != nil {
			b.Fatalf("Expected no error, got '%s'", err.Error())
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			rsp, err := cl.Do(req)
			if err != nil {
				b.Fatalf("Expected no error, got '%s'", err.Error())
			}
			defer rsp.Body.Close()
		}
	})

	b.Run("with short cache", func(b *testing.B) {
		hdlr := healthz.NewHandlerWithMiddleware(
			http.NewServeMux(),
			healthz.CacheMiddleware(time.Millisecond),
		)
		srv := httptest.NewServer(hdlr)
		defer srv.Close()

		cl := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/_healthz", nil)
		if err != nil {
			b.Fatalf("Expected no error, got '%s'", err.Error())
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			rsp, err := cl.Do(req)
			if err != nil {
				b.Fatalf("Expected no error, got '%s'", err.Error())
			}
			defer rsp.Body.Close()
		}
	})
}
