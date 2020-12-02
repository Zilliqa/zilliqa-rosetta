package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionParse /construction/parse
// Parse is called on both unsigned and signed transactions to understand the intent of the formulated transaction.
// This is run as a sanity check before signing (after `/construction/payloads`) and before broadcast (after `/construction/combine`).
func (c *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	req *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {

	// convert transaction to Zilliqa Transaction object
	var txnJson map[string]interface{}
	err := json.Unmarshal([]byte(req.Transaction), &txnJson)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	rosOperations := make([]*types.Operation, 0)

	// zilliqaTransaction := &transaction.Transaction{
	// 	Version:      fmt.Sprintf("%.0f", txnJson["version"]),
	// 	Nonce:        fmt.Sprintf("%.0f", txnJson["nonce"]),
	// 	Amount:       fmt.Sprintf("%.0f", txnJson["amount"]),
	// 	GasPrice:     fmt.Sprintf("%.0f", txnJson["gasPrice"]),
	// 	GasLimit:     fmt.Sprintf("%.0f", txnJson["gasLimit"]),
	// 	ToAddr:       rosettaUtil.ToChecksumAddr(txnJson["toAddr"].(string)),
	// 	SenderPubKey: txnJson["pubKey"].(string),
	// 	Code:         txnJson["code"].(string),
	// 	Data:         txnJson["data"].(string),
	// }

	senderBech32Addr := txnJson["senderAddr"].(string)
	recipientBech32Addr := txnJson["toAddr"].(string)
	amount := fmt.Sprintf("%.0f", txnJson["amount"])
	// sender operation
	idx := 0
	senderOperation := new(types.Operation)
	senderOperation.OperationIdentifier = &types.OperationIdentifier{
		Index: int64(idx),
	}
	senderOperation.Type = config.OpTypeTransfer
	senderOperation.Status = new(string)
	meta := make(map[string]interface{},1)
	meta[rosettaUtil.PUB_KEY] = txnJson[rosettaUtil.PUB_KEY]

	senderOperation.Account = &types.AccountIdentifier{
		Address: senderBech32Addr,
		Metadata: meta,
	}
	senderOperation.Amount = rosettaUtil.CreateRosAmount(amount, true)

	// recipient operation
	recipientOperation := new(types.Operation)
	recipientOperation.OperationIdentifier = &types.OperationIdentifier{
		Index: int64(idx + 1),
	}
	recipientOperation.RelatedOperations = []*types.OperationIdentifier{
		{
			Index: int64(idx),
		},
	}

	recipientOperation.Type = config.OpTypeTransfer
	recipientOperation.Status = new(string)
	//meta = make(map[string]interface{},1)

	recipientOperation.Account = &types.AccountIdentifier{
		Address: recipientBech32Addr,
		//Metadata: meta,
	}

	recipientOperation.Amount = rosettaUtil.CreateRosAmount(amount, false)

	rosOperations = append(rosOperations, senderOperation, recipientOperation)

	if req.Signed {
		if txnJson["signature"] == nil || txnJson["signature"] == "" {
			return nil, config.SignatureInvalidError
		}
		// zilliqaTransaction.Signature = txnJson["signature"].(string)
	}

	// convert to rosetta transaction object
	// rosTransaction, err2 := rosettaUtil.CreateRosTransaction(rosettaUtil.ToCoreTransaction(zilliqaTransaction))
	// if err2 != nil {
	// 	return nil, err2
	// }

	resp := &types.ConstructionParseResponse{
		AccountIdentifierSigners: make([]*types.AccountIdentifier, 0),
		Operations:               rosOperations,
		Metadata:                 make(map[string]interface{}),
	}

	// set all the operation status to success
	// for _, operations := range rosTransaction.Operations {
	// 	operations.Status = config.StatusSuccess.Status
	// }

	// resp.Operations = rosTransaction.Operations

	if req.Signed {
		// txnJson is a signed transaction
		// assume sender is signer
		// remember to put metadata here
		meta := make(map[string]interface{},1)
		meta[rosettaUtil.PUB_KEY] = txnJson[rosettaUtil.PUB_KEY]
		ai := &types.AccountIdentifier{
			Address: resp.Operations[0].Account.Address,
			Metadata: meta,
		}
		resp.AccountIdentifierSigners = append(resp.AccountIdentifierSigners, ai)
	}

	return resp, nil
}
