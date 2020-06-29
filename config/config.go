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
package config

import (
	"encoding/json"
	"github.com/spf13/viper"
)

type Rosetta struct {
	Host string
	Port int
	Version string
}

type Network struct {
	Type string
	API string
	ChainId int
}


type Config struct {
	Rosetta Rosetta
	Networks []Network
}

func ParseConfig() (*Config, error) {
	rosetta := viper.Get("rosetta").(map[string]interface{})
	host := rosetta["host"].(string)
	port := rosetta["port"].(int)
	version := rosetta["version"].(string)

	var networks []Network
	nws := viper.Get("networks").(map[string]interface{})
	for key, value := range nws {
		v := value.(map[string]interface{})
		api := v["api"].(string)
		chainId := v["chainid"].(int)
		nw := Network{
			Type: key,
			API:     api,
			ChainId: chainId,
		}

		networks = append(networks,nw)
	}

	r := Rosetta{Host: host,Port: port,Version: version}

	return &Config{
		r,networks,
	}, nil
}

func (config *Config) Stringify() ([]byte, error) {
	return json.Marshal(config)
}