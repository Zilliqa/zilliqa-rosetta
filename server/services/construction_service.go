package services

import (
	"context"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
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
	//fmt.Println(string(pubKey))
	//publicKey := util.DecodeHex(string(pubKey))
	//if publicKey == nil {
	//	return nil, config.PublicHexError
	//}

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
	return nil, nil
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
