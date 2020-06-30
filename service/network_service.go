/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package service

import (
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type NetWorkService struct {
	Networks []*types.NetworkIdentifier
	*Service
}

func NewNetworkService(cfg *config.Config, commonService *Service) *NetWorkService {
	nws := cfg.Networks
	var networks []*types.NetworkIdentifier

	for _, nw := range nws {
		n := &types.NetworkIdentifier{
			Blockchain: "zilliqa",
			Network:    nw.Type,
			SubNetworkIdentifier: &types.SubNetworkIdentifier{Metadata: map[string]interface{}{
				"api":     nw.API,
				"chainId": nw.ChainId,
			}},
		}

		networks = append(networks, n)
	}
	return &NetWorkService{
		Networks: networks,
		Service:  commonService,
	}
}

func (s NetWorkService) ContainsIdentifier(identifier *types.NetworkIdentifier) bool {
	for _, network := range s.Networks {
		if network.Blockchain == identifier.Blockchain && network.Network == identifier.Network {
			return true
		}
	}

	return false
}

func (s NetWorkService) NodeVersion(network string) string {
	for _, nw := range s.Config.Networks {
		if nw.Type == network {
			return nw.NodeVersion
		}
	}
	return ""
}

func (s NetWorkService) NodeAPI(network string) string {
	for _, nw := range s.Config.Networks {
		if nw.Type == network {
			return nw.API
		}
	}
	return ""
}
