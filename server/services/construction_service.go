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
	schnorr "github.com/Zilliqa/gozilliqa-sdk/schnorr"
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

	r := goZilUtil.DecodeHex(txnSig[0:64])
	s := goZilUtil.DecodeHex(txnSig[64:128])

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
		ToAddr:       rosettaUtil.ToChecksumAddr(unsignedTxnJson["toAddr"].(string)),
		SenderPubKey: unsignedTxnJson["pubKey"].(string),
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
	signatureVerification := schnorr.Verify(pubKey, zilliqaTransactionBytes, r, s)

	if signatureVerification == false {
		return nil, config.SignatureInvalidError
	}

	// add signature to unsigned transaction json
	unsignedTxnJson["signature"] = txnSig
	signedTxnJson, err3 := json.Marshal(unsignedTxnJson)
	if err3 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err3.Error(),
			Retriable: false,
		}
	}

	resp := new(types.ConstructionCombineResponse)
	resp.SignedTransaction = string(signedTxnJson)

	fmt.Printf("txn signature: %v\n", txnSig)
	fmt.Printf("signed payload from request: %v\n", signedPayload)
	fmt.Printf("unsigned transaction with signature: %v\n", zilliqaTransactionBytes)
	fmt.Printf("schnorr verify result: %v\n", signatureVerification)
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

	resp := new(types.ConstructionDeriveResponse)

	if meta == nil {
		resp.Address = bech32Addr
	} else if meta[ADDRESS_TYPE] == strings.ToLower(ADDRESS_TYPE_HEX) {
		resp.Address = address
	} else if meta[ADDRESS_TYPE] == strings.ToLower(ADDRESS_TYPE_BECH32) {
		resp.Address = bech32Addr
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
	transactionPayload, err := provider.NewFromJson([]byte(req.SignedTransaction))
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	txn := transaction.NewFromPayload(transactionPayload)

	// txn.ToAddr should be in zil format
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

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

	if req.Options[rosettaUtil.PUB_KEY] == nil {
		return nil, config.ParamsError
	}

	if req.Options[rosettaUtil.TO_ADDR] == nil {
		return nil, config.ParamsError
	}

	api := c.Config.NodeAPI(req.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)

	// get the nonce from sender
	senderPubKey := req.Options[rosettaUtil.PUB_KEY].(string)
	senderAddr := keytools.GetAddressFromPublic(util.DecodeHex(senderPubKey))

	balAndNonce, err1 := rpcClient.GetBalance(senderAddr)
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
	resp.Metadata[rosettaUtil.PUB_KEY] = senderPubKey
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

	zilliqaTransaction := &transaction.Transaction{
		Version:      fmt.Sprintf("%.0f", txnJson["version"]),
		Nonce:        fmt.Sprintf("%.0f", txnJson["nonce"]),
		Amount:       fmt.Sprintf("%.0f", txnJson["amount"]),
		GasPrice:     fmt.Sprintf("%.0f", txnJson["gasPrice"]),
		GasLimit:     fmt.Sprintf("%.0f", txnJson["gasLimit"]),
		ToAddr:       rosettaUtil.ToChecksumAddr(txnJson["toAddr"].(string)),
		SenderPubKey: txnJson["pubKey"].(string),
		Code:         txnJson["code"].(string),
		Data:         txnJson["data"].(string),
	}

	if req.Signed {
		if txnJson["signature"] == nil || txnJson["signature"] == "" {
			return nil, config.SignatureInvalidError
		}
		zilliqaTransaction.Signature = txnJson["signature"].(string)
	}

	// convert to rosetta transaction object
	rosTransaction, err2 := rosettaUtil.CreateRosTransaction(rosettaUtil.ToCoreTransaction(zilliqaTransaction))
	if err2 != nil {
		return nil, err2
	}

	resp := &types.ConstructionParseResponse{
		Signers:    make([]string, 0),
		Operations: []*types.Operation{},
		Metadata:   make(map[string]interface{}),
	}

	// set all the operation status to success
	for _, operations := range rosTransaction.Operations {
		operations.Status = config.StatusSuccess.Status
	}

	resp.Operations = rosTransaction.Operations

	if req.Signed {
		// txnJson is a signed transaction
		// assume sender is signer
		resp.Signers = append(resp.Signers, resp.Operations[0].Account.Address)
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
	var senderAddr string
	transactionJson := make(map[string]interface{})

	for _, operation := range req.Operations {
		// sender
		if operation.OperationIdentifier.Index == 0 {
			senderAddr = operation.Account.Address
		}

		// recipient
		if operation.OperationIdentifier.Index == 1 {
			if operation.Metadata == nil {
				return nil, config.ParamsError
			}
			transactionJson[rosettaUtil.VERSION] = operation.Metadata[rosettaUtil.VERSION]
			transactionJson[rosettaUtil.NONCE] = operation.Metadata[rosettaUtil.NONCE]
			transactionJson[rosettaUtil.AMOUNT], _ = strconv.ParseInt(operation.Amount.Value, 10, 64)
			transactionJson[rosettaUtil.TO_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
			transactionJson[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(operation.Metadata[rosettaUtil.PUB_KEY].(string))
			transactionJson[rosettaUtil.GAS_PRICE], _ = strconv.ParseInt(operation.Metadata[rosettaUtil.GAS_PRICE].(string), 10, 64)
			transactionJson[rosettaUtil.GAS_LIMIT], _ = strconv.ParseInt(operation.Metadata[rosettaUtil.GAS_LIMIT].(string), 10, 64)
		}
	}

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

	signingPayload := &types.SigningPayload{
		Address:       senderAddr,
		Bytes:         unsignedTxnJson, //byte array of transaction
		SignatureType: rosettaUtil.SIGNATURE_TYPE,
	}
	payloads = append(payloads, signingPayload)
	resp.UnsignedTransaction = string(unsignedTxnJson)
	resp.Payloads = payloads
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
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
		}
		if operation.OperationIdentifier.Index == 1 {
			if operation.Metadata == nil {
				return nil, config.ParamsError
			}
			preProcessResp.Options[rosettaUtil.AMOUNT] = operation.Amount.Value
			preProcessResp.Options[rosettaUtil.TO_ADDR] = rosettaUtil.RemoveHexPrefix(operation.Account.Address)
		}
		if operation.Metadata != nil {
			preProcessResp.Options[rosettaUtil.PUB_KEY] = rosettaUtil.RemoveHexPrefix(operation.Metadata["senderPubKey"].(string))
		}
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
	txStr := request.SignedTransaction
	if len(txStr) == 0 {
		return nil, config.SignedTxInvalid
	}
	pl, err := provider.NewFromJson([]byte(txStr))
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}
	txn := transaction.NewFromPayload(pl)

	// txn.ToAddr should be in zil format
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

	hash, err1 := txn.Hash()
	if err1 != nil {
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
