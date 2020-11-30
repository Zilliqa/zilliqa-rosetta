package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type AccountAPIService struct {
	Config *config.Config
}

func NewAccountAPIService(config *config.Config) *AccountAPIService {
	return &AccountAPIService{
		Config: config,
	}
}

// todo impl it
func (s *AccountAPIService) AccountCoins(context.Context, *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	return nil, nil
}

func (s *AccountAPIService) IsValidAddress(addr string) (string, error) {
	if validator.IsAddress(addr) {
		if strings.HasPrefix(addr, "0x") {
			return strings.Split(addr, "0x")[1], nil
		}

		return addr, nil
	}

	if validator.IsBech32(addr) {
		checksum, err := bech32.FromBech32Addr(addr)
		if err != nil {
			return "", err
		}
		return checksum, nil

	}

	return "", errors.New("invalid address")

}

func (s *AccountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	addr := request.AccountIdentifier.Address
	addr, err := s.IsValidAddress(addr)
	if err != nil {
		return nil, config.AddressInvalid
	}

	api := s.Config.NodeAPI(request.NetworkIdentifier.Network)
	rpcClient := provider.NewProvider(api)
	balAndNonce, err1 := rpcClient.GetBalance(addr)
	if err1 != nil {
		if err1.Error() == "-5:Account is not created" {
			return &types.AccountBalanceResponse{
				BlockIdentifier: &types.BlockIdentifier{},
				Balances: []*types.Amount{
					&types.Amount{
						Value: "0",
						Currency: &types.Currency{
							Symbol:   "ZIL",
							Decimals: 12,
							Metadata: map[string]interface{}{},
						},
						Metadata: map[string]interface{}{},
					},
				},
				Metadata: nil,
			}, nil
		} else {
			return nil, &types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}
	}

	// zilliqa has no historical lookup
	// return the request block / latest block as the block identifier
	blockIdentifier := new(types.BlockIdentifier)

	if request.BlockIdentifier != nil && request.BlockIdentifier.Index != nil {
		blockIdentifier.Index = *request.BlockIdentifier.Index

		if request.BlockIdentifier.Hash != nil {
			blockIdentifier.Hash = *request.BlockIdentifier.Hash
		} else {
			txBlock, err2 := rpcClient.GetTxBlock(fmt.Sprintf("%d", blockIdentifier.Index))
			if err2 != nil {
				return nil, &types.Error{
					Code:      0,
					Message:   err2.Error(),
					Retriable: false,
				}
			}
			blockIdentifier.Hash = txBlock.Body.BlockHash
		}
	} else {
		latestTxBlock, err3 := rpcClient.GetLatestTxBlock()

		if err3 != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err3.Error(),
				Retriable: false,
			}
		}

		blocknum, _ := strconv.ParseInt(latestTxBlock.Header.BlockNum, 10, 64)
		blockIdentifier.Index = blocknum
		blockIdentifier.Hash = latestTxBlock.Body.BlockHash
	}

	metadata := make(map[string]interface{})
	metadata["nonce"] = balAndNonce.Nonce

	return &types.AccountBalanceResponse{
		BlockIdentifier: blockIdentifier,
		Balances: []*types.Amount{
			&types.Amount{
				Value: balAndNonce.Balance,
				Currency: &types.Currency{
					Symbol:   "ZIL",
					Decimals: 12,
					Metadata: map[string]interface{}{},
				},
				Metadata: map[string]interface{}{},
			},
		},
		Metadata: metadata,
	}, nil
}
