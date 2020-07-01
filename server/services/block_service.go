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

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type BlockAPIService struct {
	Config *config.Config
}

func NewBlockAPIService(config *config.Config) server.BlockAPIServicer {
	return &BlockAPIService{
		Config: config,
	}
}

// implements /block endpoint
func (s *BlockAPIService) Block(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, *types.Error) {

	api := s.Config.NodeAPI(request.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)
	inputTxBlock := fmt.Sprintf("%d", *request.BlockIdentifier.Index)
	txBlock, err := rpcClient.GetTxBlock(inputTxBlock)

	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		}
	}

	// check the hash matches
	// assume input hash is without '0x'
	if *request.BlockIdentifier.Hash == txBlock.Body.BlockHash {
		// TODO
		// transactions, err1 := rpcClient.GetTransactionsForTxBlock(inputTxBlock)

		// if err1 != nil {
		// 	return nil, &types.Error{
		// 		Code:      0,
		// 		Message:   err1.Error(),
		// 		Retriable: false,
		// 	}
		// }

		return &types.BlockResponse{
			Block: &types.Block{
				BlockIdentifier: &types.BlockIdentifier{
					Index: *request.BlockIdentifier.Index,
					Hash:  fmt.Sprintf("block %d", *request.BlockIdentifier.Index),
				},
				ParentBlockIdentifier: &types.BlockIdentifier{
					Index: *request.BlockIdentifier.Index,
					Hash:  fmt.Sprintf("block %d", *request.BlockIdentifier.Index),
				},
				Timestamp:    time.Now().UnixNano() / 1000000,
				Transactions: []*types.Transaction{},
			},
		}, nil
	} else {
		return nil, config.BlockHashInvalid
	}
}

// implements /block/transaction endpoint
func (s *BlockAPIService) BlockTransaction(ctx context.Context, request *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {

	api := s.Config.NodeAPI(request.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)
	inputTxBlock := fmt.Sprintf("%d", request.BlockIdentifier.Index)
	inputTransactionHash := request.TransactionIdentifier.Hash
	txBlock, err := rpcClient.GetTxBlock(inputTxBlock)

	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		}
	}

	if request.BlockIdentifier.Hash == txBlock.Body.BlockHash {

		transactions, err1 := rpcClient.GetTransactionsForTxBlock(inputTxBlock)

		if err1 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}

		// check if transaction exists
		isTransactionExist := false

		for _, transactionBlocks := range transactions {
			for _, transactionHash := range transactionBlocks {
				if inputTransactionHash == transactionHash {
					isTransactionExist = true
				}
			}
		}

		if isTransactionExist {
			// TODO create transaction struct
			return &types.BlockTransactionResponse{
				Transaction: &types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: "transaction 1",
					},
					Operations: []*types.Operation{
						{
							OperationIdentifier: &types.OperationIdentifier{
								Index: 0,
							},
							Type:   "Reward",
							Status: "Success",
							Account: &types.AccountIdentifier{
								Address: "account 2",
							},
							Amount: &types.Amount{
								Value: "1000",
								Currency: &types.Currency{
									Symbol:   "ZIL",
									Decimals: 12,
								},
							},
						},
					},
				},
			}, nil
		} else {
			return nil, config.TxhashInvalid
		}

	} else {
		return nil, config.BlockHashInvalid
	}
}
