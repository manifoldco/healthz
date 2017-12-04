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
