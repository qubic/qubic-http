package handlers

import (
	"context"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/business/data/tx"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type txHandler struct {
	pool *nodes.Pool
}

func (h *txHandler) SendRawTx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}

	var payload tx.SendRawTxInput
	err = web.Decode(r, &payload)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	err = tx.SendRawTx(ctx, qc, payload)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "sending raw tx"))
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (h *txHandler) GetTxStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}

	var payload tx.GetTxStatusInput
	err = web.Decode(r, &payload)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	output, err := tx.GetTxStatus(ctx, qc, payload)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	return web.Respond(ctx, w, output, http.StatusOK)
}
