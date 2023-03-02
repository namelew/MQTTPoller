package communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
	"github.com/namelew/mqtt-bm-latency/history"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func (w *Worker)Ping(){
	log.Register("Sending ka message to orquestrator")
	t := w.client.Publish("Orquestrator/Ping", byte(1), false, w.Id)
	t.Wait()
}

func (w *Worker) Start(cmdExp messages.Command, commandLiteral string, experimentId int64){
	var arg_file string = `"myconfig.conf"`
	var flag string = "-c"
	var createLogFile bool = false
	var id int64

	if !utils.FileExists("myconfig.conf"){
		file,err := os.Create("myconfig.conf")

		if err != nil {
			log.Register("Can't create mqttloader config file")
			mess,_ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
			t := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
			t.Wait()
			os.Exit(3)
		}

		file.Close()
	}
	err := os.Truncate("myconfig.conf", 0)

	if err != nil {
		log.Register("Can't to truncate mqttloader config file")
		mess,_ := json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
		t := w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		t.Wait()
		os.Exit(3)
	}

	if cmdExp.Arguments != nil{
		createLogFile,id = loadArguments("myconfig.conf", cmdExp.Arguments)
	} else{
		arg_file = ""
		flag = ""
	}

	if experimentId != -1{
		id = experimentId
	}

	if !utils.FileExists("CommandsLog/experiment_"+fmt.Sprint(id)+".json"){
		file,_ :=os.Create("CommandsLog/experiment_"+fmt.Sprint(id)+".json")
		file.Write([]byte(commandLiteral))
		file.Close()
	}

	log.Register("Start experiment "+ strconv.FormatInt(id, 10))

	cmd := exec.Command("./"+w.tool, flag ,arg_file)

	mess,_ := json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: "start",Attr:  cmdExp}) 
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

	cmd.Wait()

	experimentListMutex.Lock()
	if experimentNode.Finished{
		experimentList.Remove(id)
		experimentListMutex.Unlock()
		return
	}
	experimentList.Remove(id)
	experimentListMutex.Unlock()
	
	if err != nil{
		mess,_ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id) , Status: fmt.Sprint(err) + ": " + stderr.String(), Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		mess,_ = json.Marshal(messages.Status{Type: "Client Status", Status: "offline " + err.Error(), Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Crash experiment "+strconv.FormatInt(id, 10)+" error "+err.Error())
		os.Exit(3)
	}

	resultsExperiment := extracExperimentResults(output.String(), createLogFile)

	if resultsExperiment.Publish.AvgThroughput == 0{
		mess,_ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: "Error 10: Hardware Colapse", Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Error experiment "+strconv.FormatInt(id, 10)+" hardware colapse")
	} else {
		mess,_ = json.Marshal(messages.Status{Type: fmt.Sprintf("Experiment Status %d", id), Status: "finish", Attr: messages.Command{}})
		t = w.client.Publish(w.Id+"/Experiments/Status", byte(1), true, string(mess))
		t.Wait()

		log.Register("Finish experiment "+strconv.FormatInt(id, 10))
	}

	resultsExperiment.Meta.ID = uint64(id)

	results,_ := json.Marshal(resultsExperiment)

	t = w.client.Publish(w.Id+"/Experiments/Results", byte(1), false, string(results))
	t.Wait()

	os.Remove("CommandsLog/experiment_"+fmt.Sprint(id)+".json")
}

func (w *Worker) Info(arguments messages.Info){
	var result messages.InfoDisplay

	log.Register("Collecting info")
	
	if arguments.DiscDisplay{
		rootPath := "/"
		if !w.isUnix{
			rootPath = "\\"
		}
		diskStat,_ := disk.Usage(rootPath)
		result.Disk = diskStat.Total/1024/1024
	}
	if arguments.MemoryDisplay{
		vmStat, _ := mem.VirtualMemory()
		result.Ram = vmStat.Total/1024/1024
	}
	if arguments.CpuDisplay{
		cpuStat, _ := cpu.Info()
		result.Cpu = cpuStat[0].ModelName
	}

	resp,_ := json.Marshal(result)

	log.Register("Sending info")

	t := w.client.Publish(w.Id + "/Info", byte(1), false, string(resp))
	t.Wait()
}

func (w *Worker) KeepAlive(){
	for {
		w.Ping()
		time.Sleep(time.Second)
	}
}