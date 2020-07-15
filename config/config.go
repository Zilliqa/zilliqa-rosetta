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
package config

import (
	"encoding/json"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/viper"
)

var (
	// todo carefully check
	NetworkIdentifierError = &types.Error{
		Code:      400,
		Message:   "network identifier is not supported",
		Retriable: false,
	}

	BlockIdentifierNil = &types.Error{
		Code:      401,
		Message:   "block identifier is empty",
		Retriable: false,
	}

	BlockNumberInvalid = &types.Error{
		Code:      402,
		Message:   "block index is invalid",
		Retriable: false,
	}

	GetBlockFailed = &types.Error{
		Code:      403,
		Message:   "get block failed",
		Retriable: true,
	}

	BlockHashInvalid = &types.Error{
		Code:      404,
		Message:   "block hash is invalid",
		Retriable: false,
	}

	GetTransactionFailed = &types.Error{
		Code:      405,
		Message:   "get transaction failed",
		Retriable: true,
	}

	SignedTxInvalid = &types.Error{
		Code:      406,
		Message:   "signed transaction failed",
		Retriable: false,
	}

	CommitTxFailed = &types.Error{
		Code:      407,
		Message:   "commit transaction failed",
		Retriable: false,
	}
	TxhashInvalid = &types.Error{
		Code:      408,
		Message:   "transaction hash is invalid",
		Retriable: false,
	}
	UnknownBlock = &types.Error{
		Code:      409,
		Message:   "block is not exist",
		Retriable: false,
	}
	ServerNotSupport = &types.Error{
		Code:      500,
		Message:   "services not realize",
		Retriable: false,
	}
	AddressInvalid = &types.Error{
		Code:      501,
		Message:   "address is invalid",
		Retriable: true,
	}
	BalanceError = &types.Error{
		Code:      502,
		Message:   "get balance error",
		Retriable: true,
	}
	ParseIntError = &types.Error{
		Code:      503,
		Message:   "parse integer error",
		Retriable: true,
	}
	JsonMarshalError = &types.Error{
		Code:      504,
		Message:   "json marshal failed",
		Retriable: false,
	}
	InvalidPayload = &types.Error{
		Code:      505,
		Message:   "parse tx payload failed",
		Retriable: false,
	}
	CurrencyNotConfig = &types.Error{
		Code:      506,
		Message:   "currency not config",
		Retriable: false,
	}
	ParamsError = &types.Error{
		Code:      507,
		Message:   "params error",
		Retriable: true,
	}
	ContractAddressError = &types.Error{
		Code:      508,
		Message:   "contract address invalid",
		Retriable: true,
	}
	PreExecuteError = &types.Error{
		Code:      509,
		Message:   "pre execute contract failed",
		Retriable: false,
	}
	QueryBalanceError = &types.Error{
		Code:      510,
		Message:   "query balance failed",
		Retriable: true,
	}
	TxNotExistInMem = &types.Error{
		Code:      511,
		Message:   "tx not exist in mem pool",
		Retriable: false,
	}
	HeightHistoricalLessThanCurrent = &types.Error{
		Code:      512,
		Message:   "historical compute balance height less than req height",
		Retriable: false,
	}
	StoreDbError = &types.Error{
		Code:      513,
		Message:   "db store error",
		Retriable: true,
	}
)

var (
	StatusSuccess = &types.OperationStatus{
		Status:     "SUCCESS",
		Successful: true,
	}
	StatusFailed = &types.OperationStatus{
		Status:     "FAILED",
		Successful: false,
	}
)

const OpTypeTransfer = "transfer"
const OpTypeContractDeployment = "contract_deployment"
const OpTypeContractCall = "contract_call"
const OpTypeContractCallDeposit = "contract_call_deposit"

type Rosetta struct {
	Host string
	Port int
	// rosetta sepc version
	Version       string
	MiddleVersion string
}

type Network struct {
	Type    string
	API     string
	ChainId int
	// zilliqa node version
	NodeVersion string
}

type Config struct {
	Rosetta  Rosetta
	Networks []Network
}

func ParseConfig() (*Config, error) {
	rosetta := viper.Get("rosetta").(map[string]interface{})
	host := rosetta["host"].(string)
	port := rosetta["port"].(int)
	version := rosetta["version"].(string)
	middleVersion := rosetta["middleware_version"].(string)

	var networks []Network
	nws := viper.Get("networks").(map[string]interface{})
	for key, value := range nws {
		v := value.(map[string]interface{})
		api := v["api"].(string)
		chainId := v["chain_id"].(int)
		nodeVersion := v["node_version"].(string)

		nw := Network{
			Type:        key,
			API:         api,
			ChainId:     chainId,
			NodeVersion: nodeVersion,
		}

		networks = append(networks, nw)
	}

	r := Rosetta{Host: host, Port: port, Version: version, MiddleVersion: middleVersion}

	return &Config{
		r, networks,
	}, nil
}

func (config *Config) Stringify() ([]byte, error) {
	return json.Marshal(config)
}

func (config *Config) NodeAPI(network string) string {
	for _, nw := range config.Networks {
		if nw.Type == network {
			return nw.API
		}
	}
	return ""
}

func (config *Config) GetNetworkIdentifier() []*types.NetworkIdentifier {
	nws := config.Networks
	var networks []*types.NetworkIdentifier
	for _, nw := range nws {
		n := &types.NetworkIdentifier{
			Blockchain:           "zilliqa",
			Network:              nw.Type,
			SubNetworkIdentifier: &types.SubNetworkIdentifier{Network: "empty"},
		}

		networks = append(networks, n)
	}
	return networks
}
