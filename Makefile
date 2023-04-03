all: build

docker:
	git clone https://github.com/namelew/MQTTDistributedBench images/orquestrator/dump
	git clone https://github.com/namelew/MQTTDistributedBench images/worker/dump
	docker compose up -d
	docker compose restart orquestrator
	docker compose restart workers
	docker compose restart workers-2
	docker compose restart workers-3
build:
	go mod tidy
	go build -o bin/orquestrator cmd/orquestrator/main.go
	go build -o bin/worker cmd/worker/main.go
clean:
	rm -rf bin
	rm -rf images/orquestrator/dump
	rm -rf images/worker/dump
	docker compose down --rmi all --volumes
	
