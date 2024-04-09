# MQTTPoller
 A MQTT Network Benckmark Tool based on MQTT Loader focus on distributed experiments
## Architecture
![]()
## Build
* The orquestrator and the worker code must be compile before use to the target OS
```
go mod tidy
go build -o bin/orquestrator src/core/cmd/orquestrator/main.go
go build -o bin/worker src/core/cmd/worker/main.go
```
* The web interface also must be compile to a javascript code runnable on Node.js or equivalent
```
npm install
npm run build
node serve.js
```
PS: Look the Dockerfile of each process to more details
## Running
All the processes have they on Dockerfile to build a personalized Docker Image. Therefore, a docker compose example is available.
### Docker
* To run the compose example
```
make example
```
* To run each container separeted is necessary to build the application images on images directory
```
cd ./images/[PROCESS_NAME]
docker build -t [IMAGE_NAME]:[VERSION_NAME] .
docker run [IMAGE_NAME]:[VERSION_NAME] -e [LIST_OF_ENVS]
```
### Local
* To run the Golang codes directly
```
go run ./src/core/cmd/[PROCESS_NAME]/main.go
```
* To build the Golang codes and run
```
go build -o bin/[PROCESS_NAME] src/core/cmd/[PROCESS_NAME]/main.go
./bin/[PROCESS_NAME] [OPTIONS]
```
## Environment Variables
Bellow, the envorinment variables of each Docker Image
### Orquestrator
```
TOLERANCE=5 # Tolerance for the worker messages in seconds
BROKER=tcp://localhost:1883 # Control broker to manage the workers
PORT=8000 # Control API port
```
### Worker
```
LTIMEOUT=30 # Timeout for each login attemps
LTHRESHOUT=-1 # Login attemps before quit
TOOL=./tools/mqttloader/bin/mqttloader # Tool filepath
BROKER=tcp://localhost:1883 # Control broker
```
### Web Interface
```
ORQUESTRATOR_ADRESS http://orquestrator:8000/ # Control API Address
```
## Compatible Tools
* [MQTTLoader](https://github.com/dist-sys/mqttloader)
