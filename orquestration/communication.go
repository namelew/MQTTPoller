package orquestration

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/databases"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/databases/services/experiments"
	seworkers "github.com/namelew/mqtt-bm-latency/databases/services/workers"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/logs"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
)

var oLog = logs.Build("orquestrator.log")
var serviceWorkers = seworkers.Build(oLog)
var serviceExperiments = experiments.Build(oLog)
var infos = make([]output.Info, 0, 10)
var workers = make([]messages.Worker, 1, 10)
var rexp []output.ExperimentResult
var rexpMutex sync.Mutex
var waitQueueMutex sync.Mutex
var waitQueue []output.ExperimentResult
var expWG sync.WaitGroup
var client mqtt.Client

func GetWorkers(filter *filters.Worker) []models.Worker {
	return serviceWorkers.List(filter)
}

func Init(broker string, t_interval int) error {
	var clientID string = "Orquestrator"
	ka, err := time.ParseDuration(strconv.Itoa(10000) + "s")

	if err != nil {
		return err
	}

	oLog.Create()

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(ka).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	client = mqtt.NewClient(opts)

	oLog.Register("connect paho mqtt client to broker " + broker)

	tokenConnection := client.Connect()

	tokenConnection.Wait()

	oLog.Register("Starting database")

	databases.Connect(oLog)

	token := client.Subscribe("Orquestrator/Register", byte(1), Register)
	token.Wait()

	token = client.Subscribe("Orquestrator/Login", byte(1), Login)
	token.Wait()

	token = client.Subscribe("Orquestrator/Ping", byte(1), func(c mqtt.Client, m mqtt.Message) { Ping(c, m, t_interval) })
	token.Wait()

	return nil
}

func End() {
	oLog.Register("disconnect mqtt.client")

	client.Disconnect(0)

	oLog.Register("shutdown")
}

func setMessageHandler(t *string) {
	token := client.Subscribe(*t+"/Experiments/Results", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperiment(m)
	})
	token.Wait()
	token = client.Subscribe(*t+"/Experiments/Status", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperimentStatus(m)
	})
	token.Wait()
}

func StartExperiment(arg input.Start) ([]output.ExperimentResult, error) {
	expid := time.Now().Unix()

	rexpMutex.Lock()
	rexp = nil
	rexpMutex.Unlock()

	var cmd messages.Command
	var experiment messages.CommandExperiment

	cmd.Name = "start"
	cmd.CommandType = "experiment command"

	experiment.Expid = expid
	experiment.Attempts = arg.Description.Attempts
	experiment.Broker = arg.Description.Broker
	experiment.ExecTime = arg.Description.ExecTime
	experiment.Interval = arg.Description.Interval
	experiment.LogLevel = arg.Description.LogLevel
	experiment.MqttVersion = arg.Description.MqttVersion
	experiment.Ntp = arg.Description.Ntp
	experiment.NumMessages = arg.Description.NumMessages
	experiment.NumPublishers = arg.Description.NumPublishers
	experiment.NumSubscriber = arg.Description.NumSubscriber
	experiment.Output = arg.Description.Output
	experiment.Password = arg.Description.Password
	experiment.Payload = arg.Description.Payload
	experiment.Port = arg.Description.Port
	experiment.QosPublisher = arg.Description.QosPublisher
	experiment.QosSubscriber = arg.Description.QosSubscriber
	experiment.RampDown = arg.Description.RampDown
	experiment.RampUp = arg.Description.RampUp
	experiment.Retain = arg.Description.Retain
	experiment.SharedSubscrition = arg.Description.SharedSubscrition
	experiment.SubscriberTimeout = arg.Description.SubscriberTimeout
	experiment.TlsKeystore = arg.Description.TlsKeystore
	experiment.TlsKeystorePassword = arg.Description.TlsKeystorePassword
	experiment.TlsTrustsore = arg.Description.TlsTrustsore
	experiment.TlsTruststorePassword = arg.Description.TlsTruststorePassword
	experiment.Tool = arg.Description.Tool
	experiment.Topic = arg.Description.Topic
	experiment.User = arg.Description.User

	err := experiment.Attach(&cmd)

	if err != nil {
		return rexp, err
	}

	msg, err := json.Marshal(cmd)

	if err != nil {
		return rexp, err
	}

	if arg.Id[0] == -1 {
		nw := len(workers)
		expWG.Add(nw)
		for i := 0; i < nw; i++ {
			if !workers[i].Status {
				log.Printf("Worker %d is off, skipping\n", i)
				continue
			}

			workers[i].Historic.Add(expid, cmd, arg.Description.Attempts)

			token := client.Publish(workers[i].Id+"/Command", byte(1), false, msg)
			token.Wait()

			go receiveControl(i, arg.Description.ExecTime*5)

			log.Printf("Requesting experiment in worker %d\n", i)
		}
	} else {
		argTam := len(arg.Id)
		expWG.Add(argTam)
		for i := 0; i < argTam; i++ {
			if !workers[arg.Id[i]].Status {
				if argTam > 1 {
					log.Printf("Worker %d is off, skipping\n", arg.Id[i])
					continue
				} else {
					log.Printf("Worker %d is off, aborting experiment\n", arg.Id[i])
					break
				}
			}

			workers[i].Historic.Add(expid, cmd, arg.Description.Attempts)

			token := client.Publish(workers[arg.Id[i]].Id+"/Command", byte(1), false, msg)
			token.Wait()

			go receiveControl(arg.Id[i], arg.Description.ExecTime*5)

			log.Printf("Requesting experiment in worker %d\n", arg.Id[i])
		}
	}

	expWG.Wait()

	rexpMutex.Lock()
	waitQueueMutex.Lock()
	rexp = append(rexp, waitQueue...)
	rexpMutex.Unlock()

	waitQueue = nil
	waitQueueMutex.Unlock()

	return rexp, nil
}

func CancelExperiment(id int, expid int64) error {
	arg := make(map[string]interface{})
	arg["id"] = expid
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: arg}
	msg, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	token := client.Publish(workers[id].Id+"/Command", byte(1), false, msg)
	token.Wait()

	workers[id].ReceiveConfirmation = true
	exp := workers[id].Historic.Search(expid)
	exp.Finished = true

	return nil
}

func GetInfo(arg input.Info) ([]output.Info, error) {
	var infoCommand messages.Command
	infos = nil

	infoCommand.Name = "info"
	infoCommand.CommandType = "command moderation"
	infoCommand.Arguments = map[string]interface{}{"cpuDisplay": arg.CpuDisplay, "discDisplay": arg.DiscDisplay, "memoryDisplay": arg.MemoryDisplay}

	msg, err := json.Marshal(&infoCommand)

	if err != nil {
		return infos, err
	}

	if arg.Id[0] == -1 {
		for i := 0; i < len(workers); i++ {
			if !workers[i].Status {
				log.Printf("Worker %d isn't report, skipping\n", i)
				continue
			}
			token := client.Subscribe(workers[i].Id+"/Info", byte(1), func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, i)
			})
			token.Wait()

			token = client.Publish(workers[i].Id+"/Command", byte(1), false, msg)
			token.Wait()

			for !workers[i].ReceiveConfirmation {
				if !workers[i].Status {
					log.Printf("Worker %d isn't report, skipping\n", i)
					break
				}
				time.Sleep(time.Second)
			}
			workers[i].ReceiveConfirmation = false
			token = client.Unsubscribe(workers[i].Id + "/Info")
			token.Wait()
		}
	} else {
		argTam := len(arg.Id)

		for i := 0; i < argTam; i++ {
			if !workers[arg.Id[i]].Status {
				if argTam > 1 {
					log.Printf("Worker %d is off, skipping\n", arg.Id[i])
					continue
				} else {
					log.Printf("Worker %d is off, aborting request\n", arg.Id[i])
					break
				}
			}

			token := client.Subscribe(workers[arg.Id[i]].Id+"/Info", byte(1), func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, arg.Id[i])
			})
			token.Wait()

			token = client.Publish(workers[arg.Id[i]].Id+"/Command", byte(1), false, msg)
			token.Wait()

			for !workers[arg.Id[i]].ReceiveConfirmation {
				if !workers[arg.Id[i]].Status {
					if argTam > 1 {
						log.Printf("Worker %d is off, skipping\n", arg.Id[i])
					} else {
						log.Printf("Worker %d is off, aborting request\n", arg.Id[i])
					}
					break
				}
				time.Sleep(time.Second)
			}
			workers[arg.Id[i]].ReceiveConfirmation = false
			token = client.Unsubscribe(workers[i].Id + "/Info")
			token.Wait()
		}
	}

	return infos, nil
}
