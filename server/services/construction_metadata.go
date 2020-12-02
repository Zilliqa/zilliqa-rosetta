package services

import (
	"context"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionMetadata /construction/metadata
// options structure is from /preprocess
// fetch online metadata
func (c *ConstructionAPIService) ConstructionMetadata(
	ctx context.Context,
	req *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	resp := &types.ConstructionMetadataResponse{
		Metadata: make(map[string]interface{}),
	}

	if req.Options[rosettaUtil.AMOUNT] == nil {
		return nil, config.ParamsError
	}

	if req.Options[rosettaUtil.GAS_LIMIT] == nil {
		return nil, config.ParamsError
	}

	if req.Options[rosettaUtil.GAS_PRICE] == nil {
		return nil, config.ParamsError
	}

	// if req.Options[rosettaUtil.PUB_KEY] == nil {
	// 	return nil, config.ParamsError
	// }

	if req.Options[rosettaUtil.TO_ADDR] == nil {
		return nil, config.ParamsError
	}

	api := c.Config.NodeAPI(req.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)

	// get the nonce from sender
	fmt.Printf("getting nonce from address: %v\n", req.Options[rosettaUtil.SENDER_ADDR])
	senderBech32Addr := req.Options[rosettaUtil.SENDER_ADDR].(string)
	balAndNonce, err1 := rpcClient.GetBalance(rosettaUtil.ToChecksumAddr(senderBech32Addr))
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	// get the networkID (chainID) to compute the version
	networkID, err2 := rpcClient.GetNetworkId()
	if err2 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err2.Error(),
			Retriable: false,
		}
	}

	resp.Metadata[rosettaUtil.VERSION] = rosettaUtil.GetVersion(networkID)
	resp.Metadata[rosettaUtil.NONCE] = balAndNonce.Nonce + 1
	resp.Metadata[rosettaUtil.AMOUNT] = req.Options[rosettaUtil.AMOUNT]
	resp.Metadata[rosettaUtil.GAS_LIMIT] = req.Options[rosettaUtil.GAS_LIMIT]
	resp.Metadata[rosettaUtil.GAS_PRICE] = req.Options[rosettaUtil.GAS_PRICE]
	resp.Metadata[rosettaUtil.PUB_KEY] = req.Options[rosettaUtil.PUB_KEY]
	resp.Metadata[rosettaUtil.SENDER_ADDR] = senderBech32Addr
	resp.Metadata[rosettaUtil.TO_ADDR] = req.Options[rosettaUtil.TO_ADDR]

	return resp, nil
}
