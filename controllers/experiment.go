package controllers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/orquestration"
)

func StartExperiment(c echo.Context) error {
	var request input.Start

	err := json.NewDecoder(c.Request().Body).Decode(&request)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	response := orquestration.StartExperiment(request)

	return c.JSON(200,response)
}

func CancelExperiment(c echo.Context) error {
	id,err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	expid,err := strconv.Atoi(c.Param("expid"))

	if err != nil {
		return err
	}

	orquestration.CancelExperiment(id, int64(expid))

	return c.JSON(200, nil)
}