package communication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/history"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/utils"
)

var logMutex sync.Mutex
var experimentListMutex sync.Mutex
var experimentList history.OngoingExperiments

func loadArguments(file string, arg map[string]interface{}) (bool, int64){
	var arguments messages.CommandExperiment
	jsonObj,_ := json.Marshal(arg)
	json.Unmarshal(jsonObj, &arguments)

	f, _ := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

	var argf string = ""

	isNull := func (param string) string{if param == ""{ return "#"}; return ""}

	argf += fmt.Sprintf("broker=%s\n", arguments.Broker)
	argf += fmt.Sprintf("broker_port=%d\n", arguments.Port)
	argf += fmt.Sprintf("mqtt_version=%d\n", arguments.MqttVersion)
	argf += fmt.Sprintf("num_publishers=%d\n", arguments.NumPublishers)
	argf += fmt.Sprintf("num_subscribers=%d\n", arguments.NumSubscriber)
	argf += fmt.Sprintf("qos_publisher=%d\n", arguments.QosPublisher)
	argf += fmt.Sprintf("qos_subscriber=%d\n", arguments.QosSubscriber)
	argf += fmt.Sprintf("shared_subscription=%t\n", arguments.SharedSubscrition)
	argf += fmt.Sprintf("retain=%t\n", arguments.Retain)
	argf += fmt.Sprintf("topic=%s\n", arguments.Topic)
	argf += fmt.Sprintf("payload=%d\n", arguments.Payload)
	argf += fmt.Sprintf("num_messages=%d\n", arguments.NumMessages)
	argf += fmt.Sprintf("ramp_up=%d\n", arguments.RampUp)
	argf += fmt.Sprintf("ramp_down=%d\n", arguments.RampDown)
	argf += fmt.Sprintf("interval=%d\n", arguments.Interval)
	argf += fmt.Sprintf("subscriber_timeout=%d\n", arguments.SubscriberTimeout)
	argf += fmt.Sprintf("log_level=%s\n", arguments.LogLevel)
	argf += fmt.Sprintf("exec_time=%d\n", arguments.Exec_time)
	argf += fmt.Sprintf("%sntp=%s\n", isNull(arguments.Ntp),arguments.Ntp)
	if arguments.Output {argf += fmt.Sprintf("output=%s\n", "output")}
	argf += fmt.Sprintf("%suser_name=%s\n", isNull(arguments.User),arguments.User)
	argf += fmt.Sprintf("%spassword=%s\n", isNull(arguments.Password),arguments.Password)
	argf += fmt.Sprintf("%stls_truststore=%s\n", isNull(arguments.TlsTrustsore), arguments.TlsTrustsore)
	argf += fmt.Sprintf("%stls_truststore_pass=%s\n", isNull(arguments.TlsTruststorePassword), arguments.TlsTruststorePassword)
	argf += fmt.Sprintf("%stls_keystore=%s\n", isNull(arguments.TlsKeystore),arguments.TlsKeystore)
	argf += fmt.Sprintf("%stls_keystore_pass=%s\n", isNull(arguments.TlsKeystorePassword),arguments.TlsKeystorePassword)

    // Write bytes to file
    byteSlice := []byte(argf)
    _, err := f.Write(byteSlice)
    if err != nil {
		logMutex.Lock()
		f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("crash "+err.Error()+"\n")
		logMutex.Unlock()
		f.Close()
    }

	return arguments.Output,arguments.Expid
}

func extracExperimentResults(output string, createLog bool) messages.ExperimentResult{
	results:= messages.ExperimentResult{}
	results.Meta.Literal = output

	temp := [12]string{}

	i := 0

	for _,s := range strings.Split(output, "\n"){
		if s != "" && s[0] != '-'{
			temp[i] = strings.Split(s, ": ")[1]
			i++
		}
	}
	results.Meta.ToolName = "mqttLoader"
	results.Meta.ExperimentError = ""
	
	results.Publish.Throughput,_ = strconv.ParseFloat(strings.Replace(temp[2], ",",".", 1), 64) 
	results.Publish.AvgThroughput,_ = strconv.ParseFloat(strings.Replace(temp[3], ",",".", 1), 64)
	results.Publish.PubMessages,_ = strconv.Atoi(temp[4])

	for _,s := range strings.Split(temp[5], ", "){
		aux,_ := strconv.Atoi(s)
		results.Publish.PerSecondThrouput = append(results.Publish.PerSecondThrouput, aux) 
	}

	results.Subscribe.Throughput,_ = strconv.ParseFloat(strings.Replace(temp[6], ",",".", 1), 64)
	results.Subscribe.AvgThroughput,_ = strconv.ParseFloat(strings.Replace(temp[7], ",",".", 1), 64)
	results.Subscribe.ReceivedMessages,_ = strconv.Atoi(temp[8])

	for _,s := range strings.Split(temp[9], ", "){
		aux,_ := strconv.Atoi(s)
		results.Subscribe.PerSecondThrouput = append(results.Subscribe.PerSecondThrouput, aux) 
	}

	results.Subscribe.Latency,_ = strconv.ParseFloat(strings.Replace(temp[10], ",",".", 1), 64)
	results.Subscribe.AvgLatency,_ = strconv.ParseFloat(strings.Replace(temp[11], ",",".", 1), 64)

	if createLog {
		var files []string

		err := filepath.Walk("output", func(path string, info os.FileInfo, err error) error {
			if err != nil{
				logMutex.Lock()
				f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				f.WriteString("extract failure "+err.Error()+"\n")
				f.Close()
				logMutex.Unlock()
				return nil
			}

			if !info.IsDir(){
				files = append(files, path)
			}

			return nil
		})

		if err != nil{
			results.Meta.ExperimentError = "Error X: Failed in upload log file"
		}

		for _,f := range files{
			aux := strings.Split(f, "/")
			name := aux[len(aux) - 1]
			if name[0:10] == "mqttloader"{
				buffer,_ := ioutil.ReadFile(f)
				results.Meta.LogFile.Data = buffer
				results.Meta.LogFile.Name = name
				results.Meta.LogFile.Extension = strings.Split(name, ".")[1]
				os.Remove(f)
			}
		}
	}

	return results
}

func workerKeepAlive(client mqtt.Client, msg string){
	for {
		Ping(client, msg)
		time.Sleep(time.Second)
	}
}

func Init(broker string, tool string,loginTimeout int, isUnix bool) {
	var clientID string = "Client_"
	var seed rand.Source
	var random *rand.Rand
	var makeRegister bool = false
	var login_confirmation bool = false
	var register_confirmation bool = false

	if !utils.FileExists("worker.log"){
		f,_ := os.Create("worker.log")
		f.Close()
	} else{
		os.Truncate("worker.log", 0)
	}

	if !utils.FileExists("token.bin"){
		for i := 0; i < 10; i++{
			seed = rand.NewSource(time.Now().UnixNano())
			random = rand.New(seed)
			clientID += fmt.Sprintf("%d", random.Int() % 10)
		}
		makeRegister = true
		login_confirmation = true
	} else{
		data,_ := os.ReadFile("token.bin")
		clientID = strings.Split(string(data), "\n")[0]
		register_confirmation = true
	}

	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")

	opts := mqtt.NewClientOptions().
			AddBroker(broker).
			SetClientID(clientID).
			SetCleanSession(true).
			SetAutoReconnect(true).
			SetKeepAlive(ka).
			SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
			SetConnectionLostHandler(func(client mqtt.Client, reason error) {})
	
	client := mqtt.NewClient(opts)

	tokenConnection := client.Connect()

	tokenConnection.Wait()

	token := client.Subscribe(clientID+"/Login/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		login_confirmation = true
	})
	token.Wait()

	token = client.Subscribe("Orquestrator/Register/Log", byte(1), func(c mqtt.Client, m mqtt.Message) {
		response := strings.Split(string(m.Payload()), "-")

		if response[0] != clientID {
			mess,_ := json.Marshal(messages.Status{Type:"Client Status", Status: "offline registration fail", Attr: messages.Command{}})
			token = client.Publish(clientID+"/Status", byte(1), true, string(mess))
			token.Wait()
			client.Disconnect(0)
			logMutex.Lock()
			f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			f.WriteString("shutdown register failure\n")
			f.Close()
			logMutex.Unlock()
			os.Exit(0)
			return
		}

		var f *os.File

		if!utils.FileExists("token.bin"){
			f,_ = os.Create("token.bin")
		} else{
			f,_ = os.Open("token.bin")
		}
		f.Truncate(0)
		f.Write([]byte(response[1]))

		f.Close()
		register_confirmation = true
	})
	token.Wait()

	if makeRegister {
		token = client.Publish("Orquestrator/Register", byte(1), false, clientID)
		token.Wait()

		cd := 0
		for !register_confirmation{
			if cd >= loginTimeout{
				break
			}
			time.Sleep(time.Second)
			cd++
		}
	}else{
		token = client.Publish("Orquestrator/Login", byte(1), false, clientID)
		token.Wait()

		cd := 0
		for !login_confirmation{
			if cd >= loginTimeout{
				break
			}
			time.Sleep(time.Second)
			cd++
		}
	}

	if !login_confirmation{
		mess,_ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline login fail",Attr: messages.Command{}})
		token = client.Publish(clientID+"/Status", byte(1), true, string(mess))
		token.Wait()
		client.Disconnect(0)
		logMutex.Lock()
		f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("shutdown login failure\n")
		f.Close()
		logMutex.Unlock()
		os.Exit(0)
	}

	if !register_confirmation{
		mess,_ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline registration fail", Attr: messages.Command{}})
		token = client.Publish(clientID+"/Status", byte(1), true, string(mess))
		token.Wait()
		client.Disconnect(0)
		logMutex.Lock()
		f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("shutdown register failure\n")
		f.Close()
		logMutex.Unlock()
		os.Exit(0)
	}

	mess,_ := json.Marshal(messages.Status{Type: "Client Status", Status: "online", Attr: messages.Command{}})
	token = client.Publish(clientID+"/Status", byte(1), true, string(mess))
	token.Wait()

	token = client.Subscribe(clientID+"/Command", byte(1), func(c mqtt.Client, m mqtt.Message) {
		if m.Retained(){
			return
		}
		message := m.Payload()

		var commd messages.Command
		err := json.Unmarshal(message, &commd)
		if err != nil {
			mess,_ = json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
			t := client.Publish(clientID+"/Status", byte(1), true, string(mess))
			t.Wait()
			logMutex.Lock()
			f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			f.WriteString("crash "+err.Error()+"\n")
			f.Close()
			logMutex.Unlock()
			os.Exit(3)
		}

		switch commd.Name{
			case "info":
				var arguments messages.Info
				jsonObj,_ := json.Marshal(commd.Arguments)
				json.Unmarshal(jsonObj, &arguments)

				go Info(client, arguments, isUnix, clientID)
			case "start":
				go Start(client, clientID, tool, commd, string(m.Payload()), -1)
			case "cancel":
				experimentListMutex.Lock()
				node := experimentList.Search(int64(commd.Arguments["id"].(float64)))
				if node != nil{
					node.Finished = true
					node.Proc.Kill()
				}
				experimentListMutex.Unlock()
		}
	})

	token.Wait()

	token = client.Subscribe(clientID+"/Ping", byte(1), func(c mqtt.Client, m mqtt.Message) {
		Ping(client, clientID)
	})

	token.Wait()

	go workerKeepAlive(client, clientID)

	c := make(chan os.Signal, 1)

    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<- c
	mess,_ = json.Marshal(messages.Status{Type: "Client messages.Status", Status: "offline", Attr: messages.Command{}})
	token = client.Publish(clientID+"/Status", byte(1), true, string(mess))
	token.Wait()
	client.Disconnect(0)

	logMutex.Lock()
	f,_ := os.OpenFile("worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString("shutdown\n")
	f.Close()
	logMutex.Unlock()
	os.Exit(1)
}