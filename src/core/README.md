# MQTTPoller - Application Core
O core da aplicação são os processos principais do MQTT Poller, aqui estão os códigos do orquestrador e dos trabalhadores (workers) e um exemplo de uso. O core utiliza o MQTT para a comunicação entre trabalhadores e orquestrador. Sua saída padrão é via API Rest no formato JSON, que é consumido pelo Frontend. Abaixo, segue as prinicipais dependências, e as rotas de controle do orquestrador.
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
| broker | `tcp://localhost:1883` | broker de comunicação entre orquestrador e trabalhadores|
| port   | 8000 | porta de comunicação da API Rest de comandos do orquestrador |

Após isso, será iniciado uma API Rest com as rotas abaixo

## Rotas
### Descrição
| Route | Method | Description |
|:-----|:--------------|:------------|
| `/orquestrator/worker` | GET | Retorna todos os workers conhecidos, listando seus ids e seus estados atuais |
| `/orquestrator/worker/:id` | GET | Retorna os dados apenas do worker com o `id` selecinado |
| `/orquestrator/experiment` | GET | Retorna todos os experimentos requisitados aos orquestrador, listando seus parâmetros, em quais workers foram executados e seus resultados |
| `/orquestrator/experiment/:id` | GET | Retorna os dados apenas do experimento com o `id` selecinado |
| `/orquestrator/experiment/start` | POST | Executa um experimento em um ou mais workers selecionados |
| `/orquestrator/experiment/cancel/:id` | POST | Cancela um experimento de id `id` em todos os workers que o estão executando |
| `/orquestrator/experiment/:id` | DELETE | Exclui o experimento de id `id` do registro do orquestrador |

### Mensagens (API)
* Requisição de Experimento

| Nome | Valor padrão | Descrição |
|:-----|:--------------|:------------|
| `id` | [] | Lista de workers que executaram o experimento. Caso vazio, irá executar em todos os workers registrados que estiverem online |
| `description` | {} | JSON com parâmetros de definição do experimento utilizado pelo MQTTLoader |

* Exemplo
```
{
    "id": ["id1", "id2"],
    "description":{
        "broker":	"test_broker",
        "attempts": 0,
        "port":	1883,
        "mqttVersion":	3,
        "numPublishers": 1000,
        "numSubscribers":	1000,
        "qosPublisher":	2,
        "qosSubscriber":	2,
        "retain":	false,
        "topic":	"test",
        "payload":	1000000,
        "numMessages":	100000000,
        "interval":	0,
        "execTime":	240
    }
}
```

* Resultado do Experimento

## Worker
### Utilização
Para iniciar o worker execute o binário worker que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| login_t | 30 | Tempo de espera máximo do worker para a resposta de login do orquestrador |
| login_th | -1 | Tentativas de login antes de desconectar (-1 equivale a infinito) |
| broker | `tcp://localhost:1883` | Broker de comunicação com o orquestrador|
| tool | `source/tools/mqttloader/bin/mqttloader` | Caminho do binário do MQTTLoader|
