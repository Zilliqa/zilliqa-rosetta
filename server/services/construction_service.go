package services

import (
	"github.com/Zilliqa/zilliqa-rosetta/config"
)

const (
	ADDRESS_TYPE        = "type"
	ADDRESS_TYPE_HEX    = "hex"
	ADDRESS_TYPE_BECH32 = "bech32"

	METHOD_TYPE = "method"
)

type ConstructionAPIService struct {
	Config *config.Config
}

func NewConstructionAPIService(config *config.Config) *ConstructionAPIService {
	return &ConstructionAPIService{
		Config: config,
	}
}
