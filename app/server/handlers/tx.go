package handlers

import (
	"context"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/business/data/tx"
	"github.com/qubic/qubic-http/external/opensearch"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type txHandler struct {
	pool *nodes.Pool
	opensearchClient *opensearch.Client
}

func (h *txHandler) SendRawTx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var payload tx.SendRawTxInput
	err := web.Decode(r, &payload)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	ips := h.pool.GetMaxTargetRandomIPs(3)

	err = broadcastTxToMultiple(ctx, ips, payload)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "sending raw tx"))
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func broadcastTxToMultiple(ctx context.Context, ips []string, input tx.SendRawTxInput) error {
	nrFails := 0
	for i := 0; i < len(ips); i ++ {
		qc, err := tcp.NewQubicConnection(ctx, ips[i], "21841")
		if err != nil {
			nrFails ++
			continue
		}
		err = tx.SendRawTx(ctx, qc, input)
		if err != nil {
			nrFails ++
			continue
		}
	}

	if nrFails == len(ips) {
		return errors.New("broadcasting tx failed for all peers")
	}

	return nil
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

func (h *txHandler) GetTx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	txID, ok:= params["txID"]
	if !ok {
		return web.NewRequestError(errors.New("request should have the tx id of the address in the endpoint"), http.StatusBadRequest)
	}
	transaction, err := h.opensearchClient.GetTx(ctx, txID)
	if err != nil {
		return errors.Wrap(err, "getting tx by id")
	}

	return web.Respond(ctx, w, transaction, http.StatusOK)
}

func (h *txHandler) GetBx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	bxID, ok := params["bxID"]
	if !ok {
		return web.NewRequestError(errors.New("request should have the bx id of the address in the endpoint"), http.StatusBadRequest)
	}
	bx, err := h.opensearchClient.GetBx(ctx, bxID)
	if err != nil {
		return errors.Wrap(err, "getting bx by id")
	}

	return web.Respond(ctx, w, bx, http.StatusOK)
}


