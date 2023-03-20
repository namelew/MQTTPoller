package orquestration

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
)

func messageHandlerInfos(m mqtt.Message, id int) {
	var response messages.Info
	var output output.Info

	json.Unmarshal(m.Payload(), &response)
	ws[id].ReceiveConfirmation = true

	output.Id = id
	output.NetId = ws[id].Id
	output.Infos = response

	infos = append(infos, output)
}
