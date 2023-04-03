# MQTTDistributedBench - broker latency measure tool
Ferramenta feita para gerar testes distribuidos de latência, vazão e perda em redes baseadas no protocolo MQTT. Para isso, é utilizado um nodo principal, chamado orquestrador, para gerenciar nodos trabalhadores, denomidados workers, reponsáveis por executar os experimentos de forma concorrente através de uma ferramenta que pode ser escolhida pelo usuário.

Para a comunicação entre orquestrador e workers, é utilizado o protocolo MQTT. Deve existir um Broker de comunicação conhecido entre ambas as partes para que a aplicação funcione corretamente. O orquestrador deve ser iniciado antes dos workers, pois os workers são incapazes de esperar até ele estar disponível.
## Dependências
* Golang
* Git
* Make
## Instalação
 ```
git clone https://github.com/namelew/MQTTDistributedBench MQTTDB
cd MQTTDB
make
 ```
## Example - Docker
```
make docker
```
## Orquestrador
### Utilização
Para iniciar o orquestrador execute o binário orquestrator, que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| tl | 5 | tempo de tolerância para batidas de coração dos workers|
| broker | `tcp://localhost:1883` | communication broker to worker - orquestrator relation|
| port   | 8000 | api rest communication port |

Após isso, será iniciado uma API Rest com as rotas abaixo

## Rotas
### Descrição
| Route | Method | Description |
|:-----|:--------------|:------------|
| `/orquestrator/worker` | GET | Retorna todos os workers conhecidos, seus estados e o experimento que estão executando |
| `/orquestrator/info` | GET | Retorna informações sobre o máquina onde o worker está executando|
| `/orquestrator/experiment/start` | POST | Executa um experimento em um ou mais workers selecionados|
| `/orquestrator/experiment/cancel/:id/:expid` | DELETE | Cancela um experimento de id `expid` que está executando num worker `id`|
### Mensagens
#### /orquestrator/worker/:id
* request
```
    vars in url
```
* response
```
[
    {
        "Id": int,
        "NetId": string
        "Online": bool
        "History": [
            {
                "Command": description
                "Finished": bool
                "Id": int
            },
            ...
        ]
    },
    ...
]
```
or 
```
{
    "Id": int,
    "NetId": string
    "Online": bool
    "History": [
        {
            "Command": description
            "Finished": bool
            "Id": int
        },
        ...
    ]
}
```
#### /orquestrator/experiment/start
* request
```
{
    "id": [ids],
    "description":{
        "tool":	tool string name,
        "broker":	broker ip/dns,
        "attempts": int,
        "port":	broker port,
        "mqttVersion":	3|5,
        "numPublishers": int,
        "numSubscribers":	int,
        "qosPublisher":	0|1|2,
        "qosSubscriber":	0|1|2,
        "sharedSubscription":	bool,
        "retain":	bool,
        "topic":	topic name,
        "payload":	message size,
        "numMessages":	int,
        "ramUp":	ramp up time,
        "rampDown": ramp down time,
        "interval":	interval beetwen messages,
        "subscriberTimeout": int second,
        "execTime":	int second,
        "logLevel":	"INFO"|"SEVERE"|"WARNING"|"ALL",
        "ntp": ntp server adress,
        "output": get output file bool,
        "username": mqtt client username string,
        "password": mqtt client password string,
        "tlsTruststore": string file path,
        "tlsTruststorePass": string key file path,
        "tlsKeystore": string file path,
        "tlsKeystorePass": string key file path
    },
    "exeMode":unused attribute
}
```
* response
```
[
    {
        "meta": meta object
        "publish": publishers result object
        "subscribe": subscribers result object
    },
    ...
]
```
#### /orquestrator/experiment/cancel/:id/:expid
* request
```
    vars in url
```
* response
```
    request status code
    null | error
```
## Worker
### Utilização
Para iniciar o worker execute o binário worker que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| timeout | 5 | Number of time, in minutes, that the client will remain running after trigger|
| login_t | 30 | Time that client will wait Orquestrator login confirmation, in seconds|
| broker | `tcp://localhost:1883` | Communication broker to worker - orquestrator relation|
| isunix | true | flag that confirm if the client is running on a unixlike machine|
| tool | `source/tools/mqttloader/bin/mqttloader` | localization of the mqtt benckmark tool to the client experiment|
