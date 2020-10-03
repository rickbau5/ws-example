BIN ?= bin/ws-example
GO_MOD_MODE ?= vendor
GOLANG_CI_LINT_VERSION ?= v1.31

.PHONY: all
all: lint test build

.PHONY: build
build:
	go build -mod=$(GO_MOD_MODE) -o $(BIN) -v .

.PHONY: install
install: .git/hooks/pre-commit

.git/hooks/pre-commit:
	ln pre-commit.sh .git/hooks/pre-commit

.PHONY: lint
lint:
	if gofmt -l -d -e . | grep "^" ; then exit 1 ; fi

.PHONY: test
test: vendor
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html coverage.out -o coverage.html
	go tool cover -func coverage.out

vendor: vendor/modules.txt

vendor/modules.txt: go.mod
	go mod vendor

.PHONY: clean
clean: .clean-go

.PHONY: .clean-go
.clean-go:
	rm $(BIN)
	rm coverage.out coverage.html
	rm -r vendor
