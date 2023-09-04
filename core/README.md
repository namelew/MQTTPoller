# MQTTPoller - Application Core
O core da aplicação é são os processos principais do MQTT Poller, aqui estão os códigos do Orquestrador e dos trabalhadores e um exemplo de uso. O core utiliza o MQTT para a comunicação entre trabalhadores e orquestrador. Sua saída padrão é via API Rest no formato JSON, que é consumido pelo Frontend. Abaixo, segue as prinicipais dependências, e as rotas de controle do orquestrador.
## Dependências
* Golang
* Docker (Para o exemplo)
* Make
## Instalação
 ```
git clone https://github.com/namelew/MQTTDistributedBench MQTTPoller
cd MQTTPoller/core
make
 ```
## Example - Docker
```
make example
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
    "id": [string],
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
| login_th | -1 | Login attemps before auth fail shutdown (-1 equals to one attemp) |
| broker | `tcp://localhost:1883` | Communication broker to worker - orquestrator relation|
| tool | `source/tools/mqttloader/bin/mqttloader` | localization of the mqtt benckmark tool to the client experiment|
