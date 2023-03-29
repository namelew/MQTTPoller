package workers

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
)

type Workers struct {
	Orquestrator *orquestrator.Orquestrator
}

func (w Workers) Get(c echo.Context) (messages.Worker, error) {
	wid, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return messages.Worker{}, echo.ErrBadRequest
	}

	worker := w.Orquestrator.GetWorker(wid)

	response := messages.Worker{Id: int(worker.ID), NetId: worker.Token, Online: worker.Online, History: nil}

	return response, nil
}

func (w Workers) List(c echo.Context) ([]messages.Worker, error) {
	workers := w.Orquestrator.ListWorkers(nil)
	response := make([]messages.Worker, 0)

	for i := range workers {
		response = append(response, messages.Worker{Id: int(workers[i].ID), NetId: workers[i].Token, Online: workers[i].Online, History: nil})
	}

	return response, nil
}
