package utils

import (
	"encoding/json"
	"os"

	"github.com/namelew/mqtt-bm-latency/messages"
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
	data, _ := os.ReadFile(file)

	var temp messages.Command
	var jsonArg messages.CommandExperiment

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
