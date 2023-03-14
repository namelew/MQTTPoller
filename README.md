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
git clone -b orquestrator https://github.com/namelew/MQTTDistributedBenck images/orquestrator/dump
git clone -b worker https://github.com/namelew/MQTTDistributedBenck images/worker/dump
docker compose up -d
```
## Utilização
Para iniciar o orquestrador execute o binário orquestrator, que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| tl | 5 | tempo de tolerância para batidas de coração dos workers|
| broker | `tcp://localhost:1883` | communication broker to worker - orquestrator relation|
| port   | 8000 | api rest communication port |

Após isso, abrirá um shell interativo de controle para a ferramenta. Ele aceita, quatro commandos diferentes:
| Command | Description |
|:-----|:------------|
| ls | lista todos os workers cadastrados durante a sessão atual|
| start | inicia experimentos em 1 ou mais worker|
| info  | recuperação informações de hardware e do sistema operacional de um ou mais workers |
| cancel | cancela um experimento em execução em um worker|

### Exemplos