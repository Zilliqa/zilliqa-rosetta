package services

import (
	"context"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	goZilUtil "github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionSubmit /construction/submit
func (c *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Printf("/submit - signed txn: %v\n\n", request.SignedTransaction)
	txStr := request.SignedTransaction
	if len(txStr) == 0 {
		return nil, config.SignedTxInvalid
	}
	pl, err := provider.NewFromJson([]byte(txStr))
	if err != nil {
		fmt.Println("error trying to new from json")
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}
	txn := transaction.NewFromPayload(pl)

	fmt.Printf("/submit - txn: %v\n\n", txn)

	// txn.ToAddr should be in zil format
	// New.FromPayload adds uneeded '0x' prefix
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

	fmt.Printf("/submit - txn: %v\n\n", txn.ToAddr)

	hash, err1 := txn.Hash()
	if err1 != nil {
		fmt.Println("error trying to convert to hash")
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	hexHash := goZilUtil.EncodeHex(hash)
	txn.ID = hexHash

	err2 := c.MemPoolService.AddTransaction(ctx, request.NetworkIdentifier, txn)
	if err2 != nil {
		fmt.Println("error trying to add transaction")
		return nil, &types.Error{
			Code:      0,
			Message:   err2.Error(),
			Retriable: false,
		}
	}

	resp := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hexHash,
		},
	}
	return resp, nil
}
