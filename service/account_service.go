package service

import (
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"strings"
)

type AccountService struct {
	*Service
}

func NewAccountService(cfg *config.Config, commonService *Service) *AccountService {
	return &AccountService{
		commonService,
	}
}

func (s *AccountService) IsValidAddress(addr string) (string, error) {
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
