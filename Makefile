all: build

docker: build
	git clone https://github.com/namelew/MQTTDistributedBenck images/orquestrator/dump
	git clone https://github.com/namelew/MQTTDistributedBenck images/worker/dump
	docker compose up -d

all: build

build:
	go mod tidy
	go build -o bin/orquestrator /cmd/orquestrator/main.go
	go build -o bin/worker /cmd/worker/main.go
clean:
	rm -rf bin
	rm -rf images/orquestrator/dump
	rm -rf images/worker/dump
	docker compose down --rmi all --volumes
	