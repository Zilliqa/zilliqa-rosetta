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

	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
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
	inputTxBlockHash := request.BlockIdentifier.Hash
	inputTransactionHash := request.TransactionIdentifier.Hash
	txBlock, err := rpcClient.GetTxBlock(inputTxBlock)

	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		}
	}

	if inputTxBlockHash == txBlock.Body.BlockHash {

		transactionDetails, err1 := rpcClient.GetTransaction(inputTransactionHash)

		if err1 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}

		if transactionDetails.Receipt.EpochNum != inputTxBlock {
			return nil, config.TxhashInvalid
		}

		rosTransaction, err2 := createRosTransaction(transactionDetails)

		if err2 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err2.Error(),
				Retriable: false,
			}
		}

		rosBlockTransactionReponse := new(types.BlockTransactionResponse)
		rosBlockTransactionReponse.Transaction = rosTransaction

		return rosBlockTransactionReponse, nil

	} else {
		return nil, config.BlockHashInvalid
	}
}

// convert zilliqa transaction object to rosetta transaction object
// Transaction
//    - TransactionIdentifier
//    - Operation[]
//        - OperationIdentifier
//        - related_operations
//        - type
//        - status
//        - AccountIdentifier
//        - Amount
//        - metadata
func createRosTransaction(ctx *core.Transaction) (*types.Transaction, error) {
	rosTransaction := new(types.Transaction)
	rosTransactionIdentifier := &types.TransactionIdentifier{Hash: ctx.ID}
	rosTransaction.TransactionIdentifier = rosTransactionIdentifier
	rosOperations := make([]*types.Operation, 0)

	idx := 0

	// if transaction is payment
	if ctx.Code == "" && ctx.Data == nil {
		// sender operation
		senderOperation := new(types.Operation)
		senderOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx),
		}
		senderOperation.Type = config.OpTypeTransfer
		senderOperation.Status = config.StatusSuccess.Status
		senderOperation.Account = &types.AccountIdentifier{
			Address: keytools.GetAddressFromPublic(util.DecodeHex(ctx.SenderPubKey)),
		}
		senderOperation.Amount = createRosAmount(ctx.Amount)

		// recipient operation
		recipientOperation := new(types.Operation)
		recipientOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx + 1),
		}
		recipientOperation.RelatedOperations = []*types.OperationIdentifier{
			{
				Index: int64(idx),
				// previous index = senderOperation
			},
		}
		recipientOperation.Type = config.OpTypeTransfer
		recipientOperation.Status = config.StatusSuccess.Status
		recipientOperation.Account = &types.AccountIdentifier{
			Address: ctx.ToAddr,
		}
		recipientOperation.Amount = createRosAmount(ctx.Amount)

		rosOperations = append(rosOperations, senderOperation, recipientOperation)
	}

	// fmt.Println(ctx.Code)
	// fmt.Println(ctx.Data)

	// TODO
	// add metadata for payment transaction
	// if transaction is contract deployment
	// if contract is contract call

	rosTransaction.Operations = rosOperations
	return rosTransaction, nil
}

func createRosAmount(amount string) *types.Amount {
	return &types.Amount{
		Value: amount,
		Currency: &types.Currency{
			Symbol:   "ZIL",
			Decimals: 12,
		},
	}
}
