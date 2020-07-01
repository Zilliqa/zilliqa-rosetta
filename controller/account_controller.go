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
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/kataras/iris"
)

type AccountController struct {
	*Controller
}

func NewAccountController(app *iris.Application, commonController *Controller) *AccountController {
	c := &AccountController{
		commonController,
	}

	fmt.Println(c.NetworkService)

	app.Post("/account/balance", c.AccountBalance)

	return c
}

func (c *AccountController) AccountBalance(ctx iris.Context) {
	var req types.AccountBalanceRequest

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

	addr := req.AccountIdentifier.Address
	addr, err := c.AccountService.IsValidAddress(addr)
	if err != nil {
		_, _ = ctx.JSON(config.AddressInvalid)
		return
	}

	api := c.NetworkService.NodeAPI(req.NetworkIdentifier.Network)

	rpcClient := provider.NewProvider(api)

	balAndNonce, err1 := rpcClient.GetBalance(addr)

	if err1 != nil {
		if err1.Error() == "-5:Account is not created" {
			_, _ = ctx.JSON(&types.AccountBalanceResponse{
				BlockIdentifier: &types.BlockIdentifier{},
				Balances: []*types.Amount{
					&types.Amount{
						Value: "0",
						Currency: &types.Currency{
							Symbol:   "ZIL",
							Decimals: 12,
							Metadata: map[string]interface{}{},
						},
						Metadata: map[string]interface{}{},
					},
				},
				Metadata: nil,
			})
		} else {
			_, _ = ctx.JSON(&types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: true,
			})
		}

		return
	}

	_, _ = ctx.JSON(&types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{},
		Balances: []*types.Amount{
			&types.Amount{
				Value: balAndNonce.Balance,
				Currency: &types.Currency{
					Symbol:   "ZIL",
					Decimals: 12,
					Metadata: nil,
				},
				Metadata: map[string]interface{}{},
			},
		},
		Metadata: map[string]interface{}{},
	})
}
