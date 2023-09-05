package experiments

import (
	"encoding/json"
	"fmt"
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

	experiment := e.Orquestrator.StartExperiment(request)

	if experiment.Error != "" {
		return experiment.Results, echo.NewHTTPError(
			500,
			fmt.Sprintf("Experiment Error: %s on workers %v", experiment.Error, experiment.WorkerIDs),
		)
	}

	return experiment.Results, nil
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
