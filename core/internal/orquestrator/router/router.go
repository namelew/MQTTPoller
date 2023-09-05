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
	api.GET("/orquestrator/experiment", controller.List)
	api.GET("/orquestrator/experiment/:id", controller.Get)
	api.POST("/orquestrator/experiment/start", controller.Procedure)
	api.POST("/orquestrator/experiment/cancel/:expid", controller.Procedure)
	api.DELETE("/orquestrator/experiment/:id", controller.Delete)

	api.Logger.Fatal(api.Start(":" + port))
}
