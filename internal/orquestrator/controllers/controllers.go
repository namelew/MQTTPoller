package controllers

import (
	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/controllers/experiments"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/controllers/workers"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
)

type Controller struct {
	exp     experiments.Experiments
	workers workers.Workers
}

func Build(o *orquestrator.Orquestrator) Controller {
	return Controller{
		exp:     experiments.Experiments{Orquestrator: o},
		workers: workers.Workers{Orquestrator: o},
	}
}

func (cs Controller) Get(c echo.Context) error {
	switch x := c.Request().URL.Path; {
	case x[:20] == "/orquestrator/worker":
		resp, err := cs.workers.Get(c)
		if err != nil {
			return err
		}
		return c.JSON(200, resp)
	case x[:18] == "/orquestrator/info":
	}
	return echo.ErrBadRequest
}

func (cs Controller) List(c echo.Context) error {
	switch c.Request().URL.Path {
	case "/orquestrator/worker":
		resp, err := cs.workers.List(c)
		if err != nil {
			return err
		}
		return c.JSON(200, resp)
	}
	return echo.ErrBadRequest
}

func (cs Controller) Procedure(c echo.Context) error {
	switch c.Request().URL.Path {
	case "/orquestrator/experiment/start":
		resp, err := cs.exp.StartExperiment(c)

		if err != nil {
			return err
		}

		return c.JSON(200, resp)
	case "/orquestrator/experiment/cancel":
		err := cs.exp.CancelExperiment(c)

		if err != nil {
			return err
		}

		return c.JSON(200, nil)
	}
	return echo.ErrBadRequest
}