package identity

import (
	"context"
	"github.com/0xluk/go-qubic"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
)

func GetIdentity(ctx context.Context, qc *tcp.QubicConnection, id string) (GetIdentityOutput, error) {
	res, err := qubic.GetIdentity(ctx, qc, id)
	if err != nil {
		return GetIdentityOutput{}, errors.Wrap(err, "getting tx status from qubic node")
	}

	var output GetIdentityOutput

	return output.fromQubicModel(res), nil
}
