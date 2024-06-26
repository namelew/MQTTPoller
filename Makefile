all: build

example:
	docker compose up -d
	docker compose logs -f
build:
	go mod tidy
	go build -o bin/orquestrator src/core/cmd/orquestrator/main.go
	go build -o bin/worker src/core/cmd/worker/main.go
clean:
	rm -rf bin
	rm -f *.bin
	rm -f *.data
	rm -f *.log
	rm -f *.conf
	rm -f *.db
	docker compose down --rmi all --volumes
	
