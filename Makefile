build:
	go build ./cmd/main.go

run:
	go build ./cmd/main.go
	go run ./cmd/main.go

test:
	CGO_ENABLED=1 go test --short --cover ./...

up: build
	docker-compose up -d app

down:
	docker-compose down -v
