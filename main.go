package main

import (
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/orquestration"
	"github.com/namelew/mqtt-bm-latency/output"
)

func retWorker(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	workers := orquestration.GetWorkers()

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

		temp_hist := make([]interface{}, 1)
		workers[wid].Historic.Print(temp_hist)
		response := output.Worker{Id: wid, NetId: workers[wid].Id, Online: workers[wid].Status, History: temp_hist}

		return c.JSON(200, response)

	case []interface{}:
		workersid, ok := json_map["wid"].([]interface{})
		response := make([]output.Worker, len(workers))

		if !ok {
			return echo.ErrInternalServerError
		}

		for i := 0; i < len(workers); i++ {
			tempid, ok := workersid[i].(float64)

			if !ok {
				return echo.ErrInternalServerError
			}

			wid := int(tempid)
			temp_hist := make([]interface{}, 1)
			workers[wid].Historic.Print(temp_hist)
			wj := output.Worker{Id: wid, NetId: workers[wid].Id, Online: workers[wid].Status, History: temp_hist}

			if response[0].NetId == "" {
				response[0] = wj
			} else {
				response = append(response, wj)
			}
		}

		return c.JSON(200, response)
	default:
		response := make([]output.Worker, len(workers))
		for i := 0; i < len(workers); i++ {
			temp_hist := make([]interface{}, 1)
			workers[i].Historic.Print(temp_hist)
			wj := output.Worker{Id: i, NetId: workers[i].Id, Online: workers[i].Status, History: temp_hist}
			if response[0].NetId == "" {
				response[0] = wj
			} else {
				response = append(response, wj)
			}
		}
		return c.JSON(200, response)
	}
}

func retInfo(c echo.Context) error {
	var request input.Info
	err := json.NewDecoder(c.Request().Body).Decode(&request)

	if err != nil {
		return err
	}

	reponse := orquestration.GetInfo(request)

	return c.JSON(200,reponse)
}

func main() {
	orquestration.Init()

	api := echo.New()
	api.GET("/orquestrator/worker", retWorker)
	api.GET("/orquestrator/info", retInfo)
	api.POST("/orquestrator/experiment/start", nil)
	api.POST("/orquestrator/experiment/cancel", nil)
	api.Logger.Fatal(api.Start(":8080"))

	orquestration.End()
}
