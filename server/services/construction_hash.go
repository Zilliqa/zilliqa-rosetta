package services

import (
	"context"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	goZilUtil "github.com/Zilliqa/gozilliqa-sdk/util"
	rosettaUtil "github.com/Zilliqa/zilliqa-rosetta/util"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionHash /construction/hash
func (c *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	req *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	fmt.Printf("signed txn: %v\n", req.SignedTransaction)
	transactionPayload, err := provider.NewFromJson([]byte(req.SignedTransaction))
	fmt.Printf("transaction payload: %v\n", transactionPayload)
	if err != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err.Error(),
			Retriable: false,
		}
	}

	txn := transaction.NewFromPayload(transactionPayload)

	fmt.Printf("/hash - txn :%v\n\n", txn)

	// txn.ToAddr should be in zil format
	// NewFromPayload will add uneeded 0x prefix
	// remove uneeded 0x prefix
	// convert to checksum format
	txn.ToAddr = rosettaUtil.ToChecksumAddr(rosettaUtil.RemoveHexPrefix(txn.ToAddr))

	fmt.Printf("/hash - txn to address :%v\n\n", txn.ToAddr)

	hash, err1 := txn.Hash()
	if err1 != nil {
		return nil, &types.Error{
			Code:      0,
			Message:   err1.Error(),
			Retriable: false,
		}
	}

	transactionHash := goZilUtil.EncodeHex(hash)

	resp := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: transactionHash,
		},
	}
	return resp, nil
}
