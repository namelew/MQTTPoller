package orquestration

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
	"github.com/namelew/mqtt-bm-latency/utils"
)

var workers = make([]messages.Worker, 1, 10)
var infos = make([]output.Info, 0, 10)
var rexp []output.ExperimentResult
var expWG sync.WaitGroup
var client mqtt.Client

func GetWorkers() []messages.Worker {
	return workers
}

func Init() {
	var (
		broker     = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		t_interval = flag.Int("tl", 5, "orquestrator tolerance interval")
	)
	flag.Parse()
	var currentSession messages.Session
	currentSession.Finish = true
	var clientID string = "Orquestrator"
	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")

	if !utils.FileExists("orquestrator.log") {
		f, _ := os.Create("orquestrator.log")
		f.Close()
	} else {
		os.Truncate("orquestrator.log", 0)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(*broker).
		SetClientID(clientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(ka).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	client = mqtt.NewClient(opts)

	f, _ := os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString("connect mqtt.client\n")
	f.Close()

	tokenConnection := client.Connect()

	tokenConnection.Wait()

	token := client.Subscribe("Orquestrator/Sessions", byte(1), func(c mqtt.Client, m mqtt.Message) {
		err := json.Unmarshal(m.Payload(), &currentSession)
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	token.Wait()

	token = client.Unsubscribe("Orquestrator/Sessions")
	token.Wait()

	token = client.Subscribe("Orquestrator/Register", byte(1), func(c mqtt.Client, m mqtt.Message) {
		var clientID string = ""
		var seed rand.Source
		var random *rand.Rand

		for i := 0; i < 10; i++ {
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			clientID += fmt.Sprintf("%d", random.Int()%10)
		}

		f, _ := os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("worker " + clientID + " register\n")
		f.Close()

		if workers[0].Id == "" {
			workers[0] = messages.Worker{Id: clientID, Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}}

			t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
			t.Wait()

			setMessageHandler(0)
			return
		}
		workers = append(workers, messages.Worker{Id: clientID, Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}})

		t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
		t.Wait()

		setMessageHandler(len(workers) - 1)
	})
	token.Wait()

	token = client.Subscribe("Orquestrator/Login", byte(1), func(c mqtt.Client, m mqtt.Message) {
		f, _ := os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("worker " + string(m.Payload()) + " login\n")
		f.Close()
		if !utils.IsIn(workers, string(m.Payload())) {
			if workers[0].Id == "" {
				workers[0] = messages.Worker{Id: string(m.Payload()), Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}}

				t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
				t.Wait()

				setMessageHandler(0)
				return
			}
			workers = append(workers, messages.Worker{Id: string(m.Payload()), Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}})

			t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
			t.Wait()

			setMessageHandler(len(workers) - 1)
		} else {
			id := 0
			for i := 0; i < len(workers); i++ {
				if workers[i].Id == string(m.Payload()) {
					id = i
					break
				}
			}
			workers[id].Status = true
			setMessageHandler(id)
		}
	})
	token.Wait()

	token = client.Subscribe("Orquestrator/Ping", byte(1), func(c mqtt.Client, m mqtt.Message) {
		id := 0
		for i := 0; i < len(workers); i++ {
			if workers[i].Id == string(m.Payload()) {
				workers[i].Status = true
				id = i
				break
			}
		}
		if workers[id].TestPing {
			workers[id].TestPing = false
			go watcher(id, *t_interval)
		} else {
			workers[id].TestPing = true
			go watcher(id, *t_interval)
			workers[id].TestPing = false
		}
	})

	token.Wait()
}

func End() {
	f, _ := os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString("disconnect mqtt.client\n")

	client.Disconnect(0)

	f.WriteString("shutdown")
	f.Close()
}

func setMessageHandler(id int) {
	token := client.Subscribe(workers[id].Id+"/Experiments/Results", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperiment(m, id)
	})
	token.Wait()
	token = client.Subscribe(workers[id].Id+"/Experiments/Status", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperimentStatus(m, id)
	})
	token.Wait()
}
// argumentos estão chegando corretamente mas continua dando erro
func StartExperiment(arg input.Start) ([]output.ExperimentResult){
	expid := time.Now().Unix()
	rexp = nil
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
	
	experiment.Attach(&cmd)

	msg,err := json.Marshal(cmd)

	if err != nil{
		log.Fatal(err.Error())
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

	return rexp
}

func cancelExperiment(id int, expid int64) {
	arg := make(map[string]interface{})
	arg["id"] = expid
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: arg}
	msg, _ := json.Marshal(cmd)

	token := client.Publish(workers[id].Id+"/Command", byte(1), false, msg)
	token.Wait()

	workers[id].ReceiveConfirmation = true
	exp := workers[id].Historic.Search(expid)
	exp.Finished = true
}

func GetInfo(arg input.Info) []output.Info{
	var infoCommand messages.Command
	infos = nil

	infoCommand.Name = "info"
	infoCommand.CommandType = "command moderation"
	infoCommand.Arguments = map[string]interface{}{"cpuDisplay": arg.CpuDisplay, "discDisplay": arg.DiscDisplay, "memoryDisplay": arg.MemoryDisplay}

	msg, _ := json.Marshal(&infoCommand)

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

	return infos
}
