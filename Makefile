all: build

example:
	docker compose up -d
	docker compose restart workers
	docker compose restart workers-2
	docker compose restart workers-3
dev:
	docker compose -f ./docker-compose-dev.yaml up -d
build:
	go mod tidy
	go build -o bin/orquestrator cmd/orquestrator/main.go
	go build -o bin/worker cmd/worker/main.go
clean:
	rm -rf bin
	docker compose down --rmi all --volumes
	
