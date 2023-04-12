package router

import (
	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/controllers"
)

func Route(o *orquestrator.Orquestrator, port string) {
	api := echo.New()

	controller := controllers.Build(o)

	api.GET("/orquestrator/worker", controller.List)
	api.GET("/orquestrator/worker/:id", controller.Get)
	api.GET("/orquestrator/info", controller.List)
	api.POST("/orquestrator/experiment/start", controller.Procedure)
	api.DELETE("/orquestrator/experiment/cancel/:id/:expid", controller.Procedure)

	api.Logger.Fatal(api.Start(":" + port))
}