package mempool

import (
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/kataras/golog"
	"time"
)

var log = golog.New()

type PoolTransaction struct {
	SendAddr string
	Nonce    string
	Txn      *transaction.Transaction
}

type MemPool struct {
	FreshTime   int
	PendingPool map[string]*PoolTransaction
	SentPool    map[string]*PoolTransaction
	Network     config.Network
	Client *provider.Provider
}

type MemPools []*MemPool

func NewMemPools(intNum int, fresh int, cfg *config.Config) MemPools {
	var pools MemPools
	for _, nw := range cfg.Networks {
		pool := NewMemPool(intNum, fresh, nw)
		pools = append(pools, pool)
		pool.Start()
	}

	return pools
}

func NewMemPool(initNum int, fresh int, network config.Network) *MemPool {
	c := provider.NewProvider(network.API)
	return &MemPool{
		PendingPool: make(map[string]*PoolTransaction, initNum),
		SentPool:    make(map[string]*PoolTransaction, initNum),
		FreshTime:   fresh,
		Network:     network,
		Client: c,
	}
}

func (pool *MemPool) Start() {
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(pool.FreshTime))
			log.Printf("%s handle mem pool\n", pool.Network.Type)

			// 1. check if transactions within sent pool already get confirmed


			// 2. check if any transaction inside pending pool need to be sent out
		}
	}()
}
