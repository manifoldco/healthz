LINTERS=\
	gofmt \
	golint \
	gosimple \
	vet \
	misspell \
	ineffassign \
	deadcode

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
	github.com/alecthomas/gometalinter

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	gometalinter --install

vendor:

.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Test and linting
#################################################

test: vendor $(GENERATED_NAMING_FILES)
	@CGO_ENABLED=0 go test -v $$(go list ./... | grep -v vendor)

lint: vendor $(LINTERS)

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	 ./... --enable

$(LINTERS): vendor $(GENERATED_NAMING_FILES)
	$(METALINT) $@

.PHONY: $(LINTERS) test

release:
ifneq ($(shell git rev-parse --abbrev-ref HEAD),master)
	$(error You are not on the master branch)
endif
ifneq ($(shell git status --porcelain),)
	$(error You have uncommitted changes on your branch)
endif
ifndef VERSION
	$(error You need to specify the version you want to tag)
endif
	git commit -m "Tagging v$(VERSION)"
	git tag v$(VERSION)
	git push
	git push --tags
