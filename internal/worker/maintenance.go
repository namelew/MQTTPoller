package worker

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
	"github.com/namelew/mqtt-bm-latency/internal/worker/history"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
	"github.com/namelew/mqtt-bm-latency/packages/utils"
)

type Worker struct {
	Id     string
	tool   string
	broker string
	loginTimeout time.Duration
	loginThreshold int
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

	argf += fmt.Sprintf("%sbroker=%s\n", isNull(arguments.Declaration.Broker), arguments.Declaration.Broker)
	argf += fmt.Sprintf("%sbroker_port=%d\n", isNull(arguments.Declaration.Port), arguments.Declaration.Port)
	argf += fmt.Sprintf("%smqtt_version=%d\n", isNull(arguments.Declaration.MqttVersion), arguments.Declaration.MqttVersion)
	argf += fmt.Sprintf("%snum_publishers=%d\n", isNull(arguments.Declaration.NumPublishers), arguments.Declaration.NumPublishers)
	argf += fmt.Sprintf("%snum_subscribers=%d\n", isNull(arguments.Declaration.NumSubscriber), arguments.Declaration.NumSubscriber)
	argf += fmt.Sprintf("%sqos_publisher=%d\n", isNull(arguments.Declaration.QosPublisher), arguments.Declaration.QosPublisher)
	argf += fmt.Sprintf("%sqos_subscriber=%d\n", isNull(arguments.Declaration.QosSubscriber), arguments.Declaration.QosSubscriber)
	argf += fmt.Sprintf("%sshared_subscription=%t\n", isNull(arguments.Declaration.SharedSubscrition), arguments.Declaration.SharedSubscrition)
	argf += fmt.Sprintf("%sretain=%t\n", isNull(arguments.Declaration.Retain), arguments.Declaration.Retain)
	argf += fmt.Sprintf("%stopic=%s\n", isNull(arguments.Declaration.Topic), arguments.Declaration.Topic)
	argf += fmt.Sprintf("%spayload=%d\n", isNull(arguments.Declaration.Payload), arguments.Declaration.Payload)
	argf += fmt.Sprintf("%snum_messages=%d\n", isNull(arguments.Declaration.NumMessages), arguments.Declaration.NumMessages)
	argf += fmt.Sprintf("%sramp_up=%d\n", isNull(arguments.Declaration.RampUp), arguments.Declaration.RampUp)
	argf += fmt.Sprintf("%sramp_down=%d\n", isNull(arguments.Declaration.RampDown), arguments.Declaration.RampDown)
	argf += fmt.Sprintf("%sinterval=%d\n", isNull(arguments.Declaration.Interval), arguments.Declaration.Interval)
	argf += fmt.Sprintf("%ssubscriber_timeout=%d\n", isNull(arguments.Declaration.SubscriberTimeout), arguments.Declaration.SubscriberTimeout)
	argf += fmt.Sprintf("%slog_level=%s\n", isNull(arguments.Declaration.LogLevel), arguments.Declaration.LogLevel)
	argf += fmt.Sprintf("%sexec_time=%d\n", isNull(arguments.Declaration.ExecTime), arguments.Declaration.ExecTime)
	argf += fmt.Sprintf("%sntp=%s\n", isNull(arguments.Declaration.Ntp), arguments.Declaration.Ntp)
	if arguments.Declaration.Output {
		argf += fmt.Sprintf("output=%s\n", "output")
	}
	argf += fmt.Sprintf("%suser_name=%s\n", isNull(arguments.Declaration.User), arguments.Declaration.User)
	argf += fmt.Sprintf("%spassword=%s\n", isNull(arguments.Declaration.Password), arguments.Declaration.Password)
	argf += fmt.Sprintf("%stls_truststore=%s\n", isNull(arguments.Declaration.TlsTrustsore), arguments.Declaration.TlsTrustsore)
	argf += fmt.Sprintf("%stls_truststore_pass=%s\n", isNull(arguments.Declaration.TlsTruststorePassword), arguments.Declaration.TlsTruststorePassword)
	argf += fmt.Sprintf("%stls_keystore=%s\n", isNull(arguments.Declaration.TlsKeystore), arguments.Declaration.TlsKeystore)
	argf += fmt.Sprintf("%stls_keystore_pass=%s\n", isNull(arguments.Declaration.TlsKeystorePassword), arguments.Declaration.TlsKeystorePassword)

	// Write bytes to file
	byteSlice := []byte(argf)
	_, err := f.Write(byteSlice)
	if err != nil {
		log.Register("Load arguments error " + err.Error())
		f.Close()
	}

	return arguments.Declaration.Output, arguments.Expid
}

func extracExperimentResults(output string, createLog bool) messages.ExperimentResult {
	results := messages.ExperimentResult{}
	results.Meta.Literal = output

	temp := [12]string{}

	i := 0
	for _, s := range strings.Split(output, "\n") {
		if s != "" && s[0] != '-' {
			data := strings.Split(s, ": ")
			if data[0] == "ERROR" {
				results.Meta.ExperimentError = data[1]
				return results
			}
			temp[i] = data[1]
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

func getToken() (string, bool) {
	var seed rand.Source
	var random *rand.Rand
	var token string = "Client_"

	makeRegister := !utils.FileExists("token.bin")

	if makeRegister {
		for i := 0; i < 10; i++ {
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			token += fmt.Sprintf("%d", random.Int()%10)
		}
	} else {
		data, _ := os.ReadFile("token.bin")
		token = strings.Split(string(data), "\n")[0]
	}

	return token, makeRegister
}

func authentication(w *Worker, threshold int) {
	confirmation := make(chan bool, 1)

	token := w.client.Subscribe(w.Id+"/Login/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		confirmation <- true
	})
	token.Wait()

	token = w.client.Publish("Orquestrator/Login", byte(1), false, w.Id)
	token.Wait()

	go func ()  {
		<-time.After(w.loginTimeout)
		confirmation <- false	
	}()

	if !<-confirmation {
		if threshold > 1 {
			authentication(w, threshold - 1)
			return
		}

		if w.loginThreshold < 0 {
			authentication(w, threshold)
			return
		}
		
		mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline auth fail", Attr: messages.Command{}})
		token := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		token.Wait()
		w.client.Disconnect(0)
		log.Register("Shutdown auth failure")
		os.Exit(0)
	}

	log.Register("Auth sucess")
}

func createClient(w *Worker) mqtt.Client{
	log.Register("Configuring mqtt paho client")
	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")

	opts := mqtt.NewClientOptions().
		AddBroker(w.broker).
		SetClientID(w.Id).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(ka).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	c := mqtt.NewClient(opts)

	log.Register("Connecting new mqtt paho client to " + w.broker + " broker")
	token := c.Connect()
	token.Wait()

	return c
}

func register(w *Worker, threshold int) {
	log.Register("Registering worker to orquestrator")

	confirmation := make(chan bool, 1)

	token := w.client.Subscribe("Orquestrator/Register/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		response := strings.Split(string(m.Payload()), "-")

		if response[0] != w.Id {
			mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline registration fail", Attr: messages.Command{}})
			token := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
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
		w.Id = response[1]
		confirmation <- true
	})
	token.Wait()

	token = w.client.Publish("Orquestrator/Register", byte(1), false, w.Id)
	token.Wait()

	go func ()  {
		<-time.After(w.loginTimeout)
		confirmation <- false	
	}()

	if !<-confirmation {
		if threshold > 1 {
			register(w, threshold - 1)
			return
		}

		if w.loginThreshold < 0 {
			register(w, threshold)
			return
		}
		
		mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline auth fail", Attr: messages.Command{}})
		token := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		token.Wait()
		w.client.Disconnect(0)
		log.Register("Shutdown register failure")
		os.Exit(0)
	}

	log.Register("Register sucess")

	w.client.Disconnect(0)
}

func connect(w *Worker){
	log.Register("Check if token exists")

	wtoken, makeRegister := getToken()

	if makeRegister {
		w.Id = wtoken
		w.client = createClient(w)
		register(w, w.loginThreshold)	
	}

	w.client = createClient(w)
	log.Register("Authentication worker")
	authentication(w, w.loginThreshold)
}

func disconnect(worker *Worker) {
	mess, _ := json.Marshal(messages.Status{Type: "Client messages.Status", Status: "offline", Attr: messages.Command{}})
	token := worker.client.Publish(worker.Id+"/Status", byte(1), true, string(mess))
	token.Wait()
	worker.client.Disconnect(0)

	log.Register("Shutdown")
	os.Exit(1)
}

func Init(broker string, tool string, loginTimeout, loginThreshold int) {
	log.Create()

	log.Register("Preparing authentication token")

	worker := Worker {
		broker: broker,
		tool: tool,
		loginTimeout: time.Duration(loginTimeout) * time.Second,
		loginThreshold: loginThreshold,
	}
	
	connect(&worker)

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

	disconnect(&worker)
}
