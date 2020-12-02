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
	versioni := int64(req.Metadata[rosettaUtil.VERSION].(float64))
	version := strconv.FormatInt(versioni, 10)
	nonce := fmt.Sprintf("%v", req.Metadata[rosettaUtil.NONCE].(float64))
	tx := &transaction.Transaction{
		Version:  version,
		Nonce:    nonce,
		Code:     "",
		Data:     "",
		Priority: false,
	}

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
			// transactionJson[rosettaUtil.AMOUNT], _ = strconv.ParseInt(operation.Amount.Value, 10, 64)
			tx.Amount = operation.Amount.Value
			//transactionJson[rosettaUtil.TO_ADDR] = operation.Account.Address
			tx.ToAddr = operation.Account.Address
			tx.SenderPubKey = operation.Metadata[rosettaUtil.PUB_KEY].(string)
			// transactionJson[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(operation.Metadata[rosettaUtil.PUB_KEY].(string))
		}
	}

	tx.GasPrice = req.Metadata[rosettaUtil.GAS_PRICE].(string)
	tx.GasLimit = req.Metadata[rosettaUtil.GAS_LIMIT].(string)

	// transactionJson[rosettaUtil.GAS_PRICE], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_PRICE].(string), 10, 64)
	// transactionJson[rosettaUtil.GAS_LIMIT], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_LIMIT].(string), 10, 64)

	bytes, err3 := tx.Bytes()
	unsigned, _ := json.Marshal(tx)

	if err3 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err3.Error(),
			Retriable: false,
		}
	}

	ai := &types.AccountIdentifier{
		Address: transactionJson[rosettaUtil.SENDER_ADDR].(string),
	}
	signingPayload := &types.SigningPayload{
		AccountIdentifier: ai,
		Bytes:             bytes, //byte array of transaction
		SignatureType:     rosettaUtil.SIGNATURE_TYPE,
	}
	payloads = append(payloads, signingPayload)
	resp.UnsignedTransaction = string(unsigned)
	resp.Payloads = payloads
	fmt.Printf("/payloads - unsigned transaction: %v\n\n", string(unsigned))
	return resp, nil
}
