package worker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/namelew/mqtt-bm-latency/internal/worker/history"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
	"github.com/namelew/mqtt-bm-latency/packages/utils"
)

type Worker struct {
	Id             string
	tool           string
	broker         string
	loginTimeout   time.Duration
	loginThreshold int
	client         mqtt.Client
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

func sanitizeStrings(s string) string {
	garbage := []string{"\n", "\t", "\r", "\a", "\f", "\v", "\b", " "}

	for _, str := range garbage {
		s = strings.ReplaceAll(s, str, "")
	}

	return s
}

func extracExperimentResults(output string, logs string, createLog bool) messages.ExperimentResult {
	results := messages.ExperimentResult{}
	results.Meta.Literal = logs + output
	results.Meta.ToolName = "mqttLoader"

	if output == "" {
		results.Meta.ExperimentError = "Tool runtime error"
		return results
	}

	temp := [12]string{}
	var err error

	i := 0

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			data := strings.Split(sanitizeStrings(line), ":")
			if len(data) < 2 || data[0] == "ERROR" {
				log.Register("Experiment Error: " + output)
				results.Meta.ExperimentError = output
				return results
			}
			temp[i] = data[1]
			i++
		}
	}
	results.Meta.ExperimentError = ""

	results.Publish.Throughput, err = strconv.ParseFloat(strings.Replace(temp[2], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	results.Publish.AvgThroughput, err = strconv.ParseFloat(strings.Replace(temp[3], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	results.Publish.PubMessages, err = strconv.Atoi(temp[4])

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	for _, s := range strings.Split(temp[5], ",") {
		aux, err := strconv.Atoi(s)

		if err != nil {
			errorMessage := "Unable to parse data. " + err.Error()
			log.Register("Experiment Error: " + errorMessage)
			results.Meta.ExperimentError = errorMessage
			return results
		}

		results.Publish.PerSecondThrouput = append(results.Publish.PerSecondThrouput, aux)
	}

	results.Subscribe.Throughput, err = strconv.ParseFloat(strings.Replace(temp[6], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	results.Subscribe.AvgThroughput, err = strconv.ParseFloat(strings.Replace(temp[7], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	results.Subscribe.ReceivedMessages, err = strconv.Atoi(temp[8])

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	for _, s := range strings.Split(temp[9], ",") {
		aux, err := strconv.Atoi(s)

		if err != nil {
			errorMessage := "Unable to parse data. " + err.Error()
			log.Register("Experiment Error: " + errorMessage)
			results.Meta.ExperimentError = errorMessage
			return results
		}

		results.Subscribe.PerSecondThrouput = append(results.Subscribe.PerSecondThrouput, aux)
	}

	results.Subscribe.Latency, err = strconv.ParseFloat(strings.Replace(temp[10], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

	results.Subscribe.AvgLatency, err = strconv.ParseFloat(strings.Replace(temp[11], ",", ".", 1), 64)

	if err != nil {
		errorMessage := "Unable to parse data. " + err.Error()
		log.Register("Experiment Error: " + errorMessage)
		results.Meta.ExperimentError = errorMessage
		return results
	}

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
				buffer, err := os.ReadFile(f)

				if err != nil {
					errorMessage := "Unable to read data from output file. " + err.Error()
					log.Register("Experiment Error: " + errorMessage)
					results.Meta.ExperimentError = errorMessage
					return results
				}

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
	token := uuid.New().String()

	makeRegister := !utils.FileExists("token.bin")

	if !makeRegister {
		data, err := os.ReadFile("token.bin")

		if err != nil {
			log.Register("Unable to open token file. " + err.Error())
			return token, true
		}

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

	go func() {
		<-time.After(w.loginTimeout)
		confirmation <- false
	}()

	if !<-confirmation {
		if threshold > 1 {
			authentication(w, threshold-1)
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

func createClient(w *Worker, opts *mqtt.ClientOptions) mqtt.Client {
	log.Register("Configuring mqtt paho client")

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
		response := strings.Split(string(m.Payload()), " ")

		if response[0] != w.Id {
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

	go func() {
		<-time.After(w.loginTimeout)
		confirmation <- false
	}()

	if !<-confirmation {
		if threshold > 1 {
			register(w, threshold-1)
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

func connect(w *Worker) {
	log.Register("Check if token exists")

	wtoken, makeRegister := getToken()

	w.Id = wtoken
	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")
	opts := mqtt.NewClientOptions().
		AddBroker(w.broker).
		SetClientID(w.Id).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(ka).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			log.Register("Connection lost. Reason: " + reason.Error())
		})

	if makeRegister {
		w.client = createClient(w, opts)
		register(w, w.loginThreshold)
	}

	opts.SetCleanSession(false)

	w.client = createClient(w, opts)

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

	worker := Worker{
		broker:         broker,
		tool:           tool,
		loginTimeout:   time.Duration(loginTimeout) * time.Second,
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
			expid, ok := commd.Arguments["id"].(float64)

			if !ok {
				log.Register("Unabel to read experiment id from orquestrator")
			}

			experimentListMutex.Lock()

			node := experimentList.Search(int64(expid))

			if node != nil {
				log.Register(fmt.Sprintf("Canceling experiment %d", int64(expid)))
				node.Finished = true
				node.Proc.Kill()
			} else {
				log.Register(fmt.Sprintf("Unable to find experiment %d, cancel operation fail", int64(expid)))
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
