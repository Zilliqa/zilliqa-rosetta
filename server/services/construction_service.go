package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	goZilUtil "github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	ADDRESS_TYPE        = "type"
	ADDRESS_TYPE_HEX    = "hex"
	ADDRESS_TYPE_BECH32 = "bech32"

	METHOD_TYPE = "method"
)

type ConstructionAPIService struct {
	Config         *config.Config
	MemPoolService *MemoryPoolAPIService
}

func NewConstructionAPIService(config *config.Config, memPoolService *MemoryPoolAPIService) *ConstructionAPIService {
	return &ConstructionAPIService{
		Config:         config,
		MemPoolService: memPoolService,
	}
}

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

// ConstructionDerive /construction/derive
func (c *ConstructionAPIService) ConstructionDerive(
	ctx context.Context,
	req *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {

	meta := req.Metadata
	pubKey := req.PublicKey.Bytes

	address := keytools.GetAddressFromPublic(pubKey)
	bech32Addr, err := bech32.ToBech32Address(address)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	resp := &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{},
	}

	if meta == nil {
		resp.AccountIdentifier.Address = bech32Addr
	} else if meta[ADDRESS_TYPE] == strings.ToLower(ADDRESS_TYPE_HEX) {
		resp.AccountIdentifier.Address = address
	} else if meta[ADDRESS_TYPE] == strings.ToLower(ADDRESS_TYPE_BECH32) {
		resp.AccountIdentifier.Address = bech32Addr
	} else {
		return nil, config.InvalidAddressTypeError
	}

	resp.Metadata = meta
	return resp, nil
}

// ConstructionHash /construction/hash
func (c *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	req *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Printf("signed txn: %v\n", req.SignedTransaction)
	transactionPayload, err := provider.NewFromJson([]byte(req.SignedTransaction))
	fmt.Printf("transaction payload: %v\n", transactionPayload)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	txn := transaction.NewFromPayload(transactionPayload)

	fmt.Printf("/hash - txn :%v\n\n", txn)

	// txn.ToAddr should be in zil format
	// NewFromPayload will add uneeded 0x prefix
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

	fmt.Printf("/hash - txn to address :%v\n\n", txn.ToAddr)

	hash, err1 := txn.Hash()
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	transactionHash := goZilUtil.EncodeHex(hash)

	resp := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: transactionHash,
		},
	}
	return resp, nil
}

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
	senderOperation.Account = &types.AccountIdentifier{
		Address: senderBech32Addr,
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
	recipientOperation.Account = &types.AccountIdentifier{
		Address: recipientBech32Addr,
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
		ai := &types.AccountIdentifier{
			Address: resp.Operations[0].Account.Address,
		}
		resp.AccountIdentifierSigners = append(resp.AccountIdentifierSigners, ai)
	}

	return resp, nil
}

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
			// transactionJson[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(operation.Metadata[rosettaUtil.PUB_KEY].(string))
		}
	}

	transactionJson[rosettaUtil.VERSION] = req.Metadata[rosettaUtil.VERSION]
	transactionJson[rosettaUtil.NONCE] = req.Metadata[rosettaUtil.NONCE]
	transactionJson[rosettaUtil.GAS_PRICE], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_PRICE].(string), 10, 64)
	transactionJson[rosettaUtil.GAS_LIMIT], _ = strconv.ParseInt(req.Metadata[rosettaUtil.GAS_LIMIT].(string), 10, 64)
	transactionJson[rosettaUtil.CODE] = ""
	transactionJson[rosettaUtil.DATA] = ""

	unsignedTxnJson, err3 := json.Marshal(transactionJson)

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
		Bytes:             unsignedTxnJson, //byte array of transaction
		SignatureType:     rosettaUtil.SIGNATURE_TYPE,
	}
	payloads = append(payloads, signingPayload)
	resp.UnsignedTransaction = string(unsignedTxnJson)
	resp.Payloads = payloads
	fmt.Printf("/payloads - unsigned transaction: %v\n\n", string(unsignedTxnJson))
	return resp, nil
}

// ConstructionPreprocess /construction/preprocess
// create a request to fetch metadata
// TODO - support contract deployment and contract call operations
// support payment operation
func (c *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	req *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	preProcessResp := &types.ConstructionPreprocessResponse{
		Options: make(map[string]interface{}),
	}
	for _, operation := range req.Operations {
		if operation.OperationIdentifier.Index == 0 {
			preProcessResp.Options[rosettaUtil.SENDER_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
		}
		if operation.OperationIdentifier.Index == 1 {
			// if operation.Metadata == nil {
			// 	return nil, config.ParamsError
			// }
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
			preProcessResp.Options[rosettaUtil.TO_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
		}
		// if operation.Metadata != nil {
		// 	preProcessResp.Options[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(operation.Metadata["senderPubKey"].(string))
		// }
	}
	preProcessResp.Options[rosettaUtil.GAS_PRICE] = rosettaUtil.GAS_PRICE_VALUE
	preProcessResp.Options[rosettaUtil.GAS_LIMIT] = rosettaUtil.GAS_LIMIT_VALUE
	return preProcessResp, nil
}

// ConstructionSubmit /construction/submit
func (c *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Printf("/submit - signed txn: %v\n\n", request.SignedTransaction)
	txStr := request.SignedTransaction
	if len(txStr) == 0 {
		return nil, config.SignedTxInvalid
	}
	pl, err := provider.NewFromJson([]byte(txStr))
	if err != nil {
		fmt.Println("error trying to new from json")
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}
	txn := transaction.NewFromPayload(pl)

	fmt.Printf("/submit - txn: %v\n\n", txn)

	// txn.ToAddr should be in zil format
	// New.FromPayload adds uneeded '0x' prefix
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

	fmt.Printf("/submit - txn: %v\n\n", txn.ToAddr)

	hash, err1 := txn.Hash()
	if err1 != nil {
		fmt.Println("error trying to convert to hash")
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	hexHash := goZilUtil.EncodeHex(hash)
	txn.ID = hexHash

	err2 := c.MemPoolService.AddTransaction(ctx, request.NetworkIdentifier, txn)
	if err2 != nil {
		fmt.Println("error trying to add transaction")
		return nil, &types.Error{
			Code:      0,
			Message:   err2.Error(),
			Retriable: false,
		}
	}

	resp := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hexHash,
		},
	}
	return resp, nil
}
