<p align="center">
  <a href="https://www.rosetta-api.org">
    <img width="90%" alt="Rosetta" src="https://www.rosetta-api.org/img/rosetta_header.png">
  </a>
</p>
<h3 align="center">
   Zilliqa Rosetta
</h3>

## Overview
`Zilliqa-rosetta` provides a reference implementation of the [Rosetta specification](https://github.com/coinbase/rosetta-specifications) for Zilliqa in Golang.

## Documentation
Documentation of Zilliqa-rosetta can be found at [Zilliqa developers portal](https://dev.zilliqa.com/docs/exchanges/rosetta-introduction).

## Pre-requisite
To run `Zilliqa-rosetta`, you must install Docker. Please refer to [Docker offical documentation](https://docs.docker.com/get-docker/) on installation instruction.

## System Requirements
`zilliqa-rosetta` has been tested on an [AWS t2.large](https://aws.amazon.com/ec2/instance-types/t2/) instance. This instance type has 2 vCPU and 8 GB of RAM.

## Building `Zilliqa-rosetta`

### Building the zilliqa-rosetta standalone image
Please refer to [rosetta-standalone](rosetta_standalone/README.md) readme

### Building the docker image with the current Zilliqa and Scilla release
```shell script
sh ./build_docker.sh
```

### Building the docker image with a specific scilla and zilliqa tag
```shell script
docker build \
--build-arg ROSETTA_COMMIT_OR_TAG=<ROSETTA_TAG> \
--build-arg SCILLA_COMMIT_OR_TAG=<SCILLA_TAG> \
--build-arg COMMIT_OR_TAG=<ZILLIQA_TAG> \
-t rosetta:1.0 .
```

|Variable|Description|
|---|---|
|ROSETTA_COMMIT_OR_TAG|Override this to download a specific rosetta commit or version tag|
|SCILLA_COMMIT_OR_TAG|Override this to download a specific scilla commit or version tag|
|COMMIT_OR_TAG|Override this to download a specific zilliqa commit of version tag|

## Running `Zilliqa-rosetta`

### Running the zilliqa-rosetta standalone image
Please refer to [rosetta-standalone](rosetta_standalone/README.md) readme

### Configuring `Zilliqa-rosetta`
`Zilliqa-rosetta` configurations yaml allow you to configure which Zilliqa's network and endpoint to connect to. 

```yaml
* rosetta:
  * host: rosetta restful api host
  * port: resetta restful api port
  * version: rosetta sdk version
  * middleware_version: middleware version
* networks:
  * <network_name>:
    * api: api endpoint of mainnet
    * chain_id: chain id of mainnet
    * node_version: zilliqa node verion
  * <network_name>:
    * api: api endpoint of mainnet
    * chain_id: chain id of mainnet
    * node_version: zilliqa node verion
```

Default configuration files for Zilliqa testnet and mainnet has been included in Rosetta root directory.
| Network | Config file |
| ------- | ----------- |
| Testnet | `testnet.config.local.yaml` |
| Mainnet | `mainnet.config.local.yaml` |

If you do not wish to use Zilliqa seed node inside `Zilliqa-rosetta`, you can choose to connect `Zilliqa-rosetta` to any existing Zilliqa public endpoints. 
A sample of this configuration can be found in `config.local.yaml`.

```yaml
rosetta:
 host: "0.0.0.0"
 port: 8080
 version: "1.4.10"
 middleware_version: "0.0.1"

networks:
 mainnet:
    api: "https://api.zilliqa.com"
    chain_id: 1
    node_version: "v8.1.0"
 testnet:
    api: "https://dev-api.zilliqa.com"
    chain_id: 333
    node_version: "v8.1.0"
```

### Zilliqa Seed node keypair generation
You need to generate 2 sets of keys:
* seed node keypair
* whitelisting keypair (if key based whitelisting only)
* Keypair will be in format: `<public key> <private key>`

You can generate the seed node and whitelisting keypair after building the `dockerfile` by running the following command<br />
Note that you can run the docker command without the output redirection if you want to generate another keypair, but do store the output in a secure location
```shell script
mkdir secrets

docker run --rm \
--env GENKEYPAIR="true" \
rosetta:1.0 > secrets/mykey.txt

docker run --rm \
--env GENKEYPAIR="true" \
rosetta:1.0 > secrets/whitelistkey.txt
```

### Zilliqa seed node whitelisting
Zilliqa seed node IP or public key need to be whitelisted by the Zilliqa team in order receive network data. To get whitelisted, please reach out to maintainers[at]zilliqa.com and provide the IP or public key and reason for whitelisting.

Note that IP based whitelisting will be deprecated soon. It's better to start implementing key-based whitelisting.

### Running `Zilliqa-rosetta`
```shell script
docker run -d \
--env BLOCKCHAIN_NETWORK="<NETWORK_TO_USE>" \
--env IP_ADDRESS="<SEED_NODE_HOST_IP>" \
--env IP_WHITELISTING="<Y_or_N>" \
-v $(pwd)/secrets/mykey.txt:/run/zilliqa/mykey.txt \
-v $(pwd)/secrets/whitelistkey.txt:/run/zilliqa/whitelistkey.txt \
-p 4201:4201 -p 4301:4301 -p 4501:4501 -p 33133:33133 -p 8080:8080 \
--name rosetta rosetta:1.0
```

|Variable|Description|
|---|---|
|BLOCKCHAIN_NETWORK|Use either testnet or mainnet no capitals|
|IP_ADDRESS|The seed node's host public ip address|
|IP_WHITELISTING|Y if you chose IP based seed node whitelisting<br />N if you chose Key based seed node whitelisting|

**Note:**
- If your node is going for IP based whitelisting, you will need to set `IP_WHITELISTING` to `Y`.
- If your node is going for public key based whitelisting, you will need to set `IP_WHITELISTING` to `N`.
- For `TESTNET_NAME` and `BUCKET_NAME`, you can refer to [network_meta](network_meta.md). For future version of `Zilliqa-rosetta`, these environment variables will not be neccessary.

### Restarting `Zilliqa-rosetta`

```
docker stop <container name>
docker start <container name>
```

## Restful APIs

### Rosetta restful APIs
Suppored APIs and documentation can be found over at [developer portal](https://dev.zilliqa.com/docs/exchanges/rosetta-introduction).

### Unsupported APIs

These are the following lists of APIs not supported by Zilliqa blockchain:
```
/account/coins
/events/blocks
/search/transactions
```

## How to test
Install the latest rosetta-cli from https://github.com/coinbase/rosetta-cli.

At the time of writing, we are using [rosetta-cli v0.7.1](https://github.com/coinbase/rosetta-cli/releases/tag/v0.7.1).

To begin testing:
1. cd into zilliqa-rosetta folder
2. `go run main.go`
3. Open another terminal and run one of the following depending on the network:

## Testing Data API

**Zilliqa Mainnet**
```
rosetta-cli check:data --configuration-file ./config/mainnet_config.json
```

**Zilliqa Testnet**
```
rosetta-cli check:data --configuration-file ./config/testnet_config.json
```

Note: The `mainnet_config.json` specifically **disables** historical balance lookup and balance tracking options. This is due to the fact that historical balance lookup is not supported on Zilliqa blockchain. In addition, the blockchain rewards miners directly, which results in a single outflow transaction without prior inflow transactions. This will result in `negative balances` errors being raised incorrectly. Hence, both the historical balance lookup and balance tracking options are disabled.

For **testnet** tests, we begin the test from Block 1600000. Some of our much earlier testnet blocks, e.g. Block 270000++, cannot be fetch. Hence, it is recommended to skip certain sections of the testnet blocks.

### End Conditions

The end conditions for `check:data` on testnet is set to 600 seconds. Once the duration is up, the tests will terminate and a report would be displayed:
```
Success: Duration End Condition [Seconds: 600]

+--------------------+--------------------------------+------------+
|  CHECK:DATA TESTS  |          DESCRIPTION           |   STATUS   |
+--------------------+--------------------------------+------------+
| Request/Response   | Rosetta implementation         | PASSED     |
|                    | serviced all requests          |            |
+--------------------+--------------------------------+------------+
| Response Assertion | All responses are correctly    | PASSED     |
|                    | formatted                      |            |
+--------------------+--------------------------------+------------+
| Block Syncing      | Blocks are connected into a    | PASSED     |
|                    | single canonical chain         |            |
+--------------------+--------------------------------+------------+
| Balance Tracking   | Account balances did not go    | NOT TESTED |
|                    | negative                       |            |
+--------------------+--------------------------------+------------+
| Reconciliation     | No balance discrepencies were  | NOT TESTED |
|                    | found between computed and     |            |
|                    | live balances                  |            |
+--------------------+--------------------------------+------------+

+--------------------------+--------------------------------+-----------+
|     CHECK:DATA STATS     |          DESCRIPTION           |   VALUE   |
+--------------------------+--------------------------------+-----------+
| Blocks                   | # of blocks synced             |    101524 |
+--------------------------+--------------------------------+-----------+
| Orphans                  | # of blocks orphaned           |         0 |
+--------------------------+--------------------------------+-----------+
| Transactions             | # of transaction processed     |      5930 |
+--------------------------+--------------------------------+-----------+
| Operations               | # of operations processed      |     10674 |
+--------------------------+--------------------------------+-----------+
| Accounts                 | # of accounts seen             |         0 |
+--------------------------+--------------------------------+-----------+
| Active Reconciliations   | # of reconciliations performed |         0 |
|                          | after seeing an account in a   |           |
|                          | block                          |           |
+--------------------------+--------------------------------+-----------+
| Inactive Reconciliations | # of reconciliations performed |         0 |
|                          | on randomly selected accounts  |           |
+--------------------------+--------------------------------+-----------+
| Exempt Reconciliations   | # of reconciliation failures   |         0 |
|                          | considered exempt              |           |
+--------------------------+--------------------------------+-----------+
| Failed Reconciliations   | # of reconciliation failures   |         0 |
+--------------------------+--------------------------------+-----------+
| Skipped Reconciliations  | # of reconciliations skipped   |         0 |
+--------------------------+--------------------------------+-----------+
| Reconciliation Coverage  | % of accounts that have been   | 0.000000% |
|                          | reconciled                     |           |
+--------------------------+--------------------------------+-----------+
```

### Testing Construction API

```
rosetta-cli check:construction --configuration-file ./config/testnet_config.json
```

#### Testing instructions

First, prefund the address in `prefunded_accounts` section in the `testnet_config.json`.

After executing the above line, rosetta-cli would create an address for testing and expecting a minimum amount:
```
looking for balance {"value":"100000000000000","currency":{"symbol":"ZIL","decimals":12}} on account {"address":"zil1xk5shden2xq4s5dp63v3vq4vyacpux0h3z3jx5","metadata":{"base16":"35A90BB73351815851a1D4591602Ac27701E19f7"}}
```

Fund the stated `zil` address with **at least** the (value + gas fees), e.g. the stated value here is `100000000000000 Qa` = `100 ZIL`, so we would send `120 ZIL` (100 for the minimum amount and 20 for the gas fees). Please adjust the gas fees accordingly if you see a "insufficent balance to broadcast transaction" in the console.

Next, `rosetta-cli` would create a payment transaction from the created address to another created address.


### End Conditions

The end conditions for `check:construction` is set to:
```
  "create_account": 10,
  "transfer": 10
```

Once the tests created 10 new accounts and peform a transfer operation on each account, the tests will terminate and a report would be displayed:
```
Success: {"create_account":10,"transfer":10}

+--------------------------+--------------------------------+-------+
| CHECK:CONSTRUCTION STATS |          DESCRIPTION           | VALUE |
+--------------------------+--------------------------------+-------+
| Addresses Created        | # of addresses created         |    11 |
+--------------------------+--------------------------------+-------+
| Transactions Created     | # of transactions created      |    11 |
+--------------------------+--------------------------------+-------+
| Stale Broadcasts         | # of broadcasts missing after  |     0 |
|                          | stale depth                    |       |
+--------------------------+--------------------------------+-------+
| Transactions Confirmed   | # of transactions seen         |    10 |
|                          | on-chain                       |       |
+--------------------------+--------------------------------+-------+
| Failed Broadcasts        | # of transactions that         |     0 |
|                          | exceeded broadcast limit       |       |
+--------------------------+--------------------------------+-------+
+------------------------------+-------+
| CHECK:CONSTRUCTION WORKFLOWS | COUNT |
+------------------------------+-------+
| request_funds                |     0 |
+------------------------------+-------+
| create_account               |    11 |
+------------------------------+-------+
| transfer                     |    10 |
+------------------------------+-------+
| return_funds                 |     0 |
+------------------------------+-------+
```


## License

This project is available open source under the terms of the [GNU Lesser General Public License 3.0](LICENSE.md).

© 2020 Zilliqa

