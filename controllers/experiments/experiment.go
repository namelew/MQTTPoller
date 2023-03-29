package experiments

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/orquestration"
	"github.com/namelew/mqtt-bm-latency/output"
)

type Experiments struct {
	Orquestrator *orquestration.Orquestrator
}

func (e Experiments) StartExperiment(c echo.Context) ([]output.ExperimentResult, error) {
	var request input.Start

	err := json.NewDecoder(c.Request().Body).Decode(&request)

	if err != nil {
		return []output.ExperimentResult{}, echo.ErrBadRequest
	}

	response, err := e.Orquestrator.StartExperiment(request)

	if err != nil {
		return []output.ExperimentResult{}, echo.ErrInternalServerError
	}

	return response, nil
}

func (e Experiments) CancelExperiment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.ErrBadRequest
	}

	expid, err := strconv.Atoi(c.Param("expid"))

	if err != nil {
		return echo.ErrBadRequest
	}

	err = e.Orquestrator.CancelExperiment(id, int64(expid))

	if err != nil {
		return echo.ErrInternalServerError
	}

	return nil
}
