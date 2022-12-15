# Running

The client have five control flags, all have default values and they are was follow:

| Flag | Default value | Description |
|:-----|:--------------|:------------|
| timeout | 5 | Number of time, in minutes, that the client will remain running after trigger|
| login_t | 30 | Time that client will wait Orquestrator login confirmation, in seconds|
| broker | `tcp://localhost:1883` | Communication broker to worker - orquestrator relation|
| isunix | true | flag that confirm if the client is running on a unixlike machine|
| tool | `source/tools/mqttloader/bin/mqttloader` | localization of the mqtt benckmark tool to the client experiments |

Examples:
 - Setting a new client timeout
```
./worker --timeout 10
```