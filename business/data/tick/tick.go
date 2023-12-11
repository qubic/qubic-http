package tick

import (
	"context"
	"github.com/0xluk/go-qubic"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
)

func GetTickData(ctx context.Context, qc *tcp.QubicConnection, tick uint32) (GetTickDataOutput, error) {
	res, err := qubic.GetTickData(ctx, qc, tick)
	if err != nil {
		return GetTickDataOutput{}, errors.Wrap(err, "getting tick data from qubic")
	}

	var output GetTickDataOutput
	return output.fromQubicModel(res), nil
}

func GetTickTxs(ctx context.Context, qc *tcp.QubicConnection, tick uint32) (GetTickTransactionsOutput, error) {
	res, err := qubic.GetTickTransactions(ctx, qc, tick)
	if err != nil {
		return GetTickTransactionsOutput{}, errors.Wrap(err, "getting tick transactions from qubic")
	}

	var output GetTickTransactionsOutput
	return output.fromQubicModel(res), nil
}

func GetTickInfo(ctx context.Context, qc *tcp.QubicConnection) (GetTickInfoOutput, error) {
	res, err := qubic.GetTickInfo(ctx, qc)
	if err != nil {
		return GetTickInfoOutput{}, errors.Wrap(err, "get tick info from qubic node")
	}

	var output GetTickInfoOutput
	return output.fromQubicModel(res), nil
}
