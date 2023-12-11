package tx

import (
	"context"
	"github.com/0xluk/go-qubic"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
)

func GetTxStatus(ctx context.Context, qc *tcp.QubicConnection, input GetTxStatusInput) (GetTxStatusOutput, error) {
	qubicModel, err := input.toQubicModel()
	if err != nil {
		return GetTxStatusOutput{}, errors.Wrap(err, "converting input to qubic model")
	}

	res, err := qubic.GetTxStatus(ctx, qc, qubicModel.Tick, qubicModel.Digest, qubicModel.Signature)
	if err != nil {
		return GetTxStatusOutput{}, errors.Wrap(err, "getting tx status from qubic node")
	}

	var output GetTxStatusOutput

	return output.fromQubicModel(res), nil
}

func SendSignedTx(ctx context.Context, qc *tcp.QubicConnection, input SendSignedTxInput) error {
	qubicModel, err := input.toQubicModel()
	if err != nil {
		return errors.Wrap(err, "converting input to qubic model")
	}

	err = qubic.SendRawTransaction(ctx, qc, qubicModel)
	if err != nil {
		return errors.Wrap(err, "sending signed tx to qubic model")
	}

	return nil
}
