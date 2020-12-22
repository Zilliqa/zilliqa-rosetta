package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/Zilliqa/zilliqa-rosetta/config"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	// metadata
	AMOUNT          = "amount"
	CODE            = "code"
	DATA            = "data"
	GAS_LIMIT       = "gasLimit"
	GAS_PRICE       = "gasPrice"
	NONCE           = "nonce"
	PRIORITY        = "priority"
	PUB_KEY         = "pubKey"
	SIGNATURE       = "signature"
	TO_ADDR         = "toAddr"
	Base16          = "base16"
	VERSION         = "version"
	SENDER_ADDR     = "senderAddr"
	GAS_LIMIT_VALUE = "1"
	GAS_PRICE_VALUE = "2000000000"

	// others
	// set to ecdsa to bypass the signature type check for now
	SIGNATURE_TYPE = "schnorr_1"
)

// CreateRosTransaction convert zilliqa transaction object to rosetta transaction object
// "Operations" contain all balance-changing information within a transaction
// Transaction
//    - TransactionIdentifier
//    - Operation[]
//        - OperationIdentifier
//        - related_operations
//        - type
//        - status
//        - AccountIdentifier
//        - Amount
//        - metadata
func CreateRosTransaction(ctx *core.Transaction) (*types.Transaction, *types.Error) {
	rosTransaction := new(types.Transaction)
	rosTransactionIdentifier := &types.TransactionIdentifier{Hash: ctx.ID}
	rosTransaction.TransactionIdentifier = rosTransactionIdentifier
	rosOperations := make([]*types.Operation, 0)

	idx := 0

	// if transaction is contract call - transitions is present under receipt

	// ----------------------------------------------------------------------
	// payment
	// ----------------------------------------------------------------------
	if (ctx.Code == "" && ctx.Data == nil) || (ctx.Code == "" && ctx.Data == "") {
		// if transaction is payment - code and data is empty
		// -----------------
		// sender operation
		// -----------------
		senderOperation := new(types.Operation)
		senderOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx),
		}
		senderAddr := keytools.GetAddressFromPublic(util.DecodeHex(ctx.SenderPubKey))
		senderBech32Addr, err := bech32.ToBech32Address(senderAddr)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: false,
			}
		}

		status := getTransactionStatus(ctx.Receipt.Success)
		senderOperation.Type = config.OpTypeTransfer
		senderOperation.Status = &status

		senderOperation.Account = &types.AccountIdentifier{
			Address: senderBech32Addr,
			Metadata: map[string]interface{}{
				Base16: ToChecksumAddr(senderBech32Addr),
			},
		}
		// deduct from sender account
		// add negative sign
		senderOperation.Amount = CreateRosAmount(ctx.Amount, true)

		// -------------------
		// recipient operation
		// -------------------
		recipientOperation := new(types.Operation)
		recipientOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx + 1),
		}
		recipientOperation.RelatedOperations = []*types.OperationIdentifier{
			{
				Index: int64(idx),
			},
		}

		recipientBech32Addr, err := bech32.ToBech32Address(ctx.ToAddr)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: false,
			}
		}

		recipientChecksum, _ := bech32.FromBech32Addr(recipientBech32Addr)

		recipientStatus := getTransactionStatus(ctx.Receipt.Success)
		recipientOperation.Type = config.OpTypeTransfer
		recipientOperation.Status = &recipientStatus
		recipientOperation.Account = &types.AccountIdentifier{
			Address: recipientBech32Addr,
			Metadata: map[string]interface{}{
				Base16: ToChecksumAddr(recipientChecksum),
			},
		}

		recipientOperation.Amount = CreateRosAmount(ctx.Amount, false)
		recipientOperation.Metadata = createMetadata(ctx)

		rosOperations = append(rosOperations, senderOperation, recipientOperation)
		rosTransaction.Operations = rosOperations
		return rosTransaction, nil

	} else if ctx.ToAddr == "0000000000000000000000000000000000000000" {
		// ----------------------------------------------------------------------
		// contract deployment
		// ----------------------------------------------------------------------

		// -----------------
		// sender operation
		// -----------------
		senderOperation := new(types.Operation)
		senderOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx),
		}
		senderAddr := keytools.GetAddressFromPublic(util.DecodeHex(ctx.SenderPubKey))
		senderBech32Addr, err := bech32.ToBech32Address(senderAddr)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: false,
			}
		}

		status := getTransactionStatus(ctx.Receipt.Success)
		senderOperation.Type = config.OpTypeContractDeployment
		senderOperation.Status = &status
		senderOperation.Account = &types.AccountIdentifier{
			Address: senderBech32Addr,
			Metadata: map[string]interface{}{
				Base16: ToChecksumAddr(senderBech32Addr),
			},
		}
		// deduct from sender account
		// add negative sign
		senderOperation.Amount = CreateRosAmount(ctx.Amount, true)
		senderOperation.Metadata = createMetadata(ctx)

		rosOperations = append(rosOperations, senderOperation)
		rosTransaction.Operations = rosOperations
		return rosTransaction, nil

	} else if ctx.Data != nil || ctx.Data != "" {
		// ----------------------------------------------------------------------
		// contract call
		// ----------------------------------------------------------------------

		//  check if transition is a smart contract deposit , e.g. contains AddFunds
		isSmartContractDeposit := false
		for _, transition := range ctx.Receipt.Transitions {
			if transition.Msg.Tag == "AddFunds" || transition.Msg.Tag == "" {
				if transition.Msg.Amount != "0" {
					isSmartContractDeposit = true
					break
				}
			}
		}

		// -----------------
		// initiator operation
		// -----------------
		initiatorOperation := new(types.Operation)
		initiatorOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx),
		}

		initiatorAddr := keytools.GetAddressFromPublic(util.DecodeHex(ctx.SenderPubKey))
		initiatorBech32Addr, err := bech32.ToBech32Address(initiatorAddr)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: false,
			}
		}

		status := getTransactionStatus(ctx.Receipt.Success)
		initiatorOperation.Type = config.OpTypeContractCall
		initiatorOperation.Status = &status
		initiatorOperation.Account = &types.AccountIdentifier{
			Address: initiatorBech32Addr,
			Metadata: map[string]interface{}{
				Base16: ToChecksumAddr(initiatorBech32Addr),
			},
		}

		// if it is not smart contract deposit, ie, no transition, means there is only one operation
		// add metadata if only one and only operation
		if isSmartContractDeposit {
			initiatorOperation.Amount = CreateRosAmount("0", true)
		} else {
			initiatorOperation.Amount = CreateRosAmount(ctx.Amount, true)
			initiatorOperation.Metadata = createMetadataContractCall(ctx)
		}

		rosOperations = append(rosOperations, initiatorOperation)
		idx += 1

		for _, transition := range ctx.Receipt.Transitions {
			if transition.Msg.Tag == "AddFunds" || transition.Msg.Tag == "" {
				// -----------------
				// from operation
				// -----------------
				fromOperation := new(types.Operation)
				fromOperation.OperationIdentifier = &types.OperationIdentifier{
					Index: int64(idx),
				}
				fromOperation.RelatedOperations = []*types.OperationIdentifier{
					{
						Index: int64(idx - 1),
					},
				}

				fromBech32Addr, err := bech32.ToBech32Address(transition.Addr)

				if err != nil {
					return nil, &types.Error{
						Code:      0,
						Message:   err.Error(),
						Retriable: false,
					}
				}

				status := getTransactionStatus(ctx.Receipt.Success)
				fromOperation.Type = config.OpTypeContractCallTransfer
				fromOperation.Status = &status
				fromOperation.Account = &types.AccountIdentifier{
					Address: fromBech32Addr,
					Metadata: map[string]interface{}{
						Base16: ToChecksumAddr(fromBech32Addr),
					},
				}
				fromOperation.Amount = CreateRosAmount(transition.Msg.Amount, true)
				// fromOperation.Metadata = createMetadataContractCall(ctx)

				rosOperations = append(rosOperations, fromOperation)
				idx += 1

				// -----------------
				// to operation
				// -----------------
				toOperation := new(types.Operation)
				toOperation.OperationIdentifier = &types.OperationIdentifier{
					Index: int64(idx),
				}
				toOperation.RelatedOperations = []*types.OperationIdentifier{
					{
						Index: int64(idx - 1),
					},
				}
				toBech32Addr, err := bech32.ToBech32Address(transition.Msg.Recipient)

				if err != nil {
					return nil, &types.Error{
						Code:      0,
						Message:   err.Error(),
						Retriable: false,
					}
				}

				toStatus := getTransactionStatus(ctx.Receipt.Success)
				toOperation.Type = config.OpTypeContractCallTransfer
				toOperation.Status = &toStatus
				toOperation.Account = &types.AccountIdentifier{
					Address: toBech32Addr,
					Metadata: map[string]interface{}{
						Base16: ToChecksumAddr(toBech32Addr),
					},
				}
				toOperation.Amount = CreateRosAmount(transition.Msg.Amount, false)
				toOperation.Metadata = createMetadataContractCall(ctx)

				rosOperations = append(rosOperations, toOperation)
				idx += 1
			}
		}

		rosTransaction.Operations = rosOperations
		return rosTransaction, nil

	} else {
		// ----------------------------------------------------------------------
		// generic case
		// ----------------------------------------------------------------------
		recipientOperation := new(types.Operation)
		recipientOperation.OperationIdentifier = &types.OperationIdentifier{
			Index: int64(idx + 1),
		}
		recipientOperation.RelatedOperations = []*types.OperationIdentifier{
			{
				Index: int64(idx),
			},
		}

		recipientBech32Addr, err := bech32.ToBech32Address(ctx.ToAddr)

		if err != nil {
			return nil, &types.Error{
				Code:      0,
				Message:   err.Error(),
				Retriable: false,
			}
		}

		receipientStatus := getTransactionStatus(ctx.Receipt.Success)
		recipientOperation.Type = config.OpTypeTransfer
		recipientOperation.Status = &receipientStatus
		recipientOperation.Account = &types.AccountIdentifier{
			Address: recipientBech32Addr,
			Metadata: map[string]interface{}{
				Base16: ToChecksumAddr(recipientBech32Addr),
			},
		}

		recipientOperation.Amount = CreateRosAmount(ctx.Amount, false)
		recipientOperation.Metadata = createMetadata(ctx)

		rosOperations = append(rosOperations, recipientOperation)
		rosTransaction.Operations = rosOperations
		return rosTransaction, nil
	}
}

// create the amount identifier
// if isNegative is true, indicates that the stated amount is deducted
func CreateRosAmount(amount string, isNegative bool) *types.Amount {
	if isNegative && amount != "0" {
		amount = fmt.Sprintf("-%s", amount)
	}
	return &types.Amount{
		Value: amount,
		Currency: &types.Currency{
			Symbol:   "ZIL",
			Decimals: 12,
		},
	}
}

func createMetadata(ctx *core.Transaction) map[string]interface{} {
	metadata := make(map[string]interface{})

	if ctx.Code != "" {
		metadata["code"] = ctx.Code
	}

	if ctx.Data != nil {
		metadata["data"] = ctx.Data
	}

	metadata["gasLimit"] = ctx.GasLimit
	metadata["gasPrice"] = ctx.GasPrice
	metadata["nonce"] = ctx.Nonce
	metadata["signature"] = ctx.Signature
	metadata["receipt"] = ctx.Receipt
	metadata["senderPubKey"] = ctx.SenderPubKey
	metadata["version"] = ctx.Version
	return metadata
}

func createMetadataContractCall(ctx *core.Transaction) map[string]interface{} {
	metadata := createMetadata(ctx)
	metadata["contractAddress"] = ctx.ToAddr
	return metadata
}

func getTransactionStatus(status bool) string {
	if status == true {
		return config.StatusSuccess.Status
	}
	return config.StatusFailed.Status
}

func RemoveHexPrefix(address string) string {
	if strings.HasPrefix(address, "0x") {
		return strings.Split(address, "0x")[1]
	}
	return address
}

func ToCoreTransaction(txn *transaction.Transaction) *core.Transaction {
	return &core.Transaction{
		ID:           txn.ID,
		Version:      txn.Version,
		Nonce:        txn.Nonce,
		Amount:       txn.Amount,
		GasPrice:     txn.GasPrice,
		GasLimit:     txn.GasLimit,
		Signature:    txn.Signature,
		SenderPubKey: txn.SenderPubKey,
		ToAddr:       txn.ToAddr,
		Code:         txn.Code,
		Data:         txn.Data,
		Priority:     txn.Priority,
	}
}

func GetVersion(chainIdStr string) int {
	chainID, _ := strconv.Atoi(chainIdStr)
	return int(util.Pack(chainID, 1))
}

func ToChecksumAddr(addr string) string {
	if validator.IsBech32(addr) {
		checksum, err := bech32.FromBech32Addr(addr)
		if err != nil {
			return RemoveHexPrefix(addr)
		}
		return checksum
	}
	return RemoveHexPrefix(addr)
}
