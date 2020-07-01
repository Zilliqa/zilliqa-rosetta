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
package main

import (
	"fmt"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	router2 "github.com/Zilliqa/zilliqa-rosetta/router"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
	"net/http"
)

var log *golog.Logger

func init() {
	viper.SetConfigName("config.local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log = golog.New()
}

func main() {
	// parse config from yaml file
	cfg, err := config.ParseConfig()
	if err != nil {
		panic("Invalid config file, please check")
	}
	configString, _ := cfg.Stringify()
	log.Info("config file is: ")
	log.Info(string(configString))

	// The asserter automatically rejects incorrectly formatted
	// requests.
	networkIdentifier :=  cfg.GetNetworkIdentifier()
	asserter, err := asserter.NewServer(networkIdentifier)

	if err != nil {
		log.Fatal(err)
	}

	router := router2.NewBlockchainRouter(networkIdentifier[0], asserter,cfg)
	log.Printf("Listening on port %d\n", cfg.Rosetta.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Rosetta.Port), router))

}
