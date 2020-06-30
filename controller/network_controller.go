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
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	service2 "github.com/Zilliqa/zilliqa-rosetta/service"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/kataras/iris"
	"strconv"
)

type NetworkController struct {
	Controller
}

func NewNetworkController(app *iris.Application, networkService *service2.NetWorkService) *NetworkController {
	c := &NetworkController{Controller{
		App:            app,
		NetworkService: networkService,
	}}

	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	app.Post("/network/list", c.NetworkList)
	app.Post("/network/options", c.NetworkOptions)
	app.Post("/network/status", c.NetworkStatus)
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
	_, _ = ctx.JSON(&types.NetworkListResponse{NetworkIdentifiers: c.NetworkService.Networks})
}

// Get Network Options
func (c *NetworkController) NetworkOptions(ctx iris.Context) {
	var req types.NetworkRequest

	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		})

		return
	}

	if !c.NetworkService.ContainsIdentifier(req.NetworkIdentifier) {
		_, _ = ctx.JSON(&config.NetworkIdentifierError)
		return
	}

	version := &types.Version{
		RosettaVersion:    c.NetworkService.Config.Rosetta.Version,
		NodeVersion:       c.NetworkService.NodeVersion(req.NetworkIdentifier.Network),
		MiddlewareVersion: nil,
		Metadata:          nil,
	}

	optstatus := []*types.OperationStatus{
		config.StatusSuccess,
		config.StatusFailed,
	}

	operationTypes := []string{config.OP_TYPE_TRANSFER}

	errors := []*types.Error{
		config.NetworkIdentifierError,
		config.BlockIdentifierNil,
		config.BlockNumberInvalid,
		config.GetBlockFailed,
		config.BlockHashInvalid,
		config.GetTransactionFailed,
		config.SignedTxInvalid,
		config.CommitTxFailed,
		config.TxhashInvalid,
		config.UnknownBlock,
		config.ServerNotSupport,
		config.AddressInvalid,
		config.BalanceError,
		config.ParseIntError,
		config.JsonMarshalError,
		config.InvalidPayload,
		config.CurrencyNotConfig,
		config.ParamsError,
		config.ContractAddressError,
		config.PreExecuteError,
		config.QueryBalanceError,
	}

	allow := &types.Allow{
		OperationStatuses: optstatus,
		OperationTypes:    operationTypes,
		Errors:            errors,
	}

	_, _ = ctx.JSON(&types.NetworkOptionsResponse{Version: version, Allow: allow})
}

func (c *NetworkController) NetworkStatus(ctx iris.Context) {
	var req types.NetworkRequest

	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		})

		return
	}

	if !c.NetworkService.ContainsIdentifier(req.NetworkIdentifier) {
		_, _ = ctx.JSON(&config.NetworkIdentifierError)
		return
	}

	api := c.NetworkService.NodeAPI(req.NetworkIdentifier.Network)

	rpcClient := provider.NewProvider(api)

	txBlock, err1 := rpcClient.GetLatestTxBlock()
	if err1 != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: true,
		})

		return
	}

	blockHeight, err2 := strconv.ParseInt(txBlock.Header.BlockNum, 10, 64)

	if err2 != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   "parse block height error: " + err2.Error(),
			Retriable: true,
		})

		return
	}

	blockIden := &types.BlockIdentifier{
		Index: blockHeight,
		Hash:  txBlock.Body.BlockHash,
	}

	timestamp, err3 := strconv.ParseInt(txBlock.Header.Timestamp, 10, 64)

	if err3 != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   "parse timestamp error: " + err3.Error(),
			Retriable: true,
		})

		return
	}

	genesis,err4 := rpcClient.GetTxBlock("0")
	if err4 != nil {
		_, _ = ctx.JSON(&types.Error{
			Code:      0,
			Message:   "get gensis block error: " + err3.Error(),
			Retriable: true,
		})
	}

	genesisBlockNum,_ := strconv.ParseInt(genesis.Header.BlockNum,10,64)
	genesisIden := &types.BlockIdentifier{
		Index: genesisBlockNum,
		Hash:  genesis.Body.BlockHash,
	}

	// todo peer info
	_, _ = ctx.JSON(&types.NetworkStatusResponse{
		CurrentBlockIdentifier: blockIden,
		CurrentBlockTimestamp:  timestamp,
		GenesisBlockIdentifier: genesisIden,
	})

}
