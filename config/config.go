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
	NETWORK_IDENTIFIER_ERROR = &types.Error{
		Code:      400,
		Message:   "network identifier is not supported",
		Retriable: false,
	}

	BLOCK_IDENTIFIER_NIL = &types.Error{
		Code:      401,
		Message:   "block identifier is empty",
		Retriable: false,
	}

	BLOCK_NUMBER_INVALID = &types.Error{
		Code:      402,
		Message:   "block index is invalid",
		Retriable: false,
	}

	GET_BLOCK_FAILED = &types.Error{
		Code:      403,
		Message:   "get block failed",
		Retriable: true,
	}

	BLOCK_HASH_INVALID = &types.Error{
		Code:      404,
		Message:   "block hash is invalid",
		Retriable: false,
	}

	GET_TRANSACTION_FAILED = &types.Error{
		Code:      405,
		Message:   "get transaction failed",
		Retriable: true,
	}

	SIGNED_TX_INVALID = &types.Error{
		Code:      406,
		Message:   "signed transaction failed",
		Retriable: false,
	}

	COMMIT_TX_FAILED = &types.Error{
		Code:      407,
		Message:   "commit transaction failed",
		Retriable: false,
	}
	TXHASH_INVALID = &types.Error{
		Code:      408,
		Message:   "transaction hash is invalid",
		Retriable: false,
	}
	UNKNOWN_BLOCK = &types.Error{
		Code:      409,
		Message:   "block is not exist",
		Retriable: false,
	}
	SERVER_NOT_SUPPORT = &types.Error{
		Code:      500,
		Message:   "service not realize",
		Retriable: false,
	}
	ADDRESS_INVALID = &types.Error{
		Code:      501,
		Message:   "address is invalid",
		Retriable: true,
	}
	BALANCE_ERROR = &types.Error{
		Code:      502,
		Message:   "get balance error",
		Retriable: true,
	}
	PARSE_INT_ERROR = &types.Error{
		Code:      503,
		Message:   "parse integer error",
		Retriable: true,
	}
	JSON_MARSHAL_ERROR = &types.Error{
		Code:      504,
		Message:   "json marshal failed",
		Retriable: false,
	}
	INVALID_PAYLOAD = &types.Error{
		Code:      505,
		Message:   "parse tx payload failed",
		Retriable: false,
	}
	CURRENCY_NOT_CONFIG = &types.Error{
		Code:      506,
		Message:   "currency not config",
		Retriable: false,
	}
	PARAMS_ERROR = &types.Error{
		Code:      507,
		Message:   "params error",
		Retriable: true,
	}
	CONTRACT_ADDRESS_ERROR = &types.Error{
		Code:      508,
		Message:   "contract address invalid",
		Retriable: true,
	}
	PRE_EXECUTE_ERROR = &types.Error{
		Code:      509,
		Message:   "pre execute contract failed",
		Retriable: false,
	}
	QUERY_BALANCE_ERROR = &types.Error{
		Code:      510,
		Message:   "query balance failed",
		Retriable: true,
	}
	TX_NOT_EXIST_IN_MEM = &types.Error{
		Code:      511,
		Message:   "tx not exist in mem pool",
		Retriable: false,
	}
	HEIGHT_HISTORICAL_LESS_THAN_CURRENT = &types.Error{
		Code:      512,
		Message:   "historical compute balance height less than req height",
		Retriable: false,
	}
	STORE_DB_ERROR = &types.Error{
		Code:      513,
		Message:   "db store error",
		Retriable: true,
	}
)

var (
	STATUS_SUCCESS = &types.OperationStatus{
		Status:     "SUCCESS",
		Successful: true,
	}
	STATUS_FAILED = &types.OperationStatus{
		Status:     "FAILED",
		Successful: false,
	}
)

const OP_TYPE_TRANSFER = "transfer"

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
