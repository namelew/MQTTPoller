# MQTTDistributedBenck - broker latency measure tool
Ferramenta feita para gerar testes distribuidos de latência, vazão e perda em redes baseadas no protocolo MQTT. Para isso, é utilizado um nodo principal, chamado orquestrador, para gerenciar nodos trabalhadores, denomidados workers, reponsáveis por executar os experimentos de forma concorrente através de uma ferramenta que pode ser escolhida pelo usuário.

Para a comunicação entre orquestrador e workers, é utilizado o protocolo MQTT. Deve existir um Broker de comunicação conhecido entre ambas as partes para que a aplicação funcione corretamente. O orquestrador deve ser iniciado antes dos workers, pois os workers são incapazes de esperar até ele estar disponível.
## Dependências
* Golang
* Git
* Make
## Instalação
 ```
git clone https://github.com/namelew/MQTTDistributedBenck MQTTDB
cd MQTTDB
make
 ```
## Example - Docker
```
mkdir dump
cd dump
git clone https://github.com/namelew/MQTTDistributedBenck orquestrator
git clone -b worker https://github.com/namelew/MQTTDistributedBenck worker
cd ..
docker build -t mqttdb/orquestrator:1 -f "images/orquestrator.dockerfile" .
docker build -t mqttdb/worker:1 -f "images/worker.dockerfile" .
docker compose up -d
```
## Utilização
Para iniciar o orquestrador execute o binário orquestrator, que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| tl | 5 | tempo de tolerância para batidas de coração dos workers|
| broker | `tcp://localhost:1883` | Communication broker to worker - orquestrator relation|

Após isso, abrirá um shell interativo de controle para a ferramenta. Ele aceita, quatro commandos diferentes:
| Command | Description |
|:-----|:------------|
| ls | lista todos os workers cadastrados durante a sessão atual|
| start | Inicia experimentos em 1 ou mais worker|
| info  | Recuperação informações de hardware e do sistema operacional de um ou mais workers |
| cancel | Cancela um experimento em execução em um worker|

### Exemplos
* Listando workers cadastrados
```
ls
```
* Listando experimentos realizados pelo worker 0
```
ls -i 0
```
* Disparando um experimento em todos os workers
```
start
```
* Disparando experimento no worker 0 com arquivo de configuração command.json
```
start -i 0 -f examples/command.json
```
* Cancelando experimento de id 1000000 no worker 0
```
cancel 0 1000000
```