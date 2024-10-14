package eventlistener

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func LogHandlerOnlyTopic(el *EventListener, eventHandlers ...LogEventHandleFunc) LogHandleFunc {
	return func(ctx context.Context, txLog *types.Log) error {
		// TODO: only handle topic
		// topic[0] is always a signature when a event is topic
		eventSignature := txLog.Topics[0]
		event, err := el.Contract.Abi.EventByID(eventSignature)
		if err != nil {
			// slog.Error("fail to get even", slog.Any("err", err))
			return errors.Wrap(err, "fail to get event")
		}

		eventInfo := &Event{
			Name:          event.Name,
			IndexedParams: make([]common.Hash, len(txLog.Topics)-1),
			Data:          txLog.Data,
			Outputs:       nil,
			// BlockNumber:   txLog.BlockNumber,
			// TxHash:        txLog.TxHash,
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
			err = el.Contract.Abi.UnpackIntoMap(outputDataMap, event.Name, txLog.Data)
			if err != nil {
				// slog.Error("fail to unpack", slog.Any("err", err))
				return errors.Wrap(err, "fail to unpack")
			}
			eventInfo.Outputs = outputDataMap
		}

		slog.Debug(
			"hanle",
			slog.String("chainName", el.Config.ChainName),
			slog.String("contractName", el.Contract.Name),
			slog.String("ContractAddress", el.Contract.Address),
			slog.Any("block number", txLog.BlockNumber),
		)

		for _, handler := range eventHandlers {
			if err := handler(ctx, txLog, eventInfo); err != nil {
				return errors.Wrap(err, "call event handler error")
			}
		}

		return nil
	}

}
