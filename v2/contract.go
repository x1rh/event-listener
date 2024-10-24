package eventlistener

import (
	"context"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type EventData struct {
	Name   string
	Inputs []abi.Argument
}

type LogHandleFunc func(ctx context.Context, txLog *types.Log) error
type LogEventHandleFunc func(ctx context.Context, txLog *types.Log, event *Event) error

type ILogHandler interface {
	LogHandler() LogHandleFunc
	EventHandler() LogEventHandleFunc
}

type IContract interface {
	HandleLog() error
	HandleEvent() error
}

type Contract struct {
	Name    string
	Address string
	Abi     abi.ABI

	BlockNumber *big.Int
	Step        *big.Int
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

func (c *Contract) HandleLog(ctx context.Context, txLog *types.Log) error {
	// TODO: only handle topic
	// topic[0] is always a signature when a event is topic
	eventSignature := txLog.Topics[0]
	event, err := c.Abi.EventByID(eventSignature)
	if err != nil {
		return errors.Wrap(err, "fail to get event")
	}

	eventInfo := &Event{
		Name:          event.Name,
		IndexedParams: make([]common.Hash, len(txLog.Topics)-1),
		Data:          txLog.Data,
		Outputs:       nil,
	}
	slog.Debug("event", slog.Any("event", event))

	// topic[1:] is other indexed params in event
	if len(txLog.Topics) > 1 {
		for i, param := range txLog.Topics[1:] {
			eventInfo.IndexedParams[i] = param
			slog.Debug("", event.Inputs[i].Name, common.HexToAddress(param.Hex()))
		}
	}
	if len(txLog.Data) > 0 {
		outputDataMap := make(map[string]interface{})
		err = c.Abi.UnpackIntoMap(outputDataMap, event.Name, txLog.Data)
		if err != nil {
			return errors.Wrap(err, "fail to unpack")
		}
		eventInfo.Outputs = outputDataMap
	}

	slog.Debug(
		"hanle",
		// slog.String("chainName", el.Config.ChainName),
		slog.String("contractName", c.Name),
		slog.String("ContractAddress", c.Address),
		slog.Any("block number", txLog.BlockNumber),
	)

	if err := c.HandleEvent(ctx, txLog, eventInfo); err != nil {
		return errors.Wrap(err, "call event handler error")
	}

	return nil
}

func (c *Contract) HandleEvent(ctx context.Context, txLog *types.Log, event *Event) error {
	// use appctx do something
	_ = appctx

	slog.Info("eventInfo", slog.Any("event", event))
	switch event.Name {
	case "TokenCreated":
		var l LogTokenCreated
		if err := c.Abi.UnpackIntoInterface(&l, event.Name, event.Data); err != nil {
			return errors.Wrap(err, "fail to unpack log")
		}
		// handle indexed topic
		l.Owner = HashToAddress(event.IndexedParams[0])
		l.Token = HashToAddress(event.IndexedParams[1])
		slog.Info("TokenCreated event", slog.Any("event", l))

	default:
		// do nothing
	}
	return nil
}
