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
networks:
 mainnet:
    api: "https://api.zilliqa.com"
    chainid: 1
 testnet:
    api: "https://dev-api.zilliqa.com"
    chainid: 333
```

* rosetta:
  * host: rosetta restful api host
  * port: resetta restful api port
  * version: rosetta sdk version
* networks:
  * mainnet:
    * api: api endpoint of mainnet
    * chainid: chainid of mainnet
  * testnet:
    * api: api endpoint of community testnet
    * chainid: chainid of community testnet


## Restful API

Based on rosetta protocol, zilliqa-rosetta node provides following Restful APIs:

### Network

### Account

### Block

### Construction

### Mempool