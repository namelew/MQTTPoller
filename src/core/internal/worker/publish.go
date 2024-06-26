package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/namelew/mqtt-poller/src/core/internal/worker/history"
	"github.com/namelew/mqtt-poller/src/core/packages/messages"
	"github.com/namelew/mqtt-poller/src/core/packages/utils"
)

func (w *Worker) Ping() {
	t := w.client.Publish(w.Id+"/KeepAlive", byte(1), false, w.Id)
	t.Wait()
}

func (w *Worker) Start(cmdExp messages.Command, commandLiteral string, experimentId int64) {
	var arg_file string = `"myconfig.conf"`
	var flag string = "-c"
	var id int64

	if !utils.FileExists("myconfig.conf") {
		file, err := os.Create("myconfig.conf")

		if err != nil {
			log.Register("Can't create mqttloader config file")
			mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
			t := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
			t.Wait()
			os.Exit(3)
		}

		file.Close()
	}
	err := os.Truncate("myconfig.conf", 0)

	if err != nil {
		log.Register("Can't to truncate mqttloader config file")
		mess, _ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
		t := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		t.Wait()
		os.Exit(3)
	}

	if cmdExp.Arguments != nil {
		id = loadArguments("myconfig.conf", cmdExp.Arguments)
	} else {
		arg_file = ""
		flag = ""
	}

	if experimentId != -1 {
		id = experimentId
	}

	if !utils.FileExists("CommandsLog/experiment_" + fmt.Sprint(id) + ".json") {
		file, _ := os.Create("CommandsLog/experiment_" + fmt.Sprint(id) + ".json")
		file.Write([]byte(commandLiteral))
		file.Close()
	}

	log.Register("Start experiment " + strconv.FormatInt(id, 10))

	cmd := exec.Command(w.tool, flag, arg_file)

	mess, _ := json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: "start", Attr: cmdExp})
	t := w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
	t.Wait()
	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr
	err = cmd.Start()

	experimentNode := history.CreateRegister(id, cmd.Process)

	experimentListMutex.Lock()
	experimentList.Add(&experimentNode)
	experimentListMutex.Unlock()

	declaration, ok := cmdExp.Arguments["declaration"].(map[string]interface{})

	if !ok {
		log.Register("Unabel to read experiment declaration from orquestrator")
	}

	execTime, ok := declaration["execTime"].(float64)

	if !ok {
		log.Register("Unabel to read experiment execution time from orquestrator")
	}

	if execTime == 0 {
		execTime = 1
	}

	go func(expid int64, timeout time.Duration) {
		<-time.After(timeout)
		experimentListMutex.Lock()

		node := experimentList.Search(expid)

		if node != nil && !node.Finished {
			log.Register(fmt.Sprintf("Experiment %d excced the execution timeout", expid))
			node.Finished = true
			node.Proc.Kill()
		}

		experimentListMutex.Unlock()
	}(id, time.Second*time.Duration(execTime)*3)

	cmd.Wait()

	experimentListMutex.Lock()
	if experimentNode.Finished {
		experimentList.Remove(id)
		experimentListMutex.Unlock()
		return
	}
	experimentList.Remove(id)
	experimentListMutex.Unlock()

	if err != nil {
		mess, _ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: fmt.Sprint(err) + ": " + stderr.String(), Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		mess, _ = json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Crash experiment " + strconv.FormatInt(id, 10) + " error " + err.Error())
		os.Exit(3)
	}

	resultsExperiment := extracExperimentResults(output.String(), stderr.String())

	if resultsExperiment.Meta.ExperimentError != "" {
		mess, _ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: resultsExperiment.Meta.ExperimentError, Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Error experiment: " + resultsExperiment.Meta.ExperimentError)
	} else {
		mess, _ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: "finish", Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Finish experiment " + strconv.FormatInt(id, 10))
	}

	resultsExperiment.Meta.ID = uint64(id)

	results, err := json.Marshal(resultsExperiment)

	if err != nil {
		log.Register("Unable to marshall experiment result from " + strconv.FormatInt(id, 10) + "." + err.Error())
	}

	t = w.client.Publish(w.Id+"/Experiments/Results", byte(1), false, string(results))
	t.Wait()

	if t.Error() != nil {
		log.Register("Unable to send experiment result from " + strconv.FormatInt(id, 10) + "." + t.Error().Error())
	}

	os.Remove("CommandsLog/experiment_" + fmt.Sprint(id) + ".json")
}

func (w *Worker) KeepAlive() {
	for {
		w.Ping()
		time.Sleep(time.Second)
	}
}
