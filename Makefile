include build/Makefile.core.mk
include build/Makefile.show-help.mk

## Run suite of Golang lint checks
check:
	golangci-lint run -c .golangci.yml -v ./...

## Run Golang tests
test:
	go test --short ./... -race -coverprofile=coverage.txt -covermode=atomic

## Clean up Golang packages. Print diff.
tidy:
	go mod tidy
	git diff --exit-code go.mod go.sum

## Run Meshplay Error Code Utility. Generate error codes.
errorutil:
	go run github.com/khulnasoft/meshkit/cmd/errorutil -d . update --skip-dirs meshplay -i ./helpers -o ./helpers

## Run Meshplay Error Code Utility. Analyze only.
errorutil-analyze:
	go run github.com/khulnasoft/meshkit/cmd/errorutil -d . analyze --skip-dirs meshplay -i ./helpers -o ./helpers

## Build the Meshplay Error Code Utility. 
build-errorutil:
	go build -o errorutil cmd/errorutil/main.go
