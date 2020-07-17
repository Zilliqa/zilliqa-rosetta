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
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"strconv"
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
		transactionsList, err1 := rpcClient.GetTxnBodiesForTxBlock(inputTxBlock)

		if err1 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}

		transactions := make([]*types.Transaction, 0)

		// TODO fetch all the operations for each transaction?
		for _, txnBody := range transactionsList {
			currTransaction, _ := util.CreateRosTransaction(&txnBody)
			transactions = append(transactions, currTransaction)
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

		rosTransaction, err2 := util.CreateRosTransaction(transactionDetails)

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
