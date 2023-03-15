package communication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/communication/logs"
	"github.com/namelew/mqtt-bm-latency/history"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/utils"
)

type Worker struct {
	Id     string
	tool   string
	broker string
	isUnix bool
	client mqtt.Client
}

var log *logs.Log = logs.Build("worker.log")
var experimentListMutex sync.Mutex
var experimentList history.OngoingExperiments

func Build() *Worker {
	return &Worker{}
}

func loadArguments(file string, arg map[string]interface{}) (bool, int64) {
	var arguments messages.CommandExperiment
	jsonObj, _ := json.Marshal(arg)
	json.Unmarshal(jsonObj, &arguments)

	f, _ := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	var argf string = ""

	isNull := func(param interface{}) string {
		switch param.(type) {
		case string:
			if param == "" {
				return "#"
			}
		case int:
			if param == 0 {
				return "#"
			}
		}
		return ""
	}

	argf += fmt.Sprintf("%sbroker=%s\n", isNull(arguments.Broker), arguments.Broker)
	argf += fmt.Sprintf("%sbroker_port=%d\n", isNull(arguments.Port), arguments.Port)
	argf += fmt.Sprintf("%smqtt_version=%d\n", isNull(arguments.MqttVersion), arguments.MqttVersion)
	argf += fmt.Sprintf("%snum_publishers=%d\n", isNull(arguments.NumPublishers), arguments.NumPublishers)
	argf += fmt.Sprintf("%snum_subscribers=%d\n", isNull(arguments.NumSubscriber), arguments.NumSubscriber)
	argf += fmt.Sprintf("%sqos_publisher=%d\n", isNull(arguments.QosPublisher), arguments.QosPublisher)
	argf += fmt.Sprintf("%sqos_subscriber=%d\n", isNull(arguments.QosSubscriber), arguments.QosSubscriber)
	argf += fmt.Sprintf("%sshared_subscription=%t\n", isNull(arguments.SharedSubscrition), arguments.SharedSubscrition)
	argf += fmt.Sprintf("%sretain=%t\n", isNull(arguments.Retain), arguments.Retain)
	argf += fmt.Sprintf("%stopic=%s\n", isNull(arguments.Topic), arguments.Topic)
	argf += fmt.Sprintf("%spayload=%d\n", isNull(arguments.Payload), arguments.Payload)
	argf += fmt.Sprintf("%snum_messages=%d\n", isNull(arguments.NumMessages), arguments.NumMessages)
	argf += fmt.Sprintf("%sramp_up=%d\n", isNull(arguments.RampUp), arguments.RampUp)
	argf += fmt.Sprintf("%sramp_down=%d\n", isNull(arguments.RampDown), arguments.RampDown)
	argf += fmt.Sprintf("%sinterval=%d\n", isNull(arguments.Interval), arguments.Interval)
	argf += fmt.Sprintf("%ssubscriber_timeout=%d\n", isNull(arguments.SubscriberTimeout), arguments.SubscriberTimeout)
	argf += fmt.Sprintf("%slog_level=%s\n", isNull(arguments.LogLevel), arguments.LogLevel)
	argf += fmt.Sprintf("%sexec_time=%d\n", isNull(arguments.Exec_time), arguments.Exec_time)
	argf += fmt.Sprintf("%sntp=%s\n", isNull(arguments.Ntp), arguments.Ntp)
	if arguments.Output {
		argf += fmt.Sprintf("output=%s\n", "output")
	}
	argf += fmt.Sprintf("%suser_name=%s\n", isNull(arguments.User), arguments.User)
	argf += fmt.Sprintf("%spassword=%s\n", isNull(arguments.Password), arguments.Password)
	argf += fmt.Sprintf("%stls_truststore=%s\n", isNull(arguments.TlsTrustsore), arguments.TlsTrustsore)
	argf += fmt.Sprintf("%stls_truststore_pass=%s\n", isNull(arguments.TlsTruststorePassword), arguments.TlsTruststorePassword)
	argf += fmt.Sprintf("%stls_keystore=%s\n", isNull(arguments.TlsKeystore), arguments.TlsKeystore)
	argf += fmt.Sprintf("%stls_keystore_pass=%s\n", isNull(arguments.TlsKeystorePassword), arguments.TlsKeystorePassword)

	// Write bytes to file
	byteSlice := []byte(argf)
	_, err := f.Write(byteSlice)
	if err != nil {
		log.Register("Load arguments error " + err.Error())
		f.Close()
	}

	return arguments.Output, arguments.Expid
}

func extracExperimentResults(output string, createLog bool) messages.ExperimentResult {
	results := messages.ExperimentResult{}
	results.Meta.Literal = output

	temp := [12]string{}

	i := 0

	for _, s := range strings.Split(output, "\n") {
		if s != "" && s[0] != '-' {
			temp[i] = strings.Split(s, ": ")[1]
			i++
		}
	}
	results.Meta.ToolName = "mqttLoader"
	results.Meta.ExperimentError = ""

	results.Publish.Throughput, _ = strconv.ParseFloat(strings.Replace(temp[2], ",", ".", 1), 64)
	results.Publish.AvgThroughput, _ = strconv.ParseFloat(strings.Replace(temp[3], ",", ".", 1), 64)
	results.Publish.PubMessages, _ = strconv.Atoi(temp[4])

	for _, s := range strings.Split(temp[5], ", ") {
		aux, _ := strconv.Atoi(s)
		results.Publish.PerSecondThrouput = append(results.Publish.PerSecondThrouput, aux)
	}

	results.Subscribe.Throughput, _ = strconv.ParseFloat(strings.Replace(temp[6], ",", ".", 1), 64)
	results.Subscribe.AvgThroughput, _ = strconv.ParseFloat(strings.Replace(temp[7], ",", ".", 1), 64)
	results.Subscribe.ReceivedMessages, _ = strconv.Atoi(temp[8])

	for _, s := range strings.Split(temp[9], ", ") {
		aux, _ := strconv.Atoi(s)
		results.Subscribe.PerSecondThrouput = append(results.Subscribe.PerSecondThrouput, aux)
	}

	results.Subscribe.Latency, _ = strconv.ParseFloat(strings.Replace(temp[10], ",", ".", 1), 64)
	results.Subscribe.AvgLatency, _ = strconv.ParseFloat(strings.Replace(temp[11], ",", ".", 1), 64)

	if createLog {
		var files []string

		err := filepath.Walk("output", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Register("Extract failure " + err.Error())
				return nil
			}

			if !info.IsDir() {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			results.Meta.ExperimentError = "Error X: Failed in upload log file"
		}

		for _, f := range files {
			aux := strings.Split(f, "/")
			name := aux[len(aux)-1]
			if name[0:10] == "mqttloader" {
				buffer, _ := ioutil.ReadFile(f)
				results.Meta.LogFile.Data = buffer
				results.Meta.LogFile.Name = name
				results.Meta.LogFile.Extension = strings.Split(name, ".")[1]
				os.Remove(f)
			}
		}
	}

	return results
}

func getToken() (string, bool, bool, bool) {
	var seed rand.Source
	var random *rand.Rand
	var makeRegister bool
	var login_confirmation bool
	var register_confirmation bool
	var token string = "Client_"

	if !utils.FileExists("token.bin") {
		for i := 0; i < 10; i++ {
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			token += fmt.Sprintf("%d", random.Int()%10)
		}
		makeRegister = true
		login_confirmation = true
	} else {
		data, _ := os.ReadFile("token.bin")
		token = strings.Split(string(data), "\n")[0]
		register_confirmation = true
	}

	return token, makeRegister, login_confirmation, register_confirmation
}

func authentication(w *Worker, loginTimeout int, makeRegister bool, login_confirmation bool, register_confirmation bool) {
	token := w.client.Subscribe(w.Id+"/Login/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		login_confirmation = true
	})
	token.Wait()

	token = w.client.Subscribe("Orquestrator/Register/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		response := strings.Split(string(m.Payload()), "-")

		if response[0] != w.Id {
			mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline registration fail", Attr: messages.Command{}})
			token = w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
			token.Wait()
			w.client.Disconnect(0)
			log.Register("Shutdown register failure")
			os.Exit(0)
			return
		}

		var f *os.File

		if !utils.FileExists("token.bin") {
			f, _ = os.Create("token.bin")
		} else {
			f, _ = os.Open("token.bin")
		}
		f.Truncate(0)
		f.Write([]byte(response[1]))

		f.Close()
		register_confirmation = true
	})
	token.Wait()

	if makeRegister {
		token = w.client.Publish("Orquestrator/Register", byte(1), false, w.Id)
		token.Wait()

		cd := 0
		for !register_confirmation {
			if cd >= loginTimeout {
				break
			}
			time.Sleep(time.Second)
			cd++
		}
	} else {
		token = w.client.Publish("Orquestrator/Login", byte(1), false, w.Id)
		token.Wait()

		cd := 0
		log.Register("Waiting login confirmation")
		for !login_confirmation {
			if cd >= loginTimeout {
				break
			}
			time.Sleep(time.Second)
			cd++
		}
	}

	if !login_confirmation {
		mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline login fail", Attr: messages.Command{}})
		token = w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		token.Wait()
		w.client.Disconnect(0)
		log.Register("Shutdown login failure")
		os.Exit(0)
	}

	if !register_confirmation {
		mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline registration fail", Attr: messages.Command{}})
		token = w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		token.Wait()
		w.client.Disconnect(0)
		log.Register("Shutdown register failure")
		os.Exit(0)
	}
	log.Register("Login sucess")
}

func Init(broker string, tool string, loginTimeout int, isUnix bool) {
	log.Create()

	log.Register("Preparing authentication token")
	wtoken, makeRegister, login_confirmation, register_confirmation := getToken()

	log.Register("Configuring mqtt paho client")
	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(wtoken).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(ka).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	worker := Build()
	worker.broker = broker
	worker.tool = tool
	worker.isUnix = isUnix
	worker.client = mqtt.NewClient(opts)
	worker.Id = wtoken

	log.Register("Connecting new mqtt paho client to " + broker + " broker")
	tokenConnection := worker.client.Connect()

	tokenConnection.Wait()

	log.Register("Authenticate token")
	authentication(worker, loginTimeout, makeRegister, login_confirmation, register_confirmation)

	log.Register("Warning orquestrator that worker is online")
	mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "online", Attr: messages.Command{}})
	token := worker.client.Publish(worker.Id+"/Status", byte(1), true, string(mess))
	token.Wait()

	log.Register("Subscribing command topic")
	token = worker.client.Subscribe(worker.Id+"/Command", byte(1), func(c mqtt.Client, m mqtt.Message) {
		if m.Retained() {
			return
		}
		message := m.Payload()

		var commd messages.Command
		err := json.Unmarshal(message, &commd)
		if err != nil {
			mess, _ = json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
			t := worker.client.Publish(worker.Id+"/Status", byte(1), true, string(mess))
			t.Wait()
			log.Register("Crash " + err.Error())
			os.Exit(3)
		}

		switch commd.Name {
		case "info":
			var arguments messages.Info
			jsonObj, _ := json.Marshal(commd.Arguments)
			json.Unmarshal(jsonObj, &arguments)

			go worker.Info(arguments)
		case "start":
			go worker.Start(commd, string(m.Payload()), -1)
		case "cancel":
			experimentListMutex.Lock()
			node := experimentList.Search(int64(commd.Arguments["id"].(float64)))
			if node != nil {
				node.Finished = true
				node.Proc.Kill()
			}
			experimentListMutex.Unlock()
		}
	})

	token.Wait()

	log.Register("Starting keepalive thread")
	go worker.KeepAlive()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Register("Block main thread")
	<-c
	mess, _ = json.Marshal(messages.Status{Type: "Client messages.Status", Status: "offline", Attr: messages.Command{}})
	token = worker.client.Publish(worker.Id+"/Status", byte(1), true, string(mess))
	token.Wait()
	worker.client.Disconnect(0)

	log.Register("Shutdown")
	os.Exit(1)
}
