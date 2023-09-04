package experiments

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
)

type Experiments struct {
	Orquestrator *orquestrator.Orquestrator
}

func (e Experiments) StartExperiment(c echo.Context) ([]messages.ExperimentResult, error) {
	var request messages.Start

	err := json.NewDecoder(c.Request().Body).Decode(&request)

	if err != nil {
		return []messages.ExperimentResult{}, echo.ErrBadRequest
	}

	response, err := e.Orquestrator.StartExperiment(request)

	if err != nil {
		return []messages.ExperimentResult{}, echo.ErrInternalServerError
	}

	return response, nil
}

func (e Experiments) CancelExperiment(c echo.Context) error {
	expid, err := strconv.Atoi(c.Param("expid"))

	if err != nil {
		return echo.ErrBadRequest
	}

	err = e.Orquestrator.CancelExperiment(c.Param("id"), int64(expid))

	if err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}
