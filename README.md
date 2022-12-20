# MQTTDistributedBenck - worker
Implementação do worker responsável por realizar os experimentos requisitados pelo orquestrador, ambos trocam três tipos de mensagens: controle, comandos e autenticação. Ao iniciar, worker irá publicar uma mensagem no tópico de autenticação padrão e esperar seu id ser enviado pelo orquestrador, caso não receba resposta em um intervalo de tempo x, irá encerrar automaticamente. Não há saídas locais, o output padrão é sempre o orquestrador e suas configurações locais são feitas através das flags de configuração disponíveis.
## Dependências
* Golang
* Git
* Make
## Instalação
 ```
git clone -b worker https://github.com/namelew/MQTTDistributedBenck MQTTDBworker
cd MQTTDBworker
make
 ```
## Utilização
Para iniciar o worker execute o binário worker que será gerado no diretório bin, ele possui as seguintes flags de configuração:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| timeout | 5 | Number of time, in minutes, that the client will remain running after trigger|
| login_t | 30 | Time that client will wait Orquestrator login confirmation, in seconds|
| broker | `tcp://localhost:1883` | Communication broker to worker - orquestrator relation|
| isunix | true | flag that confirm if the client is running on a unixlike machine|
| tool | `source/tools/mqttloader/bin/mqttloader` | localization of the mqtt benckmark tool to the client experiment|