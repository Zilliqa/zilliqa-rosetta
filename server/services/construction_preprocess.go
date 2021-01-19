package services

import (
	"context"

	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionPreprocess /construction/preprocess
// create a request to fetch metadata
// TODO - support contract deployment and contract call operations
// support payment operation
func (c *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	req *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	preProcessResp := &types.ConstructionPreprocessResponse{
		Options:            make(map[string]interface{}),
		RequiredPublicKeys: []*types.AccountIdentifier{},
	}
	for _, operation := range req.Operations {
		if operation.OperationIdentifier.Index == 0 {
			preProcessResp.Options[rosettaUtil.SENDER_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
			sender := &types.AccountIdentifier{
				Address: rosettaUtil.RemoveHexPrefix(operation.Account.Address),
				Metadata: map[string]interface{}{
					rosettaUtil.Base16: operation.Account.Metadata[rosettaUtil.Base16],
				},
			}
			// to request sender public key, so we can remove pubKey from account identifier's meta data
			preProcessResp.RequiredPublicKeys = append(preProcessResp.RequiredPublicKeys, sender)
		}
		if operation.OperationIdentifier.Index == 1 {
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
			preProcessResp.Options[rosettaUtil.TO_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
		}
	}

	preProcessResp.Options[rosettaUtil.GAS_PRICE] = rosettaUtil.GAS_PRICE_VALUE
	preProcessResp.Options[rosettaUtil.GAS_LIMIT] = rosettaUtil.GAS_LIMIT_VALUE
	return preProcessResp, nil
}
