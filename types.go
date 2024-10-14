package eventlistener

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ParsedData struct {
	Method *abi.Method
	Inputs map[string]any
}

type TxInfo struct {
	Hash     common.Hash
	ChainId  *big.Int
	Value    *big.Int
	From     common.Address
	To       *common.Address
	Gas      uint64
	GasPrice *big.Int
	Nonce    uint64

	Data       []byte
	ParsedData *ParsedData // parsed by Data
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (info *TxInfo) String() string {
	return prettyPrint(info)
}

type TxReceiptInfo struct {
	Events []*Event
}

type Event struct {
	Name          string
	IndexedParams []common.Hash
	Data          []byte
	Outputs       map[string]any

	// todo: remove below fields
	// BlockNumber uint64
	// TxHash      common.Hash
}
