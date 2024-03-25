package identity

import (
	"context"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
)

func GetIdentity(ctx context.Context, qc *qubic.Client, id string) (GetIdentityOutput, error) {
	res, err := qc.GetIdentity(ctx, id)
	if err != nil {
		return GetIdentityOutput{}, errors.Wrap(err, "getting tx status from qubic node")
	}

	var output GetIdentityOutput

	return output.fromQubicModel(res), nil
}
