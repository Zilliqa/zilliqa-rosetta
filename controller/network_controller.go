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
	"github.com/Zilliqa/zilliqa-rosetta/config"
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
	app.Post("/network/options", c.NetworkOptions)
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

// Get Network Options
func (c *NetworkController) NetworkOptions(ctx iris.Context)  {
	var req types.NetworkRequest

	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		})

		return
	}

	if !c.networkService.ContainsIdentifier(req.NetworkIdentifier) {
		_, _ = ctx.JSON(&config.NETWORK_IDENTIFIER_ERROR)
		return
	}

	version := &types.Version{
		RosettaVersion:    c.networkService.Config.Rosetta.Version,
		NodeVersion:       c.networkService.NodeVersion(req.NetworkIdentifier.Network),
		MiddlewareVersion: nil,
		Metadata:          nil,
	}

	optstatus := []*types.OperationStatus{
		config.STATUS_SUCCESS,
		config.STATUS_FAILED,
	}

	operationTypes := []string{config.OP_TYPE_TRANSFER}

	errors := []*types.Error{
		config.NETWORK_IDENTIFIER_ERROR,
		config.BLOCK_IDENTIFIER_NIL,
		config.BLOCK_NUMBER_INVALID,
		config.GET_BLOCK_FAILED,
		config.BLOCK_HASH_INVALID,
		config.GET_TRANSACTION_FAILED,
		config.SIGNED_TX_INVALID,
		config.COMMIT_TX_FAILED,
		config.TXHASH_INVALID,
		config.UNKNOWN_BLOCK,
		config.SERVER_NOT_SUPPORT,
		config.ADDRESS_INVALID,
		config.BALANCE_ERROR,
		config.PARSE_INT_ERROR,
		config.JSON_MARSHAL_ERROR,
		config.INVALID_PAYLOAD,
		config.CURRENCY_NOT_CONFIG,
		config.PARAMS_ERROR,
		config.CONTRACT_ADDRESS_ERROR,
		config.PRE_EXECUTE_ERROR,
		config.QUERY_BALANCE_ERROR,
	}

	allow := &types.Allow{
		OperationStatuses: optstatus,
		OperationTypes:    operationTypes,
		Errors:            errors,
	}

	_, _ = ctx.JSON(&types.NetworkOptionsResponse{Version: version, Allow: allow})
}
