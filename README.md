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

### Construction

### Mempool

**/mempool**

*Get All Mempool Transactions*

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "mainnet",
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
    "transaction_identifiers": [
        {
            "hash": "5a8af8bd09e6b40ddcb5a0244257590724797fd19925236e0ab29baf525ea1f8"
        },
        {
            "hash": "a6713dd0bc10cea7071cd9ac2f547b067d38df2a53060fbf146a2de4d23fbc76"
        },
        {
            "hash": "37efa722fa36f82b9c1fc8ea7eee44b90e2909d950cbf829905eb05a5238a6bf"
        },
        {
            "hash": "54705f8e2c9319010acec738dd810894e0e1e6e97b23455731aeb87038451303"
        },
        {
            "hash": "e560516d65bf7c54a6fb7e723ca47f70832975d900e1f9b33de2446d9380898f"
        },
        {
            "hash": "69b4d4da4cefc56be22cc3b48dad63aff3e03f6b2dbcbd63db13ad98ef7ffc2c"
        },
        {
            "hash": "5c0d222a54f2bc3e23f0c1e27827f32408132c1ad7d4a7d8cfaba3d2c0414753"
        },
        {
            "hash": "a89c7a0f5ff5b8e176086e5e5475464945b6f5016c93c1662a44f765401c7e13"
        },
        {
            "hash": "a3977b1871125d285b5bbed1fc2d100e6752a1b800120d1b0b0f4d360547e575"
        },
        {
            "hash": "c6b8bf3ef36e5796c0c72cf20c3de0a85f69a3de300da05afb9a42a38cd5042b"
        },
        {
            "hash": "14cff40c72c329fdee306e297b8a5844a81e9283e64d1dfc80d37e3064b278a9"
        },
        {
            "hash": "74fc61c4e4ed3eed64592aac55356afc3c1aad1ebba2fd4709662c43d3223b32"
        },
        {
            "hash": "0c172d26a2e091f2e3d8a5aca9b513ec4424ab4907a2bfbdb87ab399f7b96c51"
        },
        {
            "hash": "4e4e3c70b192cc3b55e386410140040473550a0ed435f12a9d27c07d1381663a"
        },
        {
            "hash": "4861a490361fc1e71942b70986309174dbd28b5de3e8aeba586ac51ee0caa8cf"
        },
        {
            "hash": "6a38457c0b9202f892a9ef7f9f361ee73d2282ff0d8cf48d85cfe97a605e5799"
        }
    ]
}
```

