/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package controller

import (
	service2 "github.com/Zilliqa/zilliqa-rosetta/service"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/kataras/iris"
)

type NetworkController struct {
	app            *iris.Application
	networkService *service2.NetWorkService
}

func NewNetworkController(app *iris.Application, networkService *service2.NetWorkService) *NetworkController {
	c := &NetworkController{
		app:            app,
		networkService: networkService,
	}

	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	app.Post("/network/list", c.NetworkList)
	return c
}

// Get List of Available Networks
func (c *NetworkController) NetworkList(ctx iris.Context) {
	var req types.MetadataRequest

	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		})

		return
	}
	_, _ = ctx.JSON(&types.NetworkListResponse{NetworkIdentifiers: c.networkService.Networks})
}
