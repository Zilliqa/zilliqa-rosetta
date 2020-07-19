package services

import (
	"context"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
	"strings"
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

func (c *ConstructionAPIService) ConstructionCombine(
	ctx context.Context,
	req *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	return nil, nil
}

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

func (c *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	req *types.ConstructionHashRequest,
) (*types.ConstructionHashResponse, *types.Error) {
	fmt.Println(req.SignedTransaction)
	transactionPayload, err := provider.NewFromJson([]byte(req.SignedTransaction))
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	txn := transaction.NewFromPayload(transactionPayload)

	hash, err1 := txn.Hash()
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	resp := &types.ConstructionHashResponse{}

	resp.TransactionHash = util.EncodeHex(hash)
	return resp, nil
}

func (c *ConstructionAPIService) ConstructionMetadata(
	ctx context.Context,
	req *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	resp := &types.ConstructionMetadataResponse{
		Metadata: make(map[string]interface{}),
	}

	if req.Options[METHOD_TYPE] != "transfer" {
		return nil, config.ParamsError
	}

	resp.Metadata["version"] = "The decimal conversion of the bitwise concatenation of CHAIN_ID and MSG_VERSION parameters"
	resp.Metadata["nonce"] = "A transaction counter in each account. This prevents replay attacks where a transaction sending eg. " +
		"20 coins from A to B can be replayed by B over and over to continually drain A's balance"
	resp.Metadata["toAddr"] = "Recipient's account address. This is represented as a String"
	resp.Metadata["amount"] = "Transaction amount to be sent to the recipent's address. This is measured in the smallest" +
		" price unit Qa (or 10^-12 Zil) in Zilliqa"
	resp.Metadata["pubKey"] = "Sender's public key of 33 bytes"
	resp.Metadata["gasPrice"] = "An amount that a sender is willing to pay per unit of gas for processing this transaction" +
		"This is measured in the smallest price unit Qa (or 10^-12 Zil) in Zilliqa"
	resp.Metadata["gasLimit"] = "The amount of gas units that is needed to be process this transaction"
	resp.Metadata["code"] = "The smart contract source code. This is present only when deploying a new contract"
	resp.Metadata["data"] = "String-ified JSON object specifying the transition parameters to be passed to a specified smart contract"
	resp.Metadata["signature"] = "An EC-Schnorr signature of 64 bytes of the entire Transaction object as stipulated above"
	resp.Metadata["priority"] = "A flag for this transaction to be processed by the DS committee"

	return resp, nil
}

func (c *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	req *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	return nil, nil
}

func (c *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	req *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	return nil, nil
}

func (c *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	req *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	return nil, nil
}

func (c *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.ConstructionSubmitResponse, *types.Error) {
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
	hash, err1 := txn.Hash()
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	hexHash := util.EncodeHex(hash)
	txn.ID = hexHash

	err2 := c.MemPoolService.AddTransaction(ctx, request.NetworkIdentifier, txn)
	if err2 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err2.Error(),
			Retriable: false,
		}
	}

	return &types.ConstructionSubmitResponse{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: hexHash},
		Metadata:              nil,
	}, nil

}
