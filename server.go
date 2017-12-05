package healthz

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

// Logger represents the logger interface we expect our server to have.
type Logger interface {
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Server represents the HealthCheck server that runs the health check endpoint.
type Server struct {
	srv    *http.Server
	logger Logger
}

// NewServer will start a new server on the given port to start the health check
// endpoint.
// See `RunServerWithMiddleware` for more information.
func NewServer(host string, port int) *Server {
	return NewServerWithMiddleware(host, port)
}

// NewServerWithMiddleware will start a new server and attach the given
// middleware to the created `_healthz` endpoint. It is up to the user to call
// `Server.Shutdown` once the service is shutting down to terminate the server
// gracefully.
func NewServerWithMiddleware(host string, port int, mw ...middlewareFunc) *Server {
	handler := NewHandlerWithMiddleware(http.NewServeMux(), mw...)

	addr := host + ":" + strconv.Itoa(port)
	if addr == ":" {
		addr = ":80"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server{srv: srv, logger: &defaultLogger{}}
}

// RegisterLogger will register a logger for this server. If no logger is given,
// the default go logger will be used.
func (s *Server) RegisterLogger(l Logger) {
	s.logger = l
}

// Start starts the health check server. This returns an error when listening or
// serving the requests causes an error. Otherwise this will be blocking.
func (s *Server) Start() error {
	s.logger.Printf("Healthcheck listening on http://%s", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil {
		s.logger.Fatalf(err.Error())
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the health check server.
func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Printf("Shutting down health server")
	s.srv.Shutdown(ctx)
	s.logger.Printf("Health server shut down")
}

type defaultLogger struct{}

func (l *defaultLogger) Printf(format string, v ...interface{}) { log.Printf(format, v...) }
func (l *defaultLogger) Fatalf(format string, v ...interface{}) { log.Fatalf(format, v...) }
