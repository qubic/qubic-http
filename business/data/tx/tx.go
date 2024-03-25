package tx

import (
	"context"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
)

//func GetTxStatus(ctx context.Context, qc *tcp.QubicConnection, input GetTxStatusInput) (GetTxStatusOutput, error) {
//	qubicModel, err := input.toQubicModel()
//	if err != nil {
//		return GetTxStatusOutput{}, errors.Wrap(err, "converting input to qubic model")
//	}
//
//	res, err := qubic.GetTxStatus(ctx, qc, qubicModel.Tick, qubicModel.Digest, qubicModel.Signature)
//	if err != nil {
//		return GetTxStatusOutput{}, errors.Wrap(err, "getting tx status from qubic node")
//	}
//
//	var output GetTxStatusOutput
//
//	return output.fromQubicModel(res), nil
//}

func SendRawTx(ctx context.Context, qc *qubic.Client, input SendSignedTxInput) error {
	err := qc.SendRawTransaction(ctx, []byte(input.SignedTx))
	if err != nil {
		return errors.Wrap(err, "sending signed tx to qubic model")
	}

	return nil
}
