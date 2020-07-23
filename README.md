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

**/block**

*Get a Block*

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
    "block_identifier": {
    	"index": 672276,
    	"hash": "23e69657bdf3de2026f4fc9b6b6b38964bf7a7d78b3e004a412ea088116ab5cd"
    }
}
```

Response:

Sample

```json
{
    "block": {
        "block_identifier": {
            "index": 672276,
            "hash": "23e69657bdf3de2026f4fc9b6b6b38964bf7a7d78b3e004a412ea088116ab5cd"
        },
        "parent_block_identifier": {
            "index": 672275,
            "hash": "f6928b5e4017487eb4e519890908f43489999cdf950aeca7ade07ad1f113ef51"
        },
        "timestamp": 1594882462967859,
        "transactions": [
            {
                "transaction_identifier": {
                    "hash": "e26d4cb1fa01003298b626dcc78351f10bc4e19b0c8c77d12f42cbd5d9dae694"
                },
                "operations": [
                    {
                        "operation_identifier": {
                            "index": 0
                        },
                        "type": "transfer",
                        "status": "SUCCESS",
                        "account": {
                            "address": "ab45b5787a3c82d936721cdc53e4503f27f0d363"
                        },
                        "amount": {
                            "value": "-199999000000000000",
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
                            "address": "6aef6306da6de2ffda49912132985ea42ada196d"
                        },
                        "amount": {
                            "value": "199999000000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {...}
                    }
                ]
            },
            {
                "transaction_identifier": {
                    "hash": "71a0da72e03e4d6581505094e1716e02dc23d859923321b9452fd15ea6780403"
                },
                "operations": [
                    {
                        "operation_identifier": {
                            "index": 0
                        },
                        "type": "transfer",
                        "status": "SUCCESS",
                        "account": {
                            "address": "14456246cc53d31f0a2eb00fa0248f7aebd3fa89"
                        },
                        "amount": {
                            "value": "-103594390000000000",
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
                            "address": "824c1086a9620b053f64cf5c216834b03bcad4f9"
                        },
                        "amount": {
                            "value": "103594390000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {...}
                    }
                ]
            },
            {
                "transaction_identifier": {
                    "hash": "8e79839cf02fd01c1669d639756c9dcac303ae193cc44f936b8664231994ec31"
                },
                "operations": [
                    {
                        "operation_identifier": {
                            "index": 0
                        },
                        "type": "transfer",
                        "status": "SUCCESS",
                        "account": {
                            "address": "14456246cc53d31f0a2eb00fa0248f7aebd3fa89"
                        },
                        "amount": {
                            "value": "-1100888000000000",
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
                            "address": "51a7cd099fe7c7a06320f9803f04829c928811a7"
                        },
                        "amount": {
                            "value": "1100888000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {}
                    }
                ]
            }
        ]
    }
}
```

**/block/transaction**

*Get a Block Transaction - Payment*

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
    	"index": 1582509,
    	"hash": "4cc2adbb6fe5f14952b1a7043b0a3fb0a33016fe0de99d1bc2102f349e3cd3ad"
    },
    "transaction_identifier": {
    	"hash": "e03a4dcfce78a7f40a686969260bef57e0e18cead8fa1b60df05edfd69c80415"
    }
}
```

Response:

Sample
__Note__: The operation type is `transfer`.

```json
{
    "transaction": {
        "transaction_identifier": {
            "hash": "e03a4dcfce78a7f40a686969260bef57e0e18cead8fa1b60df05edfd69c80415"
        },
        "operations": [
            {
                "operation_identifier": {
                    "index": 0
                },
                "type": "transfer",
                "status": "SUCCESS",
                "account": {
                    "address": "f0b55a21ed19e6e442833afa9266c166ad639d53"
                },
                "amount": {
                    "value": "-300000000000000",
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
                    "address": "208e1e2c4130e43f8f1329b96767492650597c92"
                },
                "amount": {
                    "value": "300000000000000",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                },
                "metadata": {
                    "gasLimit": "1",
                    "gasPrice": "1000000000",
                    "nonce": "138",
                    "receipt": {
                        "accept": false,
                        "errors": null,
                        "exceptions": null,
                        "success": true,
                        "cumulative_gas": "1",
                        "epoch_num": "1582509",
                        "event_logs": null,
                        "transitions": null
                    },
                    "senderPubKey": "0x027558EDE7BA1EA7A7633F1ACA898CE3DE0F7589C6B5D8C30D91EDE457F6E552F6",
                    "signature": "0x9795884BD13CF195334B5D8E79A425C1582CF8D481091F64840D11BBC19EC893444021C1793AF242703823CF4FCC29E72060D79EEAB0478D9334A73F173A97A4",
                    "version": "21823489"
                }
            }
        ]
    }
}
```

*Get a Block Transaction - Contract Deployment*

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
    "block_identifier": {
    	"index": 670379,
    	"hash": "e71a6d73ec69accb63cb77e67ce6bdde92e6de5a9f1981d8cd9f2f4630031a7b"
    },
    "transaction_identifier": {
    	"hash": "5a3662d689468b423f050824c93343b790a7295d44a4e0f5ebee119ecc18d065"
    }
}
```

Response:

Sample
__Note__: The operation type is `contract_deployment`.

```json
{
    "transaction": {
        "transaction_identifier": {
            "hash": "5a3662d689468b423f050824c93343b790a7295d44a4e0f5ebee119ecc18d065"
        },
        "operations": [
            {
                "operation_identifier": {
                    "index": 0
                },
                "type": "contract_deployment",
                "status": "SUCCESS",
                "account": {
                    "address": "ec69f332f13923c39b3eb08c9b21b7de17289fc2"
                },
                "amount": {
                    "value": "0",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                },
                "metadata": {
                    "code": "\nscilla_version 0\nimport BoolUtils\nlibrary ResolverLib\ntype RecordKeyValue =\n  | RecordKeyValue of String String\nlet nilMessage = Nil {Message}\nlet oneMsg =\n  fun(msg: Message) =>\n    Cons {Message} msg nilMessage\nlet eOwnerSet =\n  fun(address: ByStr20) =>\n    {_eventname: \"OwnerSet\"; address: address}\n(* @deprecated eRecordsSet is emitted instead (since 0.1.1) *)\nlet eRecordSet =\n  fun(key: String) =>\n  fun(value: String) =>\n    {_eventname: \"RecordSet\"; key: key; value: value}\n(* @deprecated eRecordsSet is emitted instead (since 0.1.1) *)\nlet eRecordUnset =\n  fun(key: String) =>\n    {_eventname: \"RecordUnset\"; key: key}\nlet eRecordsSet =\n  fun(registry: ByStr20) =>\n  fun(node: ByStr32) =>\n    {_eventname: \"RecordsSet\"; registry: registry; node: node}\nlet eError =\n  fun(message: String) =>\n    {_eventname: \"Error\"; message: message}\nlet emptyValue = \"\"\nlet mOnResolverConfigured =\n  fun(registry: ByStr20) =>\n  fun(node: ByStr32) =>\n    let m = {_tag: \"onResolverConfigured\"; _amount: Uint128 0; _recipient: registry; node: node} in\n      oneMsg m\nlet copyRecordsFromList =\n  fun (recordsMap: Map String String) =>\n  fun (recordsList: List RecordKeyValue) =>\n    let foldl = @list_foldl RecordKeyValue Map String String in\n      let iter =\n        fun (recordsMap: Map String String) =>\n        fun (el: RecordKeyValue) =>\n          match el with\n          | RecordKeyValue key val =>\n            let isEmpty = builtin eq val emptyValue in\n              match isEmpty with\n              | True => builtin remove recordsMap key\n              | False => builtin put recordsMap key val\n              end\n          end\n      in\n        foldl iter recordsMap recordsList\ncontract Resolver(\n  initialOwner: ByStr20,\n  registry: ByStr20,\n  node: ByStr32,\n  initialRecords: Map String String\n)\nfield vendor: String = \"UD\"\nfield version: String = \"0.1.1\"\nfield owner: ByStr20 = initialOwner\nfield records: Map String String = initialRecords\n(* Sets owner address *)\n(* @ensures a sender address is an owner of the contract *)\n(* @param address *)\n(* @emits OwnerSet if the operation was successful *)\n(* @emits Error if a sender address has no permission for the operation *)\ntransition setOwner(address: ByStr20)\n  currentOwner <- owner;\n  isOkSender = builtin eq currentOwner _sender;\n  match isOkSender with\n  | True =>\n    owner := address;\n    e = eOwnerSet address;\n    event e\n  | _ =>\n    e = let m = \"Sender not owner\" in eError m;\n    event e\n  end\nend\n(* Sets a key value pair *)\n(* @ensures a sender address is an owner of the contract *)\n(* @param key *)\n(* @param value *)\n(* @emits RecordSet if the operation was successful *)\n(* @emits Error if a sender address has no permission for the operation *)\n(* @sends onResolverConfigured to the registry *)\ntransition set(key: String, value: String)\n  currentOwner <- owner;\n  isOkSender = builtin eq currentOwner _sender;\n  match isOkSender with\n  | True =>\n    records[key] := value;\n    e = eRecordsSet registry node;\n    event e;\n    msgs = mOnResolverConfigured registry node;\n    send msgs\n  | _ =>\n    e = let m = \"Sender not owner\" in eError m;\n    event e\n  end\nend\n(* Remove a key from records map *)\n(* @ensures a sender address is an owner of the contract *)\n(* @param key *)\n(* @emits RecordUnset if the operation was successful *)\n(* @emits Error if a sender address has no permission for the operation *)\n(* @sends onResolverConfigured to the registry *)\ntransition unset(key: String)\n  keyExists <- exists records[key];\n  currentOwner <- owner;\n  isOk =\n    let isOkSender = builtin eq currentOwner _sender in\n      andb isOkSender keyExists;\n  match isOk with\n  | True =>\n    delete records[key];\n    e = eRecordsSet registry node;\n    event e;\n    msgs = mOnResolverConfigured registry node;\n    send msgs\n  | _ =>\n    e = let m = \"Sender not owner or key does not exist\" in\n      eError m;\n    event e\n  end\nend\n(* Set multiple keys to records map *)\n(* Removes records from the map if according passed value is empty *)\n(* @ensures a sender address is an owner of the contract *)\n(* @param newRecords *)\n(* @emits RecordsSet if the operation was successful *)\n(* @emits Error if a sender address has no permission for the operation *)\n(* @sends onResolverConfigured to the registry *)\ntransition setMulti(newRecords: List RecordKeyValue)\n  currentOwner <- owner;\n  isOkSender = builtin eq currentOwner _sender;\n  match isOkSender with\n  | True =>\n    oldRecords <- records;\n    newRecordsMap = copyRecordsFromList oldRecords newRecords;\n    records := newRecordsMap;\n    e = eRecordsSet registry node;\n    event e;\n    msgs = mOnResolverConfigured registry node;\n    send msgs\n  | _ =>\n    e = let m = \"Sender not owner\" in eError m;\n    event e\n  end\nend\n",
                    "data": "[{\"vname\":\"_scilla_version\",\"type\":\"Uint32\",\"value\":\"0\"},{\"vname\":\"initialOwner\",\"type\":\"ByStr20\",\"value\":\"0x4887fb6920a8ae50886543ee8aa504da6c9f83bf\"},{\"vname\":\"registry\",\"type\":\"ByStr20\",\"value\":\"0x9611c53be6d1b32058b2747bdececed7e1216793\"},{\"vname\":\"node\",\"type\":\"ByStr32\",\"value\":\"0xd72c3c6e1e3b1b1238b5ba82ff7afe688f542b1cdbfee692a912dd88b1d31f76\"},{\"vname\":\"initialRecords\",\"type\":\"Map String String\",\"value\":[{\"key\":\"ZIL\",\"val\":\"0x803637d03997e4c29729e9ce9e4bc41c0c867354\"}]}]",
                    "gasLimit": "10000",
                    "gasPrice": "1000000000",
                    "nonce": "2007",
                    "receipt": {
                        "accept": false,
                        "errors": null,
                        "exceptions": null,
                        "success": true,
                        "cumulative_gas": "6024",
                        "epoch_num": "670379",
                        "event_logs": null,
                        "transitions": null
                    },
                    "senderPubKey": "0x032E38FCA06A680FFE1BA40956ADA08CB94236FC985B4F7571D455408A0A27E1A2",
                    "signature": "0xB665902059A7519F9A8E118B87ACE4EFDC0FB434475617B19B94E38ABAB68AE8DC85650E781C52CBE281FA527E7556EE9BC593531743A3594BF25C4330EC0165",
                    "version": "65537"
                }
            }
        ]
    }
}
```

#### Displaying Contract Calls Information for Block Transactions
A contract call can be defined in the either one of the following two forms:
1. An account has invoked a function in a contract
2. An account has invoked a function in a contract which further invokes another function in a different contract (a.k.a Chain Calling)

Depending on the functions invoked by the contract, a contract call may perform additional smart contract deposits to some accounts.
These smart contract deposits will be shown under the `operations []` idential to how typical payment transaction is displayed.
Additional metadata information related to the transaction such as the *contract address* and *gas amount* are displayed only at the __final operation block__ to reduce cluttering of metadata.

*Get a Block Transaction - Contract Call without Smart Contract Deposits*

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
    	"index": 1558244,
    	"hash": "4f00a6059b22ebd73e6a60d77fbc20f65bfa3be3f5ae57712422699e3bb031ac"
    },
    "transaction_identifier": {
    	"hash": "ad8a8aa7c1aff0a59a3d56f9c9a72176c344e8a35bbd66e69b2bc7011b44e637"
    }
}
```

Response:

Sample
__Note__: The operation type is `contract_call`.

```json
{
    "transaction": {
        "transaction_identifier": {
            "hash": "ad8a8aa7c1aff0a59a3d56f9c9a72176c344e8a35bbd66e69b2bc7011b44e637"
        },
        "operations": [
            {
                "operation_identifier": {
                    "index": 0
                },
                "type": "contract_call",
                "status": "SUCCESS",
                "account": {
                    "address": "bf6a28839a2f0c3d5d5bf30a84624b8c55897022"
                },
                "amount": {
                    "value": "0",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                },
                "metadata": {
                    "contractAddress": "c36087407e6474e038d7c316a620afe2a752ad0e",
                    "data": "{\"_tag\":\"SubmitHeaderBlock\",\"params\":[{\"vname\":\"new_hash\",\"type\":\"ByStr32\",\"value\":\"0x62c8b569f485f22878f8b31f6e159981be4ea78bdeb09062fbcdcbc4802deae2\"},{\"vname\":\"block\",\"type\":\"Uint64\",\"value\":\"178656\"}]}",
                    "gasLimit": "40000",
                    "gasPrice": "1000000000",
                    "nonce": "352",
                    "receipt": {
                        "accept": false,
                        "errors": null,
                        "exceptions": null,
                        "success": true,
                        "cumulative_gas": "841",
                        "epoch_num": "1558244",
                        "event_logs": [
                            {
                                "_eventname": "SubmitHashSuccess",
                                "address": "0xc36087407e6474e038d7c316a620afe2a752ad0e",
                                "params": [
                                    {
                                        "type": "ByStr32",
                                        "value": "0x62c8b569f485f22878f8b31f6e159981be4ea78bdeb09062fbcdcbc4802deae2",
                                        "vname": "hash"
                                    },
                                    {
                                        "type": "Int32",
                                        "value": "2",
                                        "vname": "code"
                                    }
                                ]
                            }
                        ],
                        "transitions": null
                    },
                    "senderPubKey": "0x025A5A6AFBB5797E44F29FEFA81B43EB3600C70F021B78ABCE7CF2D4D01D467AFF",
                    "signature": "0xA9FA2B79A0927B544528693D51BB7FCAD1E283146310CE3B12167EAA982AF69EB879942A8310D66F4D5E46655C930DFB4664861435F7CCC2E3DACA653A3966FF",
                    "version": "21823489"
                }
            }
        ]
    }
}
```

*Get a Block Transaction - Contract Call with Smart Contract Deposits (With Chain Calls)*

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
    	"index": 1406004,
    	"hash": "84c14dc0685e01b3c7d06f2f2dd9198880b182a82d16ed62a67752560badc6b7"
    },
    "transaction_identifier": {
    	"hash": "17c46c252569a4f3fc41ae45fc6a898892b3f75dde11d517f8b7a037caf658e3"
    }
}
```

Response:


Sample
__Note__: The operation type is `contract_call`, follow by `contract_call_transfer` for subsequent operations with a smart contract deposits.


In the sample, the sequence of operations are as follows:
- Initiator `d707d8a6f049eb7d3e1e768d22578a71b3f39553` -> Contract `8d1109594e019d9c4defe9b6a855f92691a9e4fa` (invokes a contract call to add funds)
- Contract `8d1109594e019d9c4defe9b6a855f92691a9e4fa` -> Recipient `54f57a2017f85c7f2162f833480283f040a80647` (amount is deducted from contract balance and trasnferred to recipient)

```json
{
    "transaction": {
        "transaction_identifier": {
            "hash": "17c46c252569a4f3fc41ae45fc6a898892b3f75dde11d517f8b7a037caf658e3"
        },
        "operations": [
            {
                "operation_identifier": {
                    "index": 0
                },
                "type": "contract_call",
                "status": "SUCCESS",
                "account": {
                    "address": "d707d8a6f049eb7d3e1e768d22578a71b3f39553"
                },
                "amount": {
                    "value": "0",
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
                "type": "contract_call_transfer",
                "status": "SUCCESS",
                "account": {
                    "address": "8d1109594e019d9c4defe9b6a855f92691a9e4fa"
                },
                "amount": {
                    "value": "-123073860347289351",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                }
            },
            {
                "operation_identifier": {
                    "index": 2
                },
                "related_operations": [
                    {
                        "index": 1
                    }
                ],
                "type": "contract_call_transfer",
                "status": "SUCCESS",
                "account": {
                    "address": "54f57a2017f85c7f2162f833480283f040a80647"
                },
                "amount": {
                    "value": "123073860347289351",
                    "currency": {
                        "symbol": "ZIL",
                        "decimals": 12
                    }
                },
                "metadata": {
                    "contractAddress": "8d1109594e019d9c4defe9b6a855f92691a9e4fa",
                    "data": "{\"_tag\": \"AddFunds\", \"params\": []}",
                    "gasLimit": "10000",
                    "gasPrice": "1000000000",
                    "nonce": "8",
                    "receipt": {
                        "accept": false,
                        "errors": null,
                        "exceptions": null,
                        "success": true,
                        "cumulative_gas": "1402",
                        "epoch_num": "1406004",
                        "event_logs": [
                            {
                                "_eventname": "Verifier add funds",
                                "address": "0x54f57a2017f85c7f2162f833480283f040a80647",
                                "params": [
                                    {
                                        "type": "ByStr20",
                                        "value": "0xd707d8a6f049eb7d3e1e768d22578a71b3f39553",
                                        "vname": "verifier"
                                    }
                                ]
                            }
                        ],
                        "transitions": [
                            {
                                "accept": false,
                                "addr": "0x8d1109594e019d9c4defe9b6a855f92691a9e4fa",
                                "depth": 0,
                                "msg": {
                                    "_amount": "123073860347289351",
                                    "_recipient": "0x54f57a2017f85c7f2162f833480283f040a80647",
                                    "_tag": "AddFunds",
                                    "params": [
                                        {
                                            "vname": "initiator",
                                            "type": "ByStr20",
                                            "value": "0xd707d8a6f049eb7d3e1e768d22578a71b3f39553"
                                        }
                                    ]
                                }
                            }
                        ]
                    },
                    "senderPubKey": "0x02BCD59F13A3DF40DE7D6B901B10DA416D2EFDD41E9A3631D6673809D7F5B9C4EF",
                    "signature": "0xB0B321303D8CABDC3E1AD6B3ECD5CECD90A7D0A839C69C1C23A68CC0AFD283DCC9696EB503EBAE167C3DF3943A54E6EAA1D35D28D9F414FBA44109DAAAEF4F56",
                    "version": "21823489"
                }
            }
        ]
    }
}
```

### Construction

**/construction/combine**

*Create Network Transaction from Signature*

__Note__: Before calling `/combine`, please call `/payloads` to have the `unsigned_transaction`. Next, use [goZilliqa SDK](https://github.com/Zilliqa/gozilliqa-sdk) or other Zilliqa's SDKs to craft a transaction object and sign the transaction object; print out the __*signature*__ and __*transaction object*__ in __hexadecimals__. 

Use them as request parameters as follows:


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
    "unsigned_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":1000000000,\"nonce\":184,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"toAddr\":\"4978075dd607933122f4355B220915EFa51E84c7\",\"version\":21823489}", // from /payloads
    "signatures": [
        {
            "signing_payload": {
                "address": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615", // sender account address
                "hex_bytes": "088180b40a10b8011a144978075dd607933122f4355b220915efa51e84c722230a2102e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e2a120a100000000000000000000001d1a94a200032120a100000000000000000000000003b9aca003801", // signed transaction object output in hex, from goZilliqa SDK, etc
                "signature_type": "ecdsa"
            },
            "public_key": {
                "hex_bytes": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e", // sender public key
                "curve_type": "secp256k1"
            },
            "signature_type": "ecdsa",
            "hex_bytes": "af2bb8c883d633a978cfa9b1263de0ad6e55d0f82f75317542db695c62aa50e857e05316b1f0162f4be0acc37b0dc14f2d7f2c1f0b207683be3b2bebdd89deca" // signature of the signed transaction output from goZilliqa SDK, etc
        }
    ]
}
```

Response:

Sample

```json
{
    "signed_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":1000000000,\"nonce\":184,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"signature\":\"af2bb8c883d633a978cfa9b1263de0ad6e55d0f82f75317542db695c62aa50e857e05316b1f0162f4be0acc37b0dc14f2d7f2c1f0b207683be3b2bebdd89deca\",\"toAddr\":\"4978075dd607933122f4355B220915EFa51E84c7\",\"version\":21823489}"
}
```

**/construction/metadata**

*Create a Request to Fetch Metadata*

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
    "options": {
        "method": "transfer"
    }
}
```

Response:

Sample

```json
{
    "metadata": {
        "amount": "Transaction amount to be sent to the recipent's address. This is measured in the smallest price unit Qa (or 10^-12 Zil) in Zilliqa",
        "code": "The smart contract source code. This is present only when deploying a new contract",
        "data": "String-ified JSON object specifying the transition parameters to be passed to a specified smart contract",
        "gasLimit": "The amount of gas units that is needed to be process this transaction",
        "gasPrice": "An amount that a sender is willing to pay per unit of gas for processing this transactionThis is measured in the smallest price unit Qa (or 10^-12 Zil) in Zilliqa",
        "nonce": "A transaction counter in each account. This prevents replay attacks where a transaction sending eg. 20 coins from A to B can be replayed by B over and over to continually drain A's balance",
        "priority": "A flag for this transaction to be processed by the DS committee",
        "pubKey": "Sender's public key of 33 bytes",
        "signature": "An EC-Schnorr signature of 64 bytes of the entire Transaction object as stipulated above",
        "toAddr": "Recipient's account address. This is represented as a String",
        "version": "The decimal conversion of the bitwise concatenation of CHAIN_ID and MSG_VERSION parameters"
    }
}
```

**/construction/derive**

*Derive an Address from a PublicKey*

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
    "public_key": {
        "hex_bytes": "026c7f3b8ac6f615c00c34186cbe4253a2c5acdc524b1cfae544c629d8e3564cfc",
        "curve_type": "secp256k1"
    },
    "metadata": {
    	"type": "bech32"
    }
}
```

Response:

Sample

```json
{
    "address": "zil1y9qmlzmdygfaf4eqfcka4wfx20wzghzl05xazc",
    "metadata": {
        "type": "bech32"
    }
}
```

**/construction/hash**

*Get the Hash of a Signed Transaction*

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
    "signed_transaction": "{\"version\":21823489,\"nonce\":166473,\"toAddr\":\"4BAF5faDA8e5Db92C3d3242618c5B47133AE003C\",\"amount\":1000000,\"pubKey\":\"0246e7178dc8253201101e18fd6f6eb9972451d121fc57aa2a06dd5c111e58dc6a\",\"gasPrice\":1000000000,\"gasLimit\":1,\"code\":\"\",\"data\":\"\",\"signature\":\"0e28d454535a41b2bdc36ad3eade2238e27031bdca248e87417ace809e909c1dde72a3f8e910e82cc3be36dd2c02ed90547c8518f5e329fee1f71e957078b58e\"}"
}
```

Response:

Sample

```json
{
    "transaction_hash": "044f4ac093fbd399f5829c1ecaec76e8fc6cf38367dddf8ee02eede891959d6e"
}
```

**/construction/payloads**

*Generate an Unsigned Transaction and Signing Payloads*

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
	"operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "transfer",
            "status": "",
            "account": {
                "address": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
            },
            "amount": {
                "value": "-2000000000000",
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
            "status": "",
            "account": {
                "address": "0x4978075dd607933122f4355B220915EFa51E84c7"
            },
            "amount": {
                "value": "2000000000000",
                "currency": {
                    "symbol": "ZIL",
                    "decimals": 12
                }
            },
            "metadata": {
                "gasLimit": "1",
                "gasPrice": "1000000000",
                "pubKey": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e"
            }
        }
    ],
    "metadata": {}
}
```

Response:

Sample

```json
{
    "unsigned_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":1000000000,\"nonce\":184,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"toAddr\":\"4978075dd607933122f4355B220915EFa51E84c7\",\"version\":21823489}",
    "payloads": [
        {
            "hex_bytes": "7b22616d6f756e74223a323030303030303030303030302c22636f6465223a22222c2264617461223a22222c226761734c696d6974223a312c226761735072696365223a313030303030303030302c226e6f6e6365223a3138342c227075624b6579223a22303265343465663263356332303331333836666161366361666466356636373331386363363631383731623031313261323734353865363566333761333536353565222c22746f41646472223a2234393738303735646436303739333331323266343335354232323039313545466135314538346337222c2276657273696f6e223a32313832333438397d",
            "address": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615",
            "signature_type": "ecdsa"
        }
    ]
}
```

**/construction/preprocess**

*Create a Request to Fetch Metadata*

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
	"operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "transfer",
            "status": "",
            "account": {
                "address": "f0b55a21ed19e6e442833afa9266c166ad639d53"
            },
            "amount": {
                "value": "-2000000000000",
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
            "status": "",
            "account": {
                "address": "208e1e2c4130e43f8f1329b96767492650597c92"
            },
            "amount": {
                "value": "2000000000000",
                "currency": {
                    "symbol": "ZIL",
                    "decimals": 12
                }
            },
            "metadata": {
                "gasLimit": "1",
                "gasPrice": "1000000000",
                "senderPubKey": "0x027558EDE7BA1EA7A7633F1ACA898CE3DE0F7589C6B5D8C30D91EDE457F6E552F6"
            }
        }
    ],
    "metadata": {}
}
```

Response:

Sample

```json
{
    "options": {
        "amount": "2000000000000",
        "gasLimit": "1",
        "gasPrice": "1000000000",
        "pubKey": "027558EDE7BA1EA7A7633F1ACA898CE3DE0F7589C6B5D8C30D91EDE457F6E552F6",
        "toAddr": "208e1e2c4130e43f8f1329b96767492650597c92"
    }
}
```

**/construction/submit**

*Submit a Signed Transaction*

__Note__: Before calling `/submit`, please call `/combine` to obtain the `signed_transaction` required for the  request parameters.

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
    "signed_transaction": "{\"version\":21823489,\"nonce\":30,\"toAddr\":\"4BAF5faDA8e5Db92C3d3242618c5B47133AE003C\",\"amount\":10000000,\"pubKey\":\"026c7f3b8ac6f615c00c34186cbe4253a2c5acdc524b1cfae544c629d8e3564cfc\",\"gasPrice\":1000000000,\"gasLimit\":1,\"code\":\"\",\"data\":\"\",\"signature\":\"bf1a2a5d8b7eadaebba4a6ed7f5d01c5270b0a288ba8f12fc92fa92d3da9b65b03378cfcdd7c269847d2e3ea850ee16a2533b52bb055d8501578c8f683809e49\"}"
}
```

Response:

Sample

```json
{
    "transaction_identifier": {
        "hash": "5d559dfbbfa98029b961d3f35422f8939095b17c9243fdbd6bf60aaadc41ebf4"
    }
}
```


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

**/mempool/transaction**

*Get a Mempool Transaction*

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
    "transaction_identifier": {
        "hash": "2d326d17cde4a0a4de1b9d342066c9952d015929e49c8c3ace2bebeb2817e621"
    }
}
```

Response:

Sample

```json
{
    "code": 511,
    "message": "tx not exist in mem pool",
    "retriable": false
}
```