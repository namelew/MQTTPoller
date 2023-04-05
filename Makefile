all: build

example:
	docker compose up -d
	docker compose logs -f
dev:
	docker compose -f ./docker-compose-dev.yaml up -d
	docker compose -f ./docker-compose-dev.yaml logs -f
build:
	go mod tidy
	go build -o bin/orquestrator cmd/orquestrator/main.go
	go build -o bin/worker cmd/worker/main.go
clean:
	rm -rf bin
	docker compose down --rmi all --volumes
	
