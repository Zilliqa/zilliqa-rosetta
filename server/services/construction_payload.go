package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
	"strconv"
)

// ConstructionPayloads /construction/payloads
// Generate an Unsigned Transaction and Signing Payloads
func (c *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	req *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {

	resp := new(types.ConstructionPayloadsResponse)
	payloads := make([]*types.SigningPayload, 0)

	// create the unsigned transaction json
	transactionJson := make(map[string]interface{})
	tx := &transaction.Transaction{}

	for _, operation := range req.Operations {
		// sender
		if operation.OperationIdentifier.Index == 0 {
			// senderAddr = operation.Account.Address
			transactionJson[rosettaUtil.SENDER_ADDR] = operation.Account.Address
		}

		// recipient
		if operation.OperationIdentifier.Index == 1 {
			// if operation.Metadata == nil {
			// 	return nil, config.ParamsError
			// }
			transactionJson[rosettaUtil.AMOUNT], _ = strconv.ParseInt(operation.Amount.Value, 10, 64)
			transactionJson[rosettaUtil.TO_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
			tx.Amount = operation.Amount.Value
			tx.ToAddr = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
			tx.SenderPubKey = rosettaUtil.RemoveHexPrefix(req.Metadata[rosettaUtil.PUB_KEY].(string))
		}
	}

	transactionJson[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(req.Metadata[rosettaUtil.PUB_KEY].(string))
	transactionJson[rosettaUtil.VERSION] = req.Metadata[rosettaUtil.VERSION]
	transactionJson[rosettaUtil.NONCE] = req.Metadata[rosettaUtil.NONCE]
	transactionJson[rosettaUtil.GAS_PRICE], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_PRICE].(string), 10, 64)
	transactionJson[rosettaUtil.GAS_LIMIT], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_LIMIT].(string), 10, 64)
	transactionJson[rosettaUtil.CODE] = ""
	transactionJson[rosettaUtil.DATA] = ""

	versionInt := int64(req.Metadata[rosettaUtil.VERSION].(float64))
	tx.Version = strconv.FormatInt(versionInt, 10)
	nonce := int64(req.Metadata[rosettaUtil.NONCE].(float64))
	tx.Nonce = strconv.FormatInt(nonce, 10)
	tx.GasPrice = req.Metadata[rosettaUtil.GAS_PRICE].(string)
	tx.GasLimit = req.Metadata[rosettaUtil.GAS_LIMIT].(string)
	tx.Code = ""
	tx.Data = ""

	unsignedTxnJson, err3 := json.Marshal(transactionJson)

	if err3 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err3.Error(),
			Retriable: false,
		}
	}

	unsignedByte, err4 := tx.Bytes()

	if err4 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err4.Error(),
			Retriable: false,
		}
	}

	ai := &types.AccountIdentifier{
		Address: transactionJson[rosettaUtil.SENDER_ADDR].(string),
	}
	signingPayload := &types.SigningPayload{
		AccountIdentifier: ai,
		Bytes:             unsignedByte, //byte array of transaction
		SignatureType:     rosettaUtil.SIGNATURE_TYPE,
	}
	payloads = append(payloads, signingPayload)
	resp.UnsignedTransaction = string(unsignedTxnJson)
	resp.Payloads = payloads
	fmt.Printf("/payloads - unsigned transaction: %v\n\n", string(unsignedTxnJson))
	return resp, nil
}
