all: build

build:
	git clone -b orquestrator https://github.com/namelew/MQTTDistributedBenck dump/orquestrator
	git clone -b worker https://github.com/namelew/MQTTDistributedBenck dump/worker
	docker build -t mqttdb/orquestrator:1 -f "images/orquestrator.dockerfile" .
	docker build -t mqttdb/worker:1 -f "images/worker.dockerfile" .
	docker compose up -d

clean:
	rm -rf dump
	