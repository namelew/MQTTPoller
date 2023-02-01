package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var workers = make([]worker, 1, 10)

func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func absInt(x int) int {
	if x < 0 {
		return (x * -1)
	}
	return x
}

func getJsonFromFile(file string, expid int64) (string, int, command, int) {
	if !fileExists(file) {
		return "", 0, command{}, 0
	}

	var exec_time int
	var attemps int
	data, _ := os.ReadFile(file)

	var temp command
	var jsonArg commandExperiment

	json.Unmarshal(data, &temp)

	data, _ = json.Marshal(temp.Arguments)

	json.Unmarshal(data, &jsonArg)

	jsonArg.Expid = expid
	exec_time = jsonArg.Exec_time
	attemps = jsonArg.Attempts

	data, _ = json.Marshal(jsonArg)

	json.Unmarshal(data, &temp.Arguments)

	data, _ = json.Marshal(temp)

	return string(data), exec_time, temp, attemps
}

func setMessageHandler(client mqtt.Client, id int) {
	token := client.Subscribe(workers[id].Id+"/Experiments/Results", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperiment(m, id)
	})
	token.Wait()
	token = client.Subscribe(workers[id].Id+"/Experiments/Status", byte(1), func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperimentStatus(client, m, id)
	})
	token.Wait()
}

func messageHandlerExperimentStatus(client mqtt.Client, msg mqtt.Message, id int) {
	var exps status
	json.Unmarshal(msg.Payload(), &exps)

	tokens := strings.Split(exps.Type, " ")
	expid, _ := strconv.Atoi(tokens[2])

	exp := workers[id].historic.search(int64(expid))

	switch exps.Status {
	case "start":
		return
	case "finish":
		if exp != nil {
			exp.finished = true
		}
	default:
		if exp != nil {
			exp.finished = true
			exp.err = true
			redoExperiment(client, id, exp)
		}
	}
}

func isIn(workers []worker, clientID string) bool {
	for _, w := range workers {
		if clientID == w.Id {
			return true
		}
	}
	return false
}

func receiveControl(client mqtt.Client, id int, timeout int) {
	var start int64 = time.Now().UnixMilli()
	exlog := workers[id].historic.root.findLarger()

	for absInt(int(time.Now().UnixMilli()-start)) < (timeout * 1000) {
		if workers[id].ReceiveConfirmation || !workers[id].Status || exlog.err {
			break
		}
	}

	workers[id].ReceiveConfirmation = false

	if (timeout*1000) <= absInt(int(time.Now().UnixMilli()-start)) || !workers[id].Status || exlog.err {
		log.Printf("Error in worker %d: experiment don't return\n", id)
		exlog.finished = true
		redoExperiment(client, id, exlog)
	}
}

func watcher(client mqtt.Client, id int, tl int) {
	var start int64 = time.Now().UnixMilli()

	for absInt(int(time.Now().UnixMilli()-start)) < (tl * 1000) {
		if workers[id].TestPing {
			return
		}
	}
	workers[id].Status = false
	workers[id].TestPing = true

	token := client.Unsubscribe(workers[id].Id + "/Experiments/Results")
	token.Wait()

	log.Printf("Worker %d is off\n", id)
}

func messageHandlerExperiment(m mqtt.Message, id int) {
	var output experimentResult

	worker := workers[id].historic.root.findLarger()

	if worker != nil {
		worker.finished = true
	}

	workers[id].ReceiveConfirmation = true

	json.Unmarshal(m.Payload(), &output)

	if output.Meta.LogFile.Name != "" {
		ioutil.WriteFile(workers[id].Id+output.Meta.LogFile.Name, output.Meta.LogFile.Data, 0644)
	}

	log.Printf("Experiment %d in worker %d return\n", output.Meta.ID, id)
}

func messageHandlerInfos(m mqtt.Message, id int) {
	var output infoDisplay

	json.Unmarshal(m.Payload(), &output)
	workers[id].ReceiveConfirmation = true

	//fmt.Printf("ID: %d\n", id)
	//fmt.Printf("CPU: %s\n", output.Cpu)
	//fmt.Printf("RAM: %d\n", output.Ram)
	//fmt.Printf("Storage: %d\n", output.Disk)
}

func startExperiment(client mqtt.Client, session *session, arg start) {
	expid := time.Now().Unix()
	msg, exec_t, cmd, attemps := getJsonFromFile(arg.JsonArg, expid)

	if arg.Id[0] == -1 {
		for i := 0; i < len(workers); i++ {
			if !workers[i].Status {
				log.Printf("Worker %d is off, skipping\n", i)
				continue
			}

			workers[i].historic.Add(expid, cmd, attemps)

			token := client.Publish(workers[i].Id+"/Command", byte(1), false, msg)
			token.Wait()

			go receiveControl(client, i, exec_t*5)

			log.Printf("Requesting experiment in worker %d\n", i)
		}
	} else {
		argTam := len(arg.Id)
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

			workers[i].historic.Add(expid, cmd, attemps)

			token := client.Publish(workers[arg.Id[i]].Id+"/Command", byte(1), false, msg)
			token.Wait()

			go receiveControl(client, arg.Id[i], exec_t*5)

			log.Printf("Requesting experiment in worker %d\n", arg.Id[i])
		}
	}
}

func cancelExperiment(client mqtt.Client, id int, expid int64) {
	arg := make(map[string]interface{})
	arg["id"] = expid
	cmd := command{"cancel", "moderation command", arg}
	msg, _ := json.Marshal(cmd)

	token := client.Publish(workers[id].Id+"/Command", byte(1), false, msg)
	token.Wait()

	workers[id].ReceiveConfirmation = true
	exp := workers[id].historic.search(expid)
	exp.finished = true
}

func redoExperiment(client mqtt.Client, worker int, experiment *experimentLog) {
	exp := *experiment
	workers[worker].historic.remove(experiment.id)

	if len(workers) == 1 {
		return
	}

	if exp.attempts > 0 {
		exp.attempts--
		size := len(workers)
		var sample = make([]int, 0, size)
		var timeout int

		for i := 0; i < size; i++ {
			if i != worker && workers[i].Status {
				sample = append(sample, i)
			}
		}

		cmdExp := exp.cmd.ToCommandExperiment()
		exp.id = time.Now().Unix()
		cmdExp.Expid = exp.id
		cmdExp.Attempts = exp.attempts
		timeout = cmdExp.Exec_time * 5
		cmdExp.Attach(&exp.cmd)

		nw := sample[rand.Intn(len(sample))]

		msg, _ := json.Marshal(exp.cmd)

		workers[nw].historic.Add(exp.id, exp.cmd, exp.attempts)

		token := client.Publish(workers[nw].Id+"/Command", byte(1), false, msg)
		token.Wait()

		go receiveControl(client, nw, timeout)
	}
}

func getInfo(client mqtt.Client, arg infoTerminal) {
	var infoCommand command

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
}

func retWorker(c echo.Context) error {
	type workerJson struct{
		Id  int
		NetId string
		Online bool
		History []interface{}
	}
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	switch (reflect.TypeOf(json_map["wid"])).Name(){
	case "float64":
		tempid,ok := json_map["wid"].(float64)
		if !ok {
			return  echo.ErrInternalServerError
		}

		wid := int(tempid)
		
		temp_hist := make([]interface{}, 1)
		workers[wid].historic.Print(temp_hist)
		response := workerJson{wid, workers[wid].Id, workers[wid].Status, temp_hist}

		return c.JSON(200, response)

	case "string":
		for i:=0; i < len(workers); i++{
			
		}
	default:
		response := make([]workerJson, len(workers))
		for i:=0; i < len(workers); i++{
			temp_hist := make([]interface{}, 1)
			workers[i].historic.Print(temp_hist)
			wj := workerJson{i, workers[i].Id, workers[i].Status, temp_hist}
			if response[0].NetId == ""{
				response[0] = wj
			} else{
				response = append(response, wj)
			}
		}
		return c.JSON(200,response)
	}
	return nil
}

func retInfo(c echo.Context) error {
	return nil
}

func main() {
	var (
		broker     = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		t_interval = flag.Int("tl", 5, "orquestrator tolerance interval")
	)
	flag.Parse()
	var currentSession session
	currentSession.Finish = true
	var clientID string = "Orquestrator"
	ka, _ := time.ParseDuration(strconv.Itoa(10000) + "s")

	if !fileExists("orquestrator.log") {
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

	client := mqtt.NewClient(opts)

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
			workers[0] = worker{clientID, true, false, true, experimentHistory{nil}}

			t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
			t.Wait()

			setMessageHandler(client, 0)
			return
		}
		workers = append(workers, worker{clientID, true, false, true, experimentHistory{nil}})

		t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
		t.Wait()

		setMessageHandler(client, len(workers)-1)
	})
	token.Wait()

	token = client.Subscribe("Orquestrator/Login", byte(1), func(c mqtt.Client, m mqtt.Message) {
		f, _ := os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString("worker " + string(m.Payload()) + " login\n")
		f.Close()
		if !isIn(workers, string(m.Payload())) {
			if workers[0].Id == "" {
				workers[0] = worker{string(m.Payload()), true, false, true, experimentHistory{nil}}

				t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
				t.Wait()

				setMessageHandler(client, 0)
				return
			}
			workers = append(workers, worker{string(m.Payload()), true, false, true, experimentHistory{nil}})

			t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
			t.Wait()

			setMessageHandler(client, len(workers)-1)
		} else {
			id := 0
			for i := 0; i < len(workers); i++ {
				if workers[i].Id == string(m.Payload()) {
					id = i
					break
				}
			}
			workers[id].Status = true
			setMessageHandler(client, id)
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
			go watcher(client, id, *t_interval)
		} else {
			workers[id].TestPing = true
			go watcher(client, id, *t_interval)
			workers[id].TestPing = false
		}
	})

	token.Wait()

	api := echo.New()
	api.GET("/orquestrator/worker", retWorker)
	api.GET("/orquestrator/info", retInfo)
	api.POST("/orquestrator/experiment/start", nil)
	api.POST("/orquestrator/experiment/cancel", nil)
	api.Logger.Fatal(api.Start(":8080"))

	f, _ = os.OpenFile("orquestrator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString("disconnect mqtt.client\n")

	client.Disconnect(0)

	f.WriteString("shutdown")
	f.Close()
}
