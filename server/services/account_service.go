package services

import (
	"context"
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
	"strings"
)

type AccountAPIService struct {
	network *types.NetworkIdentifier
	Config  *config.Config
}

func NewAccountAPIService(network *types.NetworkIdentifier, config *config.Config) *AccountAPIService {
	return &AccountAPIService{
		network: network,
		Config:  config,
	}
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
			return nil,&types.Error{
				Code:      0,
				Message:   err1.Error(),
				Retriable: false,
			}
		}
	}

	return &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{},
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
		Metadata: nil,
	}, nil
}
