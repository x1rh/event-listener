package eventlistener

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type EventData struct {
	Name   string
	Inputs []abi.Argument
}

type Contract struct {
	Name     string
	Address  string
	Abi      abi.ABI
	EventMap map[common.Hash]*EventData
}

func NewContract(address string, abiStr string) (*Contract, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse contract ABI")
	}

	eventMap := make(map[common.Hash]*EventData)
	for _, event := range parsedABI.Events {
		eventMap[event.ID] = &EventData{
			Name:   event.Name,
			Inputs: event.Inputs,
		}
	}

	return &Contract{
		Address:  address,
		Abi:      parsedABI,
		EventMap: eventMap,
	}, nil
}
