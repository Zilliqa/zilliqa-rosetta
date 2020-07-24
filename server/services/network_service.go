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
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	Config *config.Config
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(config *config.Config) server.NetworkAPIServicer {
	return &NetworkAPIService{
		Config: config,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: s.Config.GetNetworkIdentifier(),
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {

	api := s.Config.NodeAPI(request.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)
	txBlock, err := rpcClient.GetLatestTxBlock()
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: true,
		}
	}

	blockHeight, err1 := strconv.ParseInt(txBlock.Header.BlockNum, 10, 64)
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   "parse block height error: " + err1.Error(),
			Retriable: true,
		}
	}

	blockIden := &types.BlockIdentifier{
		Index: blockHeight,
		Hash:  txBlock.Body.BlockHash,
	}

	timestamp, err2 := strconv.ParseInt(txBlock.Header.Timestamp, 10, 64)

	if err2 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   "parse timestamp error: " + err2.Error(),
			Retriable: true,
		}
	}

	genesis, err3 := rpcClient.GetTxBlock("0")

	if err3 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   "get genesis block error: " + err3.Error(),
			Retriable: true,
		}
	}

	genesisBlockNum, _ := strconv.ParseInt(genesis.Header.BlockNum, 10, 64)
	genesisIden := &types.BlockIdentifier{
		Index: genesisBlockNum,
		Hash:  genesis.Body.BlockHash,
	}

	// todo peer info
	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: blockIden,
		CurrentBlockTimestamp:  timestamp / 1000,
		GenesisBlockIdentifier: genesisIden,
	}, nil
}

func (s *NetworkAPIService) NodeVersion(network string) string {
	for _, nw := range s.Config.Networks {
		if nw.Type == network {
			return nw.NodeVersion
		}
	}
	return ""
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {

	version := &types.Version{
		RosettaVersion:    s.Config.Rosetta.Version,
		NodeVersion:       s.NodeVersion(request.NetworkIdentifier.Network),
		MiddlewareVersion: nil,
		Metadata:          nil,
	}

	optstatus := []*types.OperationStatus{
		config.StatusSuccess,
		config.StatusFailed,
	}

	operationTypes := []string{config.OpTypeTransfer}

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
	return &types.NetworkOptionsResponse{
		Version: version,
		Allow:   allow,
	}, nil
}
