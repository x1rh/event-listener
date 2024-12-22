package eventlistener

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Event struct {
	Name          string
	IndexedParams []common.Hash
	Data          []byte
	Outputs       map[string]any
}

type ParsedLog struct {
	Log   *types.Log
	Event *Event // TODO:
}
