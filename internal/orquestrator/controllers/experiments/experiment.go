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

func (w Experiments) Get(c echo.Context) (messages.Worker, error) {
	expid, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return messages.Worker{}, echo.ErrBadRequest
	}

	worker := w.Orquestrator.GetWorker(expid)

	response := messages.Worker{Id: int(worker.ID), NetId: worker.Token, Online: worker.Online}

	return response, nil
}

func (w Experiments) List(c echo.Context) ([]messages.Worker, error) {
	workers := w.Orquestrator.ListWorkers(nil)
	response := make([]messages.Worker, 0)

	for i := range workers {
		response = append(response, messages.Worker{Id: int(workers[i].ID), NetId: workers[i].Token, Online: workers[i].Online})
	}

	return response, nil
}

func (w Experiments) Delete(c echo.Context) (messages.Worker, error) {
	expid, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return messages.Worker{}, echo.ErrBadRequest
	}

	worker := w.Orquestrator.GetWorker(expid)

	response := messages.Worker{Id: int(worker.ID), NetId: worker.Token, Online: worker.Online}

	return response, nil
}