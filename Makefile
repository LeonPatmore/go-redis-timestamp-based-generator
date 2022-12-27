setup:
	docker network create -d bridge redis-timestamp
	docker run -P --network redis-timestamp --name redis -d redis:7
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

run:
	go run cmd/$(cmd)/main.go

build:
	go build cmd/$(cmd)/main.go

lint:
	golangci-lint run --timeout=3m

format:
	gofmt -s -w .

test:
	go test -v ./...

cli:
	docker run -it --name redis-cli --network redis-timestamp --rm redis redis-cli -h redis

test-docker:
	docker run --rm -it -p 4567:80 strm/helloworld-http
