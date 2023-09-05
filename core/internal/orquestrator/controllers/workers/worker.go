package workers

import (
	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
)

type Workers struct {
	Orquestrator *orquestrator.Orquestrator
}

func (w Workers) Get(c echo.Context) (messages.Worker, error) {
	worker, err := w.Orquestrator.GetWorker(c.Param("id"))

	if err != nil {
		return messages.Worker{}, echo.NewHTTPError(500, err.Error())
	}

	response := messages.Worker{Id: worker.ID, Online: worker.Online}

	return response, nil
}

func (w Workers) List(c echo.Context) ([]messages.Worker, error) {
	workers := w.Orquestrator.ListWorkers()
	response := make([]messages.Worker, 0)

	for i := range workers {
		response = append(response, messages.Worker{Id: workers[i].ID, Online: workers[i].Online})
	}

	return response, nil
}
