package services

// @deprecated due to disabling of GetPendingTxns since Zilliqa V8.0.0
// https://github.com/Zilliqa/Zilliqa/releases/tag/v8.0.0
// https://www.rosetta-api.org/docs/all_methods.html

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	util2 "github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/Zilliqa/zilliqa-rosetta/mempool"
	"github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type MemoryPoolAPIService struct {
	Config *config.Config
	Pools  mempool.MemPools
}

func NewMemoryPoolAPIService(config *config.Config) *MemoryPoolAPIService {
	pools := mempool.NewMemPools(10, 10, config)

	return &MemoryPoolAPIService{
		Config: config,
		Pools:  pools,
	}
}

func (m *MemoryPoolAPIService) AddTransaction(ctx context.Context, network *types.NetworkIdentifier, txn *transaction.Transaction) error {
	api := m.Config.NodeAPI(network.Network)
	rpcClient := provider.NewProvider(api)

	t, _ := rpcClient.GetTransaction(txn.ID)
	if t != nil && t.ID != "" {
		return errors.New("transaction already sent and confirmed")
	}

	pool := m.Pools.GetByType(network.Network)
	if pool != nil {
		poolTxn := pool.SentPool[txn.ID]
		if poolTxn != nil {
			return errors.New("transaction already sent")
		}
	}
	pl := txn.ToTransactionPayload()
	payload, _ := json.Marshal(pl)
	fmt.Println(string(payload))
	rsp, err1 := rpcClient.CreateTransaction(pl)
	fmt.Println(rsp)
	if err1 == nil && rsp.Error == nil {
		addr := keytools.GetAddressFromPublic(util2.DecodeHex(txn.SenderPubKey))
		pool.SentPool[txn.ID] = &mempool.PoolTransaction{
			SendAddr: addr,
			Nonce:    txn.Nonce,
			Txn:      txn,
		}
		return nil
	} else {
		if err1 != nil {
			return err1
		} else {
			return errors.New(rsp.Error.Message)
		}
	}
}

// @deprecated
func (m *MemoryPoolAPIService) Mempool(ctx context.Context, req *types.NetworkRequest) (*types.MempoolResponse,
	*types.Error) {

	api := m.Config.NodeAPI(req.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)

	pendings, err := rpcClient.GetPendingTxns()
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	var mergedPendings map[string]interface{}

	mergedPendings = make(map[string]interface{}, len(pendings.Txns))
	for _, pending := range pendings.Txns {
		if strings.Contains(pending.Info, "Pending") {
			mergedPendings[pending.TxnHash] = nil
		}
	}

	pool := m.Pools.GetByType(req.NetworkIdentifier.Network)
	if pool != nil {
		localPendingMap := pool.SentPool
		for k, _ := range localPendingMap {
			mergedPendings[k] = nil
		}
	}

	transactionIdentifiers := make([]*types.TransactionIdentifier, 0)
	for k, _ := range mergedPendings {
		transactionIdentifiers = append(transactionIdentifiers, &types.TransactionIdentifier{
			k,
		})
	}

	resp := &types.MempoolResponse{
		TransactionIdentifiers: transactionIdentifiers,
	}

	return resp, nil
}

// @deprecated
func (m *MemoryPoolAPIService) MempoolTransaction(ctx context.Context, req *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	hash := req.TransactionIdentifier.Hash
	pool := m.Pools.GetByType(req.NetworkIdentifier.Network)

	if pool != nil {
		localTxn := pool.SentPool[hash]
		if localTxn != nil {
			rosettaTx, err0 := util.CreateRosTransaction(util.ToCoreTransaction(localTxn.Txn))
			if err0 != nil {
				return nil, err0
			}
			return &types.MempoolTransactionResponse{
				Transaction: rosettaTx,
			}, nil
		}
	}

	return nil, &types.Error{
		Code:      0,
		Message:   "transaction not pending",
		Retriable: false,
	}

	//api := m.Config.NodeAPI(req.NetworkIdentifier.Network)
	//rpcClient := provider.NewProvider(api)
	//
	//pendingResult, err := rpcClient.GetPendingTxn(hash)
	//if err != nil {
	//	return nil, &types.Error{
	//		Code:      0,
	//		Message:   err.Error(),
	//		Retriable: false,
	//	}
	//}
	//
	//if !strings.Contains(pendingResult.Info, "Pending") {
	//	return nil, config.TxNotExistInMem
	//}

}
