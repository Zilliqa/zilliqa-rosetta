# zilliqa-rosetta
Zilliqa node which follows Rosetta Blockchain Standard

## Build docker image

```shell script
sh ./build_docker.sh
```

## Running docker image

## How to use

### configuration

The default configuration file is config.local.yaml

```shell script
rosetta:
 host: "0.0.0.0"
 port: 8080
 version: "1.3.1"
 middleware_version: "0.0.1"

networks:
 mainnet:
    api: "https://api.zilliqa.com"
    chain_id: 1
    node_version: "v6.3.0-alpha.0"
 testnet:
    api: "https://dev-api.zilliqa.com"
    chain_id: 333
    node_version: "v6.3.0-alpha.0"
```

* rosetta:
  * host: rosetta restful api host
  * port: resetta restful api port
  * version: rosetta sdk version
  * middleware_version: middleware version
* networks:
  * mainnet:
    * api: api endpoint of mainnet
    * chain_id: chain id of mainnet
    * node_version: zilliqa node verion
  * testnet:
    * api: api endpoint of community testnet
    * chain_id: chain id of community testnet
    * node_version: zilliqa node verion
    
## Restful API

Based on rosetta protocol, zilliqa-rosetta node provides following Restful APIs:

### Network

**/network/list**

*Get List of Available Networks*

Request:

```json
{
    "metadata": {}
}
```

Response:

Sample

```json
{
    "network_identifiers": [
        {
            "blockchain": "zilliqa",
            "network": "testnet",
            "sub_network_identifier": {
                "network": "",
                "metadata": {
                    "api": "https://dev-api.zilliqa.com",
                    "chainId": 333
                }
            }
        },
        {
            "blockchain": "zilliqa",
            "network": "mainnet",
            "sub_network_identifier": {
                "network": "",
                "metadata": {
                    "api": "https://api.zilliqa.com",
                    "chainId": 1
                }
            }
        }
    ]
}
```

**/network/options**

*Get List of Available Networks*

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet",
        "sub_network_identifier": {
            "network": "empty"
        }
    },
    "metadata": {}
}
```

Response:

Sample

```json
{
    "version": {
        "rosetta_version": "1.3.1",
        "node_version": "v6.3.0-alpha.0"
    },
    "allow": {
        "operation_statuses": [
            {
                "status": "SUCCESS",
                "successful": true
            },
            {
                "status": "FAILED",
                "successful": false
            }
        ],
        "operation_types": [
            "transfer"
        ],
        "errors": [
            {
                "code": 400,
                "message": "network identifier is not supported",
                "retriable": false
            },
            {
                "code": 401,
                "message": "block identifier is empty",
                "retriable": false
            },
            {
                "code": 402,
                "message": "block index is invalid",
                "retriable": false
            },
            {
                "code": 403,
                "message": "get block failed",
                "retriable": true
            },
            {
                "code": 404,
                "message": "block hash is invalid",
                "retriable": false
            },
            {
                "code": 405,
                "message": "get transaction failed",
                "retriable": true
            },
            {
                "code": 406,
                "message": "signed transaction failed",
                "retriable": false
            },
            {
                "code": 407,
                "message": "commit transaction failed",
                "retriable": false
            },
            {
                "code": 408,
                "message": "transaction hash is invalid",
                "retriable": false
            },
            {
                "code": 409,
                "message": "block is not exist",
                "retriable": false
            },
            {
                "code": 500,
                "message": "services not realize",
                "retriable": false
            },
            {
                "code": 501,
                "message": "address is invalid",
                "retriable": true
            },
            {
                "code": 502,
                "message": "get balance error",
                "retriable": true
            },
            {
                "code": 503,
                "message": "parse integer error",
                "retriable": true
            },
            {
                "code": 504,
                "message": "json marshal failed",
                "retriable": false
            },
            {
                "code": 505,
                "message": "parse tx payload failed",
                "retriable": false
            },
            {
                "code": 506,
                "message": "currency not config",
                "retriable": false
            },
            {
                "code": 507,
                "message": "params error",
                "retriable": true
            },
            {
                "code": 508,
                "message": "contract address invalid",
                "retriable": true
            },
            {
                "code": 509,
                "message": "pre execute contract failed",
                "retriable": false
            },
            {
                "code": 510,
                "message": "query balance failed",
                "retriable": true
            }
        ]
    }
}
```

**/network/status**

*Get Network Status*

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet",
        "sub_network_identifier": {
            "network": "empty"
        }
    },
    "metadata": {}
}
```

Response:

Sample

```json
{
    "current_block_identifier": {
        "index": 1521313,
        "hash": "0d0a72c50fcfb98962d3a9ac44c3d9eea96e2a808399f3909380890d41b51893"
    },
    "current_block_timestamp": 1593478612587324,
    "genesis_block_identifier": {
        "index": 0,
        "hash": "1947718b431d25dd65c226f79f3e0a9cc96a948899dab3422993def1494a9c95"
    },
    "peers": null
}
```

### Account

**/account/balance**

*Get an Account Balance*

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet",
        "sub_network_identifier": {
            "network": "empty"
        }
    },
    "account_identifier": {
        "address": "2141bf8b6d2213d4d7204e2ddab92653dc245c5f",
        "sub_account": {
        	"address": "empty"
        },
        "metadata": {}
    },
    "block_identifier": {
    	"index": 0
    }
}
```

Response:

Sample

```json
{
    "block_identifier": {
        "index": 0,
        "hash": ""
    },
    "balances": [
        {
            "value": "979976864000000000",
            "currency": {
                "symbol": "ZIL",
                "decimals": 12
            }
        }
    ]
}
```

### Block

**/block/transaction**

*Get a Block Transaction - Payment Transaction*

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet",
        "sub_network_identifier": {
            "network": "empty"
        }
    },
    "block_identifier": {
    	"index": 1554594,
    	"hash": "70df5c44b703974e8d3c8affd94e51fb22894dc5d95febbe1b9bce1833190701"
    },
    "transaction_identifier": {
    	"hash": "8a12da8883664b45454de293471c245e75ea55127a02213ef36706d67b1f4147"
    }
}
```

Response:

Sample

```json
{
    "transaction": {
        "transaction_identifier": {
            "hash": "8a12da8883664b45454de293471c245e75ea55127a02213ef36706d67b1f4147"
        },
        "operations": [
            {
                "operation_identifier": {
                    "index": 0
                },
                "type": "transfer",
                "status": "SUCCESS",
                "account": {
                    "address": "9e7542fe226de29721f62d3301e052cdcf4db190"
                },
                "amount": {
                    "value": "500000000000000",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                }
            },
            {
                "operation_identifier": {
                    "index": 1
                },
                "related_operations": [
                    {
                        "index": 0
                    }
                ],
                "type": "transfer",
                "status": "SUCCESS",
                "account": {
                    "address": "af1fcb9b52e3cd61ad05fdf2e915c70c7ac22fbc"
                },
                "amount": {
                    "value": "500000000000000",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                }
            }
        ]
    }
}
```

### Construction

### Mempool