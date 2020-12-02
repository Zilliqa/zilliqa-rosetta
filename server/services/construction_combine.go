package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	goZilUtil "github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionCombine /construction/combine
// Create Network Transaction from Signatures
// sign the transaction using goZil or other out-of-band methods
// pass the result of the signature, and signed transaction in bytes as request for /combine
func (c *ConstructionAPIService) ConstructionCombine(
	ctx context.Context,
	req *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {

	// extract request params
	txnSig := util.EncodeHex(req.Signatures[0].Bytes)
	pubKey := req.Signatures[0].PublicKey.Bytes
	signedPayload := req.Signatures[0].SigningPayload.Bytes // not used for verification

	fmt.Printf("txn signature: %v\n", txnSig)
	fmt.Printf("pubKey: %v\n", pubKey)
	fmt.Printf("signedPayload: %v\n", signedPayload)

	encodedPubKey := goZilUtil.EncodeHex(pubKey)
	fmt.Printf("public key is: %v\n", encodedPubKey)
	// r := goZilUtil.DecodeHex(txnSig[0:64])
	// s := goZilUtil.DecodeHex(txnSig[64:128])

	// convert unsigned transaction to Zilliqa Transaction object
	var unsignedTxnJson map[string]interface{}
	err := json.Unmarshal([]byte(req.UnsignedTransaction), &unsignedTxnJson)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	zilliqaTransaction := &transaction.Transaction{
		Version:      fmt.Sprintf("%.0f", unsignedTxnJson["version"]),
		Nonce:        fmt.Sprintf("%.0f", unsignedTxnJson["nonce"]),
		Amount:       fmt.Sprintf("%.0f", unsignedTxnJson["amount"]),
		GasPrice:     fmt.Sprintf("%.0f", unsignedTxnJson["gasPrice"]),
		GasLimit:     fmt.Sprintf("%.0f", unsignedTxnJson["gasLimit"]),
		ToAddr:       unsignedTxnJson["toAddr"].(string),
		SenderPubKey: encodedPubKey,
		Code:         unsignedTxnJson["code"].(string),
		Data:         unsignedTxnJson["data"].(string),
		Signature:    txnSig, // signature from request param
	}

	zilliqaTransactionBytes, err2 := zilliqaTransaction.Bytes()
	if err2 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err2.Error(),
			Retriable: false,
		}
	}

	// not using signed payload from request directly
	// verify unsigned transaction + signature is indeed legit
	// also helps to verify integrity of unsigned transaction
	// signatureVerification := schnorr.Verify(pubKey, zilliqaTransactionBytes, r, s)

	// if signatureVerification == false {
	// 	return nil, config.SignatureInvalidError
	// }

	signedTxnJson := unsignedTxnJson
	signedTxnJson["signature"] = txnSig
	signedTxnJson["pubKey"] = encodedPubKey

	signedTxnBytes, err3 := json.Marshal(signedTxnJson)
	if err3 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err3.Error(),
			Retriable: false,
		}
	}

	resp := new(types.ConstructionCombineResponse)
	resp.SignedTransaction = string(signedTxnBytes)

	fmt.Printf("txn signature: %v\n", txnSig)
	fmt.Printf("signed txn json: %v\n\n", signedTxnJson)
	fmt.Printf("signed payload from request: %v\n", signedPayload)
	fmt.Printf("unsigned transaction with signature: %v\n", zilliqaTransactionBytes)
	// fmt.Printf("schnorr verify result: %v\n", signatureVerification)
	return resp, nil
}
