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
	"strconv"
	"strings"
)

const (
	ADDRESS_TYPE        = "type"
	ADDRESS_TYPE_HEX    = "hex"
	ADDRESS_TYPE_BECH32 = "bech32"
)

type ConstructionAPIService struct {
	Config *config.Config
}

func NewConstructionAPIService(config *config.Config) *ConstructionAPIService {
	return &ConstructionAPIService{
		Config: config,
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
		return nil,&types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	txn := transaction.Transaction{
		Version:         strconv.FormatInt(int64(transactionPayload.Version),10),
		Nonce:           strconv.FormatInt(int64(transactionPayload.Nonce),10),
		Amount:          transactionPayload.Amount,
		GasPrice:        transactionPayload.GasPrice,
		GasLimit:        transactionPayload.GasLimit,
		Signature:       transactionPayload.Signature,
		SenderPubKey:    transactionPayload.PubKey,
		ToAddr:          transactionPayload.ToAddr,
		Code:            transactionPayload.Code,
		Data:            transactionPayload.Data,
	}

	hash,err1 := txn.Hash()
	if err1 != nil {
		return nil,&types.Error{
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
	return nil, nil
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
	return nil, nil
}
