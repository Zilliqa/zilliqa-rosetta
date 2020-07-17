package router

import (
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/Zilliqa/zilliqa-rosetta/server/services"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"net/http"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(
	asserter *asserter.Asserter,
	cfg *config.Config,
) http.Handler {
	networkAPIService := services.NewNetworkAPIService(cfg)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserter,
	)

	accountAPIService := services.NewAccountAPIService(cfg)
	accountAPIController := server.NewAccountAPIController(accountAPIService, asserter)

	memoryPoolAPIService := services.NewMemoryPoolAPIService(cfg)
	memoryPoolController := server.NewMempoolAPIController(memoryPoolAPIService, asserter)

	return server.NewRouter(networkAPIController, accountAPIController, memoryPoolController)
}
