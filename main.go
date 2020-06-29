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
	"github.com/Zilliqa/zilliqa-rosetta/controller"
	service2 "github.com/Zilliqa/zilliqa-rosetta/service"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
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

	// register controller
	app := iris.New()
	networkService := service2.NewNetworkService(cfg)
	controller.NewNetworkController(app, networkService)
	_ = app.Run(iris.Addr(fmt.Sprintf("%s:%d", cfg.Rosetta.Host, cfg.Rosetta.Port)))

}
