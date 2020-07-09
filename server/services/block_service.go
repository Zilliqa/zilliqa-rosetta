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
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/core"
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
	inputTxBlockHash := *request.BlockIdentifier.Hash
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
	if inputTxBlockHash == txBlock.Body.BlockHash {
		transactionsList, err1 := rpcClient.GetTransactionsForTxBlock(inputTxBlock)

		if err1 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}

		transactions := make([]*types.Transaction, 0)

		// TODO fetch all the operations for each transaction?
		for _, shards := range transactionsList {
			for _, transactionHash := range shards {
				transactionIdentifier := &types.TransactionIdentifier{
					Hash: transactionHash,
				}
				currTransaction := &types.Transaction{
					TransactionIdentifier: transactionIdentifier,
				}
				transactions = append(transactions, currTransaction)
			}
		}

		blocknum, _ := strconv.ParseInt(txBlock.Header.BlockNum, 10, 64)
		blockIdentifier := &types.BlockIdentifier{
			Index: blocknum,
			Hash:  txBlock.Body.BlockHash,
		}

		parentBlockIdentifier := new(types.BlockIdentifier)
		if blocknum == 0 {
			parentBlockIdentifier = blockIdentifier
		} else {
			parentBlockIdentifier.Index = (blocknum - 1)
			parentBlockIdentifier.Hash = txBlock.Header.PrevBlockHash
		}

		timestamp, _ := strconv.ParseInt(txBlock.Header.Timestamp, 10, 64)

		rosBlockResponse := new(types.BlockResponse)
		rosBlockResponse.Block = &types.Block{
			BlockIdentifier:       blockIdentifier,
			ParentBlockIdentifier: parentBlockIdentifier,
			Timestamp:             timestamp,
			Transactions:          transactions,
		}

		return rosBlockResponse, nil
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
// "Operations" contain all balance-changing information within a transaction
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

	// TODO need to differentiate between different transactions?
	// need to split from and to operations?
	// if transaction is payment - code and data is empty
	// if transaction is contract deployment - toaddr is 00...0
	// if transaction is contract call - transitions is present under receipt

	// recipient operation
	recipientOperation := new(types.Operation)
	recipientOperation.OperationIdentifier = &types.OperationIdentifier{
		Index: int64(idx),
	}
	recipientOperation.Type = config.OpTypeTransfer
	recipientOperation.Status = getTransactionStatus(ctx.Receipt.Success)
	recipientOperation.Account = &types.AccountIdentifier{
		Address: ctx.ToAddr,
	}

	recipientOperation.Amount = createRosAmount(ctx.Amount)
	recipientOperation.Metadata = createMetadata(ctx)

	rosOperations = append(rosOperations, recipientOperation)

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

func createMetadata(ctx *core.Transaction) map[string]interface{} {
	metadata := make(map[string]interface{})

	if ctx.Code != "" {
		metadata["code"] = ctx.Code
	}

	if ctx.Data != nil {
		metadata["data"] = ctx.Data
	}

	metadata["gasLimit"] = ctx.GasLimit
	metadata["gasPrice"] = ctx.GasPrice
	metadata["nonce"] = ctx.Nonce
	metadata["signature"] = ctx.Signature
	metadata["receipt"] = ctx.Receipt
	metadata["senderPubKey"] = ctx.SenderPubKey
	metadata["version"] = ctx.Version
	return metadata
}

func getTransactionStatus(status bool) string {
	if status == true {
		return config.StatusSuccess.Status
	}
	return config.StatusFailed.Status
}
