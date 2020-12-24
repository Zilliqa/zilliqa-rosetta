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
 version: "1.4.9"
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
            "network": "mainnet"
        },
        {
            "blockchain": "zilliqa",
            "network": "testnet"
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
        "network": "testnet"
    },
    "metadata": {}
}
```

Response:

Sample

```json
{
    "version": {
        "rosetta_version": "1.4.9",
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
            "transfer",
            "contract_deployment",
            "contract_call",
            "contract_call_transfer"
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
            },
            {
                "code": 511,
                "message": "tx not exist in mem pool",
                "retriable": false
            },
            {
                "code": 512,
                "message": "historical compute balance height less than req height",
                "retriable": false
            },
            {
                "code": 513,
                "message": "db store error",
                "retriable": true
            },
            {
                "code": 514,
                "message": "public hex error",
                "retriable": false
            },
            {
                "code": 515,
                "message": "unsupported address format",
                "retriable": false
            },
            {
                "code": 516,
                "message": "signature provided in transaction is invalid",
                "retriable": false
            }
        ],
        "historical_balance_lookup": false
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
        "network": "testnet"
    },
    "metadata": {}
}
```

Response:

Sample

```json
{
    "current_block_identifier": {
        "index": 1668170,
        "hash": "cfe255a521942588708213129f6cce4522820fb0aaaf1bb3934f2908ca94b738"
    },
    "current_block_timestamp": 1596617124206,
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
        "network": "testnet"
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
        "hash": "1947718b431d25dd65c226f79f3e0a9cc96a948899dab3422993def1494a9c95"
    },
    "balances": [
        {
            "value": "529909051575",
            "currency": {
                "symbol": "ZIL",
                "decimals": 12
            }
        }
    ],
    "metadata": {
        "nonce": 48
    }
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
        "network": "mainnet"
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
        "timestamp": 1594882462967,
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
                            "address": "zil14dzm27r68jpdjdnjrnw98ezs8unlp5mrhwal7x",
                            "metadata": {
                                "base16": "aB45b5787A3c82d936721cDC53E4503f27F0d363"
                            }
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
                            "address": "zil1dthkxpk6dh30lkjfjysn9xz75s4d5xtd6gmv04",
                            "metadata": {
                                "base16": "6aeF6306DA6DE2FFDA49912132985Ea42ADa196d"
                            }
                        },
                        "amount": {
                            "value": "199999000000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {
                            "gasLimit": "1",
                            "gasPrice": "1000000000",
                            "nonce": "2393",
                            "receipt": {
                                "accept": false,
                                "errors": null,
                                "exceptions": null,
                                "success": true,
                                "cumulative_gas": "1",
                                "epoch_num": "672276",
                                "event_logs": null,
                                "transitions": null
                            },
                            "senderPubKey": "0x0252D19817A2956C34A3CC8C063BEF1F4A6E678FCF613711EC2E0D5ADA536FDBBB",
                            "signature": "0x45AD9BC1413379A34D45B87B54741713392904F63CD741DCDD4C1561E47B68B8113C5E877CB64585E900AC9AC7221240D7753F703F6C1D112F804A0BCE89960B",
                            "version": "65537"
                        }
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
                            "address": "zil1z3zky3kv20f37z3wkq86qfy00t4a875fxxw7sw",
                            "metadata": {
                                "base16": "14456246cc53d31f0a2EB00FA0248F7aEbD3fa89"
                            }
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
                            "address": "zil1sfxppp4fvg9s20myeawzz6p5kqau448eh5npar",
                            "metadata": {
                                "base16": "824C1086A9620B053f64cF5C216834B03bcaD4F9"
                            }
                        },
                        "amount": {
                            "value": "103594390000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {
                            "gasLimit": "1",
                            "gasPrice": "1000000000",
                            "nonce": "70611",
                            "receipt": {
                                "accept": false,
                                "errors": null,
                                "exceptions": null,
                                "success": true,
                                "cumulative_gas": "1",
                                "epoch_num": "672276",
                                "event_logs": null,
                                "transitions": null
                            },
                            "senderPubKey": "0x038274E73930301B82B43345F442A07E0A08BF8B5DCDEFC01CFD688BA3077B194C",
                            "signature": "0x04962BB97CAC6C24E076E46CC487B2959E017ACAF639CD40C8F99CF644B2C897CF44DA4795B65612F29E3E48DF4DFBBDC693FD32FBEF165AC365FBB5B00DD802",
                            "version": "65537"
                        }
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
                            "address": "zil1z3zky3kv20f37z3wkq86qfy00t4a875fxxw7sw",
                            "metadata": {
                                "base16": "14456246cc53d31f0a2EB00FA0248F7aEbD3fa89"
                            }
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
                            "address": "zil12xnu6zvlulr6qceqlxqr7pyznjfgsyd8a909t6",
                            "metadata": {
                                "base16": "51a7Cd099Fe7c7a06320F9803f04829c928811a7"
                            }
                        },
                        "amount": {
                            "value": "1100888000000000",
                            "currency": {
                                "symbol": "ZIL",
                                "decimals": 12
                            }
                        },
                        "metadata": {
                            "gasLimit": "1",
                            "gasPrice": "1000000000",
                            "nonce": "70612",
                            "receipt": {
                                "accept": false,
                                "errors": null,
                                "exceptions": null,
                                "success": true,
                                "cumulative_gas": "1",
                                "epoch_num": "672276",
                                "event_logs": null,
                                "transitions": null
                            },
                            "senderPubKey": "0x038274E73930301B82B43345F442A07E0A08BF8B5DCDEFC01CFD688BA3077B194C",
                            "signature": "0x59C7A52C85B5222F181294249AAD006B79330D405DBF77A3E9CC78DD892F96E37F8E55FE86EC1DF4F646094E22116446D462B8B12CEFFAC8647D6C4D6101F590",
                            "version": "65537"
                        }
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
        "network": "testnet"
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
                    "address": "zil17z645g0dr8nwgs5r8tafyekpv6kk882nxaqr70",
                    "metadata": {
                        "base16": "F0b55a21ED19E6E442833Afa9266C166aD639d53"
                    }
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
                    "address": "zil1yz8putzpxrjrlrcn9xukwe6fyeg9jlyjmnw70a",
                    "metadata": {
                        "base16": "208E1e2c4130E43F8F1329B96767492650597C92"
                    }
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
        "network": "mainnet"
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
                    "address": "zil1a35lxvh38y3u8xe7kzxfkgdhmctj387zs92llt",
                    "metadata": {
                        "base16": "ec69F332F13923C39B3eB08c9b21B7De17289FC2"
                    }
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
        "network": "testnet"
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
                    "address": "zil1ha4z3qu69uxr6h2m7v9ggcjt332cjupzp7c2ae",
                    "metadata": {
                        "base16": "Bf6a28839a2f0c3d5d5Bf30A84624B8c55897022"
                    }
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
        "network": "testnet"
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
- Initiator `zil16ura3fhsf84h60s7w6xjy4u2wxel892n7sq5dp` -> Contract `zil135gsjk2wqxwecn00axm2s40ey6g6ne8668046h` (invokes a contract call to add funds)
- Contract `zil135gsjk2wqxwecn00axm2s40ey6g6ne8668046h` (`8d1109594e019d9c4defe9b6a855f92691a9e4fa`) -> Recipient `zil12n6h5gqhlpw87gtzlqe5sq5r7pq2spj8x2g8pe` (amount is deducted from contract balance and trasnferred to recipient)

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
                    "address": "zil16ura3fhsf84h60s7w6xjy4u2wxel892n7sq5dp",
                    "metadata": {
                        "base16": "d707D8a6F049Eb7d3E1e768D22578A71b3f39553"
                    }
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
                    "address": "zil135gsjk2wqxwecn00axm2s40ey6g6ne8668046h",
                    "metadata": {
                        "base16": "8D1109594E019D9C4DEFe9b6A855F92691A9E4fA"
                    }
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
                    "address": "zil12n6h5gqhlpw87gtzlqe5sq5r7pq2spj8x2g8pe",
                    "metadata": {
                        "base16": "54F57A2017F85c7F2162f833480283f040a80647"
                    }
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

## Construction Flow
The construction flow is always in this sequence:
1. /construction/derive
2. /construction/preprocess
3. /construction/metadata
4. /construction/payloads
5. /construction/parse
6. /construction/combine
7. /construction/parse (to confirm correctness)
8. /construction/hash
9. /construction/submit

**/construction/combine**

*Create Network Transaction from Signature*

__Note__: Before calling `/combine`, please call `/payloads` to have the `unsigned_transaction`. Next, use [goZilliqa SDK](https://github.com/Zilliqa/gozilliqa-sdk) or other Zilliqa's SDKs to craft a transaction object and sign the transaction object; print out the __*signature*__ and __*transaction object*__ in __hexadecimals__. 

Refer to the `signRosettaTransaction.js` in the `examples` folder to craft and sign a transaction object.

Use them as request parameters as follows:

```
{
    ...,
    "unsigned_transaction": ... // from /payloads
    "signatures": [
        {
            "signing_payload": {
                "address": "string", // sender account address
                "hex_bytes": "string",  // signed transaction object in hexadecimals representation obtained after signing with goZilliqa SDK or other Zilliqa SDK
                "signature_type": "ecdsa"
            },
            "public_key": {
                "hex_bytes": "string", // sender public key
                "curve_type": "secp256k1"
            },
            "signature_type": "ecdsa",
            "hex_bytes": "string" // signature of the signed transaction object 
        }
    ]
}

```


Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet"
    },
    "unsigned_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":2000000000,\"nonce\":187,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"senderAddr\":\"zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}",
    "signatures": [
        {
            "signing_payload": {
                "account_identifier": {
                    "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
                    "metadata": {
                        "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
                    }
                },
                "hex_bytes": "088180b40a10bb011a144978075dd607933122f4355b220915efa51e84c722230a2102e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e2a120a100000000000000000000001d1a94a200032120a10000000000000000000000000773594003801",
                "signature_type": "schnorr_1"
            },
            "public_key": {
                "hex_bytes": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e",
                "curve_type": "secp256k1"
            },
            "signature_type": "schnorr_1",
            "hex_bytes": "fcb93583d963a7c11f52f04b1ecbd129aa3df896e618b47ff163dc18c53b59afc4289851fd2d5a50eaa7d7ae0763eb912797b0b34e1cf1e6d3865a218e1066b7"
        }
    ]
}
```

Response:

Sample

```json
{
    "signed_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":2000000000,\"nonce\":187,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"senderAddr\":\"zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r\",\"signature\":\"fcb93583d963a7c11f52f04b1ecbd129aa3df896e618b47ff163dc18c53b59afc4289851fd2d5a50eaa7d7ae0763eb912797b0b34e1cf1e6d3865a218e1066b7\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}"
}
```

**/construction/metadata**

*Create a Request to Fetch Metadata*

Request:

`options` is from `/construction/preprocess`
```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet"
    }, 
    "options": {
        "amount": "2000000000000",
        "gasLimit": "1",
        "gasPrice": "2000000000",
        "senderAddr": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
        "toAddr": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0"
    },
    "public_keys": [
        {
            "hex_bytes": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e",
            "curve_type": "secp256k1"
        }
    ]
}
```

Response:

Sample

```json
{
    "metadata": {
        "amount": "2000000000000",
        "gasLimit": "1",
        "gasPrice": "2000000000",
        "nonce": 187,
        "pubKey": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e",
        "senderAddr": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
        "toAddr": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0",
        "version": 21823489
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
        "network": "mainnet"
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
    "account_identifier": {
        "address": "zil1y9qmlzmdygfaf4eqfcka4wfx20wzghzl05xazc",
        "metadata": {
            "base16": "2141BF8B6D2213d4d7204E2DDAB92653dC245c5F"
        }
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
        "network": "testnet"
    },
	"signed_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":1000000000,\"nonce\":186,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"signature\":\"51c69af638ad7afd39841a7abf937d5df99e20adedc4287f43c8070d497ba78136c951192b3920914feb83b9272ccb2ca7facd835dfad10eff2b848b13616daf\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}"
}
```

Response:

Sample

```json
{
    "transaction_identifier": {
        "hash": "a17367c8bcd83cdc2d9ede4571c8e27ad74278ae195263f13e10ba84f12ab13c"
    }
}
```

**/construction/parse**

*Parse a Transaction*

__Note__: Set the `signed` flag accordingly if the transaction is signed or unsigned.

Request:

```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet"
    },
    "signed": true,
    "transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":1000000000,\"nonce\":186,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"signature\":\"51c69af638ad7afd39841a7abf937d5df99e20adedc4287f43c8070d497ba78136c951192b3920914feb83b9272ccb2ca7facd835dfad10eff2b848b13616daf\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}"
}
```

Response:

Sample

```json
{
    "signers": [
        "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r"
    ],
    "operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "transfer",
            "status": "",
            "account": {
                "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
                "metadata": {
                    "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
                }
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
                "address": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0",
                "metadata": {
                    "base16": "4978075dd607933122f4355B220915EFa51E84c7"
                }
            },
            "amount": {
                "value": "2000000000000",
                "currency": {
                    "symbol": "ZIL",
                    "decimals": 12
                }
            }
        }
    ],
    "account_identifier_signers": [
        {
            "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
            "metadata": {
                "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
            }
        }
    ]
}
```

**/construction/payloads**

*Generate an Unsigned Transaction and Signing Payloads*

Request:

`metadata` for `operation_identifier 1` is from `/construction/metadata`
```json
{
    "network_identifier": {
        "blockchain": "zilliqa",
        "network": "testnet"
    },
	"operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "transfer",
            "status": "",
            "account": {
                "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
                "metadata": {
                    "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
                }
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
                "address": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0",
                "metadata": {
                    "base16": "4978075dd607933122f4355B220915EFa51E84c7"
                }
            },
            "amount": {
                "value": "2000000000000",
                "currency": {
                    "symbol": "ZIL",
                    "decimals": 12
                }
            }
        }
    ],
    "metadata": {       // from construction/metadata
        "amount": "2000000000000",
        "gasLimit": "1",
        "gasPrice": "2000000000",
        "nonce": 187,
        "pubKey": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e",
        "senderAddr": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
        "toAddr": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0",
        "version": 21823489
    },
    "public_keys": [
        {
            "hex_bytes": "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e",
            "curve_type": "secp256k1"
        }
    ]
}
```

Response:

Sample

```json
{
    "unsigned_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":2000000000,\"nonce\":187,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"senderAddr\":\"zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}",
    "payloads": [
        {
            "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
            "hex_bytes": "088180b40a10bb011a144978075dd607933122f4355b220915efa51e84c722230a2102e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e2a120a100000000000000000000001d1a94a200032120a10000000000000000000000000773594003801",
            "account_identifier": {
                "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
                "metadata": {
                    "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
                }
            },
            "signature_type": "schnorr_1"
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
        "network": "testnet"
    },
	"operations": [
        {
            "operation_identifier": {
                "index": 0
            },
            "type": "transfer",
            "status": "",
            "account": {
                "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
                "metadata": {
                    "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
                }
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
                "address": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0",
                "metadata": {
                    "base16": "4978075dd607933122f4355B220915EFa51E84c7"
                }
            },
            "amount": {
                "value": "2000000000000",
                "currency": {
                    "symbol": "ZIL",
                    "decimals": 12
                }
            },
            "metadata": {
                "senderPubKey": "0x02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e"
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
        "gasPrice": "2000000000",
        "senderAddr": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
        "toAddr": "zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0"
    },
    "required_public_keys": [
        {
            "address": "zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r",
            "metadata": {
                "base16": "99f9d482abbdC5F05272A3C34a77E5933Bb1c615"
            }
        }
    ]
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
        "network": "testnet"
    },
    "signed_transaction": "{\"amount\":2000000000000,\"code\":\"\",\"data\":\"\",\"gasLimit\":1,\"gasPrice\":2000000000,\"nonce\":187,\"pubKey\":\"02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e\",\"senderAddr\":\"zil1n8uafq4thhzlq5nj50p55al9jvamr3s45hm49r\",\"signature\":\"fcb93583d963a7c11f52f04b1ecbd129aa3df896e618b47ff163dc18c53b59afc4289851fd2d5a50eaa7d7ae0763eb912797b0b34e1cf1e6d3865a218e1066b7\",\"toAddr\":\"zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0\",\"version\":21823489}"
}
```

Response:

Sample

```json
{
    "transaction_identifier": {
        "hash": "963a984ee255cfd881b337a52caf699d4f05799c45cc0948d8a8ce72a6a12d8e"
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
        "network": "testnet"
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
            "hash": "af6e2a81812f7834312e8e2358b51f2f9d7ca696c4d315258102ed868389a7c1"
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
        "network": "testnet"
    },
    "transaction_identifier": {
        "hash": "af6e2a81812f7834312e8e2358b51f2f9d7ca696c4d315258102ed868389a7c1"
    }
}
```

Response:

Sample

```json
{
    "code": 0,
    "message": "transaction not pending",
    "retriable": false
}
```

### Unsupported APIs

These are the following lists of APIs not supported by Zilliqa blockchain:

```
/account/coins
/events/blocks
/search/transactions

```

## How to test
Install the latest rosetta-cli from https://github.com/coinbase/rosetta-cli.

At the time of writing, we are using **rosetta-cli v0.6.6**.

To begin testing:
1. cd into zilliqa-rosetta folder
2. `go run main.go`
3. Open another terminal and run one of the following depending on the network:

### Testing Data API

**Zilliqa Mainnet**
```
rosetta-cli check:data --configuration-file <zilliqa-rosetta>/config/mainnet_config.json
```

**Zilliqa Testnet**
```
rosetta-cli check:data --configuration-file <zilliqa-rosetta>/config/testnet_config.json
```

Note: The `mainnet_config.json` specifically **disables** historical balance lookup and balance tracking options. This is due to the fact that historical balance lookup is not supported on Zilliqa blockchain. In addition, the blockchain rewards miners directly, which results in a single outflow transaction without prior inflow transactions. This will result in `negative balances` errors being raised incorrectly. Hence, both the historical balance lookup and balance tracking options are disabled.

For **testnet** tests, we begin the test from Block 1600000. Some of our much earlier testnet blocks, e.g. Block 270000++, cannot be fetch. Hence, it is recommended to skip certain sections of the testnet blocks.

### Testing Construction API

```
rosetta-cli check:construction --configuration-file ./config/testnet_config.json
```

#### How to execute
First, prefund the address in `prefunded_accounts` section in the `testnet_config.json`.

After executing the above line, rosetta-cli would create an address for testing and expecting a minimum amount:
```
looking for balance {"value":"100000000000000","currency":{"symbol":"ZIL","decimals":12}} on account {"address":"zil1xk5shden2xq4s5dp63v3vq4vyacpux0h3z3jx5","metadata":{"base16":"35A90BB73351815851a1D4591602Ac27701E19f7"}}
```

Fund the stated zil address with **at least** the (value + gas fees), e.g. the stated value here is `100000000000000 Qa` = `100 ZIL`, so we would send `120 ZIL` (100 for the minimum amount and 20 for the gas fees). Please adjust the gas fees accordingly if you see a "insufficent balance to broadcast transaction" in the console.

Next, rosetta-cli would create a payment transaction from the created address to another created address.

Lastly, the test is completed if you see this:
```
broadcast complete for job "transfer (13)" with transaction hash "ed81f9a4fab4759d9836e3ab6eeb550bab08880787f0ab0c5a464842a1662486"
```

The construction API test would continue until the funds of the created accounts are emptied. You may halt the test at any time.