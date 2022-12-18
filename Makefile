setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get .

run:
	go run cmd/$(cmd)/$(cmd).go

build:
	go build cmd/$(cmd)/$(cmd).go

lint:
	golangci-lint run --timeout=3m

format:
	gofmt -s -w .

test:
	go test -v ./...
