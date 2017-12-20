# Healthz

[Code of Conduct](./.github/CODE_OF_CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/healthz.svg?label=latest)](https://github.com/manifoldco/healthz/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/healthz)
[![Travis](https://img.shields.io/travis/manifoldco/healthz/master.svg)](https://travis-ci.org/manifoldco/healthz)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/healthz)](https://goreportcard.com/report/github.com/manifoldco/healthz)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

This is a package that allows you to set up health checks for your services. By
default, the health checks will be available at the `/_healthz` endpoint but can
be configured with the `Prefix` and `Endpoint` variables.

By default, the health check package will add a single default test. This test
doesn't validate anything, it simply returns no error.

## Registering additional tests

If you want to add more tests in case your service/worker depends on a specific
set of tools (like a JWT key mounted on a volume), you can register a new test
as follows:

```go
func init() {
	healthz.RegisterTest("jwt-key", jwtCheck)
}

func jwtCheck(ctx context.Context) error {
	_, err := os.Stat("/jwt/ecdsa-private.pem")
	return err
}
```

## Attaching the endpoints to an existing server

If you have an existing server running, you can attach the `/_healthz` endpoint
to it without having to start a separate server.

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := healthz.NewHandler(mux)
	http.ListenAndServe(":3000", handler)
}
```

This will create a new mux which listens to the `/hello` request. We then attach
the healthcheck handler by using `healthz.NewHandler(mux)`.


## Creating a standalone server

When your application isn't a web server, but you still want to add health
checks, you can use the provided server implementation.

```go
func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	srv := healthz.NewServer("0.0.0.0", 3000)
	go srv.Start()

	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}
```

## Using middleware

With both the `NewHandler` and `NewServer` implementation, we've provided a way
to add middleware. This allows you to add logging to your health checks for
example.

The middleware functions are respectively `NewHandlerWithMiddleware` and
`NewServerWithMiddleware`. They accept the arguments of their parent function
but also a set of middlewares as variadic arguments.

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := healthz.NewHandlerWithMiddleware(mux, logMiddleware)
	http.ListenAndServe(":3000", handler)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("start request")
		next.ServeHTTP(w, r)
		log.Print("end request")
	})
}
```

## Output

On the health check endpoint, we return a set of values useful to us. Extending
the example from above, these are both return values (one where the file is
present, one where it isn't).

### 200 OK
```json
{
  "checked_at": "2017-11-22T14:18:50.339Z",
  "duration_ms": 0,
  "result": "success",
  "tests": {
    "default": {
      "duration_ms": 0,
      "result": "success"
    },
    "jwt-key": {
      "duration_ms": 0,
      "result": "success"
    }
  }
}
```

### 503 Service Unavailable
```json
{
  "checked_at": "2017-11-22T14:18:50.339Z",
  "duration_ms": 1,
  "result": "failure",
  "tests": {
    "default": {
      "duration_ms": 0,
      "result": "success"
    },
    "jwt-key": {
      "duration_ms": 0,
      "result": "failure"
    }
  }
}
```

## Middlewares

We've included a set of standard middlewares that can be useful for general use.

### Cache

The cache middleware allows you to cache a response for a specific duration.
This prevents the health check to overload due to different sources asking for
the health status. This is especially useful when the health checks are used to
check the health of other services as well.

To use this middleware, simply add it to the chain:

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := healthz.NewHandlerWithMiddleware(
		mux,
		healthz.CacheMiddleware(5*time.Second),
	)
	http.ListenAndServe(":3000", handler)
}
```
