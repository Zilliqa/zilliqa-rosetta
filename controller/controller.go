package controller

import (
	service2 "github.com/Zilliqa/zilliqa-rosetta/service"
	"github.com/kataras/iris"
)


type Controller struct {
	app     *iris.Application
	service *service2.Service
}

func NewController(app *iris.Application, service *service2.Service) *Controller {
	c := &Controller{
		app:     app,
		service: service,
	}

	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{
			"message": "pong",
		})
	})
	return c
}