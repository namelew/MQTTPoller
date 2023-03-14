package controllers

import (
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/orquestration"
	"github.com/namelew/mqtt-bm-latency/output"
)

func GetWorker(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	switch json_map["wid"].(type) {
	case float64:
		tempid, ok := json_map["wid"].(float64)
		if !ok {
			return echo.ErrInternalServerError
		}

		wid := int(tempid)

		workers := orquestration.GetWorkers(&filters.Worker{WorkerID: uint64(wid)})[0]

		response := output.Worker{Id: wid, NetId: workers.Token, Online: workers.Online, History: nil}

		return c.JSON(200, response)
	default:
		workers := orquestration.GetWorkers(nil)
		response := make([]output.Worker, 0)

		for i := range workers {
			response = append(response, output.Worker{Id: int(workers[i].ID), NetId: workers[i].Token, Online: workers[i].Online, History: nil})
		}
		return c.JSON(200, response)
	}
}

func GetInfo(c echo.Context) error {
	var request input.Info
	err := json.NewDecoder(c.Request().Body).Decode(&request)

	if err != nil {
		return err
	}

	reponse, err := orquestration.GetInfo(request)

	if err != nil {
		return err
	}

	return c.JSON(200, reponse)
}
