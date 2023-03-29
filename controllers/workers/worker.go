package workers

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/orquestration"
	"github.com/namelew/mqtt-bm-latency/output"
)

type Workers struct {
	Orquestrator *orquestration.Orquestrator
}

func (w Workers) Get(c echo.Context) (output.Worker, error) {
	wid, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return output.Worker{}, echo.ErrBadRequest
	}

	worker := w.Orquestrator.GetWorker(wid)

	response := output.Worker{Id: int(worker.ID), NetId: worker.Token, Online: worker.Online, History: nil}

	return response, nil
}

func (w Workers) List(c echo.Context) ([]output.Worker, error) {
	workers := w.Orquestrator.ListWorkers(nil)
	response := make([]output.Worker, 0)

	for i := range workers {
		response = append(response, output.Worker{Id: int(workers[i].ID), NetId: workers[i].Token, Online: workers[i].Online, History: nil})
	}

	return response, nil
}
