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

	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// BlockAPIService implements the server.BlockAPIService interface.
type BlockAPIService struct {
	Config *config.Config
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewBlockAPIService(config *config.Config) server.BlockAPIServicer {
	return &BlockAPIService{
		Config: config,
	}
}

// Block implements /block endpoint
// When fetching data by BlockIdentifier, it may be possible to only specify the index or hash.
// If neither property is specified, it is assumed that the client is making a request at the current block.
func (s *BlockAPIService) Block(ctx context.Context, request *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	api := s.Config.NodeAPI(request.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)

	if request.BlockIdentifier.Index == nil && request.BlockIdentifier.Hash == nil {
		return nil, config.BlockIdentifierNil
	}

	if request.BlockIdentifier.Index != nil && *request.BlockIdentifier.Index < 0 {
		return nil, config.BlockNumberInvalid
	}

	if request.BlockIdentifier.Index != nil {
		inputTxBlock := fmt.Sprintf("%d", *request.BlockIdentifier.Index)
		txBlock, err := rpcClient.GetTxBlock(inputTxBlock)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: true,
			}
		}

		// verify the hash if present in the request
		if request.BlockIdentifier.Hash != nil && *request.BlockIdentifier.Hash != txBlock.Body.BlockHash {
			return nil, config.BlockHashInvalid
		}

		transactions := make([]*types.Transaction, 0)

		// block may have no transaction
		if txBlock.Header.NumTxns > 0 {
			transactionsList, err1 := rpcClient.GetTxnBodiesForTxBlock(inputTxBlock)

			if err1 != nil {
				if err1.Error() == "-20:Txn Hash not Present" {
					// one of the shard is null
					// use GetTransactionsForTxBlock instead
					shardList, err2 := rpcClient.GetTransactionsForTxBlock(inputTxBlock)

					if err2 != nil {
						return nil, &types.Error{
							Code:      0,
							Message:   err2.Error(),
							Retriable: false,
						}
					}

					for _, shardTxnList := range shardList {
						for _, shardTxn := range shardTxnList {
							txnDetails, err3 := rpcClient.GetTransaction(shardTxn)
							if err3 != nil {
								if err3.Error() == "-20:Txn Hash not Present" {
									// TODO remove this once the block is fixed
									// there are some issues with some of the blocks on mainnet
									// where we cannot retrieve transactions
									// skip these txn hash for now
									continue
								} else {
									return nil, &types.Error{
										Code:      0,
										Message:   err3.Error(),
										Retriable: false,
									}
								}
							}
							rosettaTxn, _ := util.CreateRosTransaction(txnDetails)
							transactions = append(transactions, rosettaTxn)
						}
					}

				} else {
					return nil, &types.Error{
						Code:      0,
						Message:   err1.Error(),
						Retriable: false,
					}
				}
			}

			// transactionList length would be zero if GetTxnBodiesForTxBlock has Txn Hash not Present error
			for _, txnBody := range transactionsList {
				currTransaction, _ := util.CreateRosTransaction(&txnBody)
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
		} else if blocknum == 1 {
			// link block 1 parent to genesis block
			genesisBlock, err2 := rpcClient.GetTxBlock("0")
			if err2 != nil {
				return nil, &types.Error{
					Code:      0,
					Message:   err2.Error(),
					Retriable: false,
				}
			}
			parentBlockIdentifier.Index = (blocknum - 1)
			parentBlockIdentifier.Hash = genesisBlock.Body.BlockHash
		} else {
			// other blocks except 0, 1
			parentBlockIdentifier.Index = (blocknum - 1)
			parentBlockIdentifier.Hash = txBlock.Header.PrevBlockHash
		}

		timestamp, _ := strconv.ParseInt(txBlock.Header.Timestamp, 10, 64)

		rosBlockResponse := new(types.BlockResponse)
		rosBlockResponse.Block = &types.Block{
			BlockIdentifier:       blockIdentifier,
			ParentBlockIdentifier: parentBlockIdentifier,
			Timestamp:             timestamp / 1e3,
			Transactions:          transactions,
		}

		return rosBlockResponse, nil
	}

	return nil, config.BlockNumberInvalid
}

// BlockTransaction implements /block/transaction endpoint
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

	}
	return nil, config.BlockHashInvalid
}
