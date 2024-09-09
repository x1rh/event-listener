package eventlistener

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

type EventData struct {
	Name   string
	Inputs []abi.Argument
}

type LogHandleFunc func(ctx context.Context, event *Event) error

type Contract struct {
	Name    string
	Address string
	Abi     abi.ABI

	BlockNumber *big.Int
	Step        *big.Int

	Handle LogHandleFunc
}

func NewContract(address string, abiStr string, blockNumber, step *big.Int) (*Contract, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse contract ABI")
	}

	return &Contract{
		Address:     address,
		Abi:         parsedABI,
		BlockNumber: blockNumber,
		Step:        step,
	}, nil
}

func (c *Contract) SetLogHandler(fn LogHandleFunc) {
	c.Handle = fn
}
