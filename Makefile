ci: lint test benchmark

.PHONY: ci

#################################################
# Bootstrapping for base golang package and tool deps
#################################################

CMD_PKGS=$(shell grep '	"' tools.go | awk -F '"' '{print $$2}')

define VENDOR_BIN_TMPL
vendor/bin/$(notdir $(1)): vendor/$(1) | vendor
	go build -a -o $$@ ./vendor/$(1)
VENDOR_BINS += vendor/bin/$(notdir $(1))
vendor/$(1): go.sum
	GO111MODULE=on go mod vendor
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))

$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: vendor/bin/%
gofmt-bin vet-bin:

vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

clean:
	go mod -n -modcache

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

#################################################
# Test and linting
#################################################

test: vendor
	@$(TEST_ENV) CGO_ENABLED=0 go test $$(go list ./... | grep -v generated)

.golangci.gen.yml: .golangci.yml
	$(shell awk '/enable:/{y=1;next} y == 0 {print}' $< > $@)

LINTERS=$(filter-out megacheck,$(shell awk '/enable:/{y=1;next} y != 0 {print $$2}' .golangci.yml))

lint: vendor/bin/golangci-lint vendor .golangci.gen.yml
	$< run -c .golangci.gen.yml $(LINTERS:%=-E %) ./...
	$< run -c .golangci.gen.yml -E megacheck ./...

benchmark: vendor
       @CGO_ENABLED=0 go test -v -run=XXX -bench=.

.PHONY: lint test benchmark
