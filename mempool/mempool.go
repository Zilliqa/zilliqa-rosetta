package mempool

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/kataras/golog"
	"time"
)

var log = golog.New()

type PoolTransaction struct {
	SendAddr string
	Nonce    string
	Txn      *core.Transaction
}

type MemPool struct {
	FreshTime   int
	PendingPool map[string]*PoolTransaction
	SentPool    map[string]*PoolTransaction
	ConfirmPool map[string]*PoolTransaction
	Network     config.Network
	Client      *provider.Provider
}

type MemPools []*MemPool

func (mps MemPools) GetByType(t string) *MemPool {
	for _, mp := range mps {
		if mp.Network.Type == t {
			return mp
		}
	}

	return nil
}

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
		Client:      c,
	}
}

func (pool *MemPool) ConfirmCheck() error {
	var hashes []string
	for hash := range pool.SentPool {
		hashes = append(hashes, hash)
	}

	if len(hashes) == 0 {
		return nil
	}

	results, err := pool.Client.GetTransactionBatch(hashes)
	if err != nil {
		return err
	}

	for _, txn := range results {
		delete(pool.SentPool, txn.ID)
		confirmed := pool.SentPool[txn.ID]
		pool.ConfirmPool[txn.ID] = confirmed
	}

	return nil
}

func (pool *MemPool) Start() {
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(pool.FreshTime))
			log.Printf("%s handle mem pool\n", pool.Network.Type)
			// 1. check if transactions within sent pool already get confirmed
			err := pool.ConfirmCheck()
			if err != nil {
				log.Error("check confirm failed ", err.Error())
			}

			// 2. todo check if any transaction inside pending pool need to be sent out
		}
	}()
}
