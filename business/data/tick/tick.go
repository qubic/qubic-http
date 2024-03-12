package tick

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/external/opensearch"
)

func GetTickData(ctx context.Context, qc *qubic.Client, tick uint32) (GetTickDataOutput, error) {
	res, err := qc.GetTickData(ctx, tick)
	if err != nil {
		return GetTickDataOutput{}, errors.Wrap(err, "getting tick data from qubic")
	}

	var output GetTickDataOutput
	return output.fromQubicModel(res), nil
}

func GetTickDataV2(ctx context.Context, osClient *opensearch.Client, tick uint32) (GetTickDataOutput, error) {
	res, err := osClient.GetTickData(ctx, tick)
	if err != nil {
		return GetTickDataOutput{}, errors.Wrap(err, "getting tick data from opensearch")
	}

	var output GetTickDataOutput
	return output.fromOpensearchModel(res), nil
}

func GetTickTxs(ctx context.Context, qc *qubic.Client, tick uint32) (GetTickTransactionsOutput, error) {
	qubicRes, err := qc.GetTickTransactions(ctx, tick)
	if err != nil {
		return GetTickTransactionsOutput{}, errors.Wrap(err, "getting tick transactions from qubic")
	}

	var output GetTickTransactionsOutput
	res, err := output.fromQubicModel(qubicRes)
	if err != nil {
		return GetTickTransactionsOutput{}, errors.Wrap(err, "parsing tick transactions")
	}
	return res, nil
}

func GetTickInfo(ctx context.Context, qc *qubic.Client) (GetTickInfoOutput, error) {
	res, err := qc.GetTickInfo(ctx)
	if err != nil {
		return GetTickInfoOutput{}, errors.Wrap(err, "get tick info from qubic node")
	}

	var output GetTickInfoOutput
	return output.fromQubicModel(res), nil
}
