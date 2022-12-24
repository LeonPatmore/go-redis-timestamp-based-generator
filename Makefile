setup:
	docker run -P --name redis -d redis:7
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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

cli:
	docker run -it --name redis-cli --network host --rm redis redis-cli

test-docker:
	docker run --rm -it -p 4567:80 strm/helloworld-http
