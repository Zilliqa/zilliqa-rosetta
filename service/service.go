package service

import "github.com/Zilliqa/zilliqa-rosetta/config"

type Service struct {
	Config    *config.Config
}

func NewService(config *config.Config) *Service {
	service := &Service{
		Config:    config,
	}

	return service
}