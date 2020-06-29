package service

import (
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type NetWorkService struct {
	Networks []*types.NetworkIdentifier
}

func NewNetworkService(cfg *config.Config) *NetWorkService {
	nws := cfg.Networks
	var networks []*types.NetworkIdentifier

	for _, nw := range nws {
		n := &types.NetworkIdentifier{
			Blockchain: "zilliqa",
			Network:    nw.Type,
			SubNetworkIdentifier: &types.SubNetworkIdentifier{Metadata: map[string]interface{}{
				"api":     nw.API,
				"chainId": nw.ChainId,
			}},
		}

		networks = append(networks, n)
	}
	return &NetWorkService{
		Networks: networks,
	}
}
