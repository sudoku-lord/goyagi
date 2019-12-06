GOTOOLS := \
	github.com/alecthomas/gometalinter \
	github.com/git-chglog/git-chglog/cmd/git-chglog \
	github.com/golang/dep/cmd/dep \
	golang.org/x/tools/cmd/cover \

GOPACKAGES := $(go list ./...)

COVERAGE_PROFILE := coverage.out

## default command
.DEFAULT_GOAL := help

.PHONY: clean
clean: ## Removes Go temporary build files build directory
	@echo "---> Cleaning..."
	go clean
	rm -rf build

.PHONY: enforce
enforce: ## Enforces code coverage
	@echo "---> Enforcing coverage..."
	./scripts/coverage.sh $(COVERAGE_PROFILE)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: html
html: ## Generates an HTML coverage report
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: install
install: ## Installs dependencies
	@echo "---> Installing dependencies..."
	dep ensure

.PHONY: lint
lint: ## Runs all linters
	@echo "---> Linting..."
	gometalinter

.PHONY: release
release: ## Creates a new release with the given tag
	@echo "---> Creating new release"
ifndef tag
	$(error tag must be specified)
endif
	git-chglog --output CHANGELOG.md --next-tag $(tag)
	git add CHANGELOG.md
	git commit -m $(tag)
	git tag $(tag)
	git push origin master --tags

.PHONY: setup
setup: ## Installs all development dependencies
	@echo "---> Setting up..."
	go get -u $(GOTOOLS)
	gometalinter --install

.PHONY: test
test: ## Runs all the tests and outputs the coverage report
	@echo "---> Testing..."
	RELEASE=test12345 go test ./... -race -coverprofile=$(COVERAGE_PROFILE)
