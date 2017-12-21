package healthz

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

// CacheMiddleware allows the health checks to be cached for a specified
// duration. This means that subsequent calls will return the cached response
// for all checks within a timespan of the given time.Duration.
// When an error occurs writing the cached data to the response, this will
// panic.
func CacheMiddleware(d time.Duration) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		h := &cacheHandler{
			duration: d,
			lastTick: time.Now(),
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if h.shouldCapture() {
				h.captureResponse(next, r)
			}

			h.writeCachedResponse(w)
		})
	}
}

type cacheHandler struct {
	sync.RWMutex

	duration time.Duration
	lastTick time.Time
	writer   *cacheWriter
}

func (h *cacheHandler) captureResponse(next http.Handler, r *http.Request) {
	h.Lock()
	defer h.Unlock()

	h.writer = &cacheWriter{
		Buffer: bytes.NewBuffer([]byte{}),
		header: http.Header{},
	}
	h.lastTick = time.Now()

	next.ServeHTTP(h.writer, r)
}

func (h *cacheHandler) writeCachedResponse(w http.ResponseWriter) {
	h.RLock()
	defer h.RUnlock()

	for hk, hvs := range h.writer.header {
		for _, hv := range hvs {
			w.Header().Add(hk, hv)
		}
	}

	if h.writer.statusCodeSet {
		w.WriteHeader(h.writer.statusCode)
	}

	if _, err := w.Write(h.writer.Bytes()); err != nil {
		panic(err)
	}
}

func (h *cacheHandler) shouldCapture() bool {
	h.RLock()
	defer h.RUnlock()

	if h.writer == nil {
		return true
	}

	return time.Since(h.lastTick) >= h.duration
}

type cacheWriter struct {
	*bytes.Buffer
	header        http.Header
	statusCodeSet bool
	statusCode    int
}

func (h *cacheWriter) Header() http.Header {
	return h.header
}

func (h *cacheWriter) WriteHeader(i int) {
	h.statusCodeSet = true
	h.statusCode = i
}
