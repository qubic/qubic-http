package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/business/data/tx"
	"github.com/qubic/qubic-http/external/opensearch"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type txHandler struct {
	pool             *qubic.Pool
	opensearchClient *opensearch.Client
}

func (h *txHandler) SendRawTx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var payload tx.SendSignedTxInput
	err := web.Decode(r, &payload)
	if err != nil {
		return web.RespondError(ctx, w, err)
	}

	nrSuccess := broadcastTxToMultiple(ctx, h.pool, payload)
	return web.Respond(ctx, w, struct{ Message string }{Message: fmt.Sprintf("Transaction broadcasted to %d peers", nrSuccess)}, http.StatusOK)
}

func broadcastTxToMultiple(ctx context.Context, pool *qubic.Pool, input tx.SendSignedTxInput) int {
	nrSuccess := 0
	for i := 0; i < 3; i++ {
		func() {
			client, err := pool.Get()
			if err != nil {
				return
			}

			decoded, err := hex.DecodeString(input.SignedTx)
			if err != nil {
				pool.Close(client)
				return
			}

			err = client.SendRawTransaction(ctx, decoded)
			if err != nil {
				pool.Close(client)
				return
			}
			pool.Put(client)
			nrSuccess++
		}()
	}

	return nrSuccess
}

//func (h *txHandler) GetTxStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
//	if err != nil {
//		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
//	}
//
//	var payload tx.GetTxStatusInput
//	err = web.Decode(r, &payload)
//	if err != nil {
//		return web.RespondError(ctx, w, err)
//	}
//
//	output, err := tx.GetTxStatus(ctx, qc, payload)
//	if err != nil {
//		return web.RespondError(ctx, w, err)
//	}
//
//	return web.Respond(ctx, w, output, http.StatusOK)
//}

func (h *txHandler) GetTx(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	txID, ok := params["txID"]
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
