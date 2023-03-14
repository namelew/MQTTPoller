package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/orquestration"
	"github.com/namelew/mqtt-bm-latency/output"
)

func GetWorker(c echo.Context) error {
	switch c.Request().URL.Path {
	case "/orquestrator/worker":
		workers := orquestration.ListWorkers(nil)
		response := make([]output.Worker, 0)

		for i := range workers {
			response = append(response, output.Worker{Id: int(workers[i].ID), NetId: workers[i].Token, Online: workers[i].Online, History: nil})
		}

		return c.JSON(200, response)
	default:
		wid, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return echo.ErrBadRequest
		}

		worker := orquestration.GetWorker(wid)

		response := output.Worker{Id: int(worker.ID), NetId: worker.Token, Online: worker.Online, History: nil}

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
