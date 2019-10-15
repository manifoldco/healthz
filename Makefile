ci: lint cover benchmark

.PHONY: ci

#################################################
# Bootstrapping for base golang package and tool deps
#################################################

vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

#################################################
# Test and linting
#################################################

test: vendor
	@$(TEST_ENV) CGO_ENABLED=0 go test $$(go list ./... | grep -v generated)

lint:
	GO111MODULE=on go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

benchmark: vendor
       @CGO_ENABLED=0 go test -v -run=XXX -bench=.

cover: vendor
	@CGO_ENABLED=0 go test -v -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v vendor)

.PHONY: $(LINTERS) test cover benchmark
