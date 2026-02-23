.PHONY: run build test clean docker-up docker-down

run:
	go run cmd/main.go

build:
	go build -o main cmd/main.go

test:
	go test -v ./tests/...

clean:
	rm -f main
	go clean

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

install:
	go mod download

migrate:
	go run cmd/main.go

lint:
	golangci-lint run
