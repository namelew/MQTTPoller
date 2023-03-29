package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/namelew/mqtt-bm-latency/packages/messages"
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetJsonFromFile(file string, expid int64) (string, int, messages.Command, int) {
	if !FileExists(file) {
		return "", 0, messages.Command{}, 0
	}

	var exec_time int
	var attemps int
	data, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err.Error())
	}

	var temp messages.Command
	var jsonArg messages.CommandExperiment

	json.Unmarshal(data, &temp)

	data, err = json.Marshal(temp.Arguments)

	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(data, &jsonArg)

	jsonArg.Expid = expid
	exec_time = jsonArg.ExecTime
	attemps = jsonArg.Attempts

	data, err = json.Marshal(jsonArg)

	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(data, &temp.Arguments)

	data, err = json.Marshal(temp)

	if err != nil {
		log.Fatal(err.Error())
	}

	return string(data), exec_time, temp, attemps
}
