package services

import (
	"context"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/Zilliqa/zilliqa-rosetta/mempool"
	"github.com/coinbase/rosetta-sdk-go/types"
	"strings"
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

func (m *MemoryPoolAPIService) Mempool(ctx context.Context, req *types.MempoolRequest) (*types.MempoolResponse,
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
		localPendingMap := pool.PendingPool
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

func (m *MemoryPoolAPIService) MempoolTransaction(ctx context.Context, req *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	// todo parse transaction
	hash := req.TransactionIdentifier.Hash
	pool := m.Pools.GetByType(req.NetworkIdentifier.Network)

	if pool != nil {
		localTxn := pool.PendingPool[hash]
		// haven't send it yet
		if localTxn != nil {
			rosettaTx := &types.Transaction{
				TransactionIdentifier: &types.TransactionIdentifier{Hash: req.TransactionIdentifier.Hash},
				Operations:            []*types.Operation{},
			}
			return &types.MempoolTransactionResponse{
				Transaction: rosettaTx,
			}, nil
		}
	}

	api := m.Config.NodeAPI(req.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)

	pendingResult, err := rpcClient.GetPendingTxn(hash)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	if !strings.Contains(pendingResult.Info, "Pending") {
		return nil, config.TxNotExistInMem
	}

	rosettaTx := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: req.TransactionIdentifier.Hash},
		Operations:            []*types.Operation{},
	}
	return &types.MempoolTransactionResponse{
		Transaction: rosettaTx,
	}, nil
}