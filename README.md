MQTT Client - broker latency measure tool
=========

A MQTT Client from a distributed broker latency tool.
Scenario: when started in a machine in a mqtt network, will wait to command to start benchmark routines from a reachable source, the orquestrator. In the end, will generate a response to the orquestrator request in a mqtt topic.
This client can generate reports with the follow measures: throughput, number of messages published or received, latency and host machine infos(memory, cpu and disc).

Installation:

```
git clone -b mqtt-client https://github.com/namelew/mqtt-bm-latency
cd mqtt-bm-latency
make
```

All dependencies are vendored with [manul](https://github.com/kovetskiy/manul).

# Client Routines
The client had four work routines: login a orquestrator, return machine infos, start a mqtt bm latency experiment and return the results and client status by ping.
## Login
Client infos the orquestrator through a knowing broker his existence and wait to orquestrator confirmation. The wait time and the communication broker can be pass has args to the client like follow:
 ```
 worker --broker "tcp://localhost:1883" --login_t 30
 ```
 ## Machine Infos
 The orquestrator send the command "Info" by a message in one predetermined mqtt topic and client will send back the machine infos: memory, cpu and disc. All messages are in JSON format and what info will be display can be controlled.

 Messages structure:
  * Orquestrator command
```
{
    "name": "info",
    "commandType": "command moderation",
    "arguments": {
        "memoryDisplay": true,
        "cpuDisplay": true,
        "discDisplay": true
    }
}
```
 * Client Response
 ```
 {
       "Cpu":"Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz",
       "Ram":7837,
       "Disk":467921
}
 ```

 ## Start Experiment
By receiving the "Start" Command from the orquestrator, client will start a benckmark experiment in a reachable broker. Multiples experiments can be start concurrently and at the end, the results will be send back to the orquestrator.

Message structure:
 * Orquestrator "Start" command:
```
{
    "name": "start",
    "type": "experiment command",
    "arguments":{
        "tool":	"mqtt loader",
        "broker":	"127.0.0.1",
        "broker_port":	1883,
        "mqtt_version":	3,
        "num_publishers": 10,
        "num_subscribers":	10,
        "qos_publisher":	1,
        "qos_subscriber":	1,
        "shared_subscription":	false,
        "retain":	false,
        "topic":	"test",
        "payload":	10,
        "num_messages":	10,
        "ramp_up":	false,
        "ramp_down": false,
        "interval":	0,
        "subscriber_timeout":3,
        "exec_time":	5,
        "log_level":	"INFO"
    }
}
```
 * Client Response
```
{ 
      "meta":
            {
                   "id":1655140257325689,
                   "error":"",
                   "tool":"mqttLoader","literal":"\nMeasurement started: 2022-06-13 14:11:04.352 BRT\nMeasurement ended: 2022-06-13 14:11:08.376 BRT\n\n-----Publisher-----\nMaximum throughput [msg/s]: 100\nAverage throughput [msg/s]: 100,000\nNumber of published messages: 100\nPer second throughput [msg/s]: 100\n\n-----Subscriber-----\nMaximum throughput [msg/s]: 1000\nAverage throughput [msg/s]: 1000,000\nNumber of received messages: 1000\nPer second throughput [msg/s]: 1000\nMaximum latency [ms]: 29,499\nAverage latency [ms]: 14,426\n",
                  "log_file":
                        {
                              "name":"",
                              "data":null,
                              "extension":""
                  }
            },
      "publish":
            {
                  "max_throughput":100,
                  "avg_throughput":100,
                  "publiqued_messages":100,
                  "per_second_throungput":100
            },
      "subscribe":
            {
                  "max_throughput":1000,"avg_throughput":1000,"received_messages":1000,"per_second_throungput":1000,"latency":29.499,
                  "avg_latency":14.426
            }
}
```
## Ping
To confirm if clients is online, orquestrator can send a ping message that will be send back from the client if he is online. Client will be set has offline for orquestrator if don't receive a confirmation after 3 attempts.

## Manuals
 * [User manual](https://github.com/namelew/mqtt-bm-latency/tree/mqtt-client/examples/usage.md)