package handlers

import (
	"context"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/business/data/tick"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
	"strconv"
)

type tickHandler struct {
	pool *nodes.Pool
}

func (h *tickHandler) GetTickInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}

	res, err := tick.GetTickInfo(ctx, qc)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick info"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (h *tickHandler) GetTickTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}

	params := web.Params(r)
	tickNr := params["tick"]
	tickNumber, err := strconv.ParseInt(tickNr, 10, 32)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "parsing input"))
	}

	res, err := tick.GetTickTxs(ctx, qc, uint32(tickNumber))
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick transactions"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (h *tickHandler) GetTickData(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := tcp.NewQubicConnection(ctx, h.pool.GetRandomIP(), "21841")
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}

	params := web.Params(r)
	tickNr := params["tick"]
	tickNumber, err := strconv.ParseInt(tickNr, 10, 32)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "parsing input"))
	}

	res, err := tick.GetTickData(ctx, qc, uint32(tickNumber))
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick data"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}
