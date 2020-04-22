SHELL=/bin/bash
.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Resolve dependencies using Go Modules
	GO111MODULE=on go mod download

.PHONY: test
test: ## Tests all code
	GO111MODULE=on go test -cover -race ./...

.PHONY: lint
lint: ## Runs static code analysis
	command -v goimports >/dev/null 2>&1 || { go get -u golang.org/x/tools/cmd/goimports; }
	diff <(goimports -d .) <(printf "")
	command -v golint >/dev/null 2>&1 || { go get -u golang.org/x/lint/golint; }
	golint -set_exit_status ./...
	command -v errcheck >/dev/null 2>&1 || { go get -u github.com/kisielk/errcheck; }
	errcheck ./...

.PHONY: count-go
count-go: ## Count number of lines of all go codes
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
