package handlers

import (
	"context"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/business/data/tick"
	"github.com/qubic/qubic-http/foundation/web"
	"log"
	"net/http"
	"strconv"
)

type tickHandler struct {
	pool *qubic.Pool
}

func (h *tickHandler) GetTickInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := h.pool.Get()
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting qubic conn from pool"))
	}

	defer func() {
		if err == nil {
			log.Printf("Putting conn back to pool")
			pErr := h.pool.Put(qc)
			if pErr != nil {
				log.Printf("Putting conn back to pool failed: %s", pErr.Error())
			}
		} else {
			log.Printf("Closing conn")
			cErr := h.pool.Close(qc)
			if cErr != nil {
				log.Printf("Closing conn failed: %s", cErr.Error())
			}
		}
	}()

	res, err := tick.GetTickInfo(ctx, qc)
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick info"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (h *tickHandler) GetTickTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := h.pool.Get()
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting qubic conn from pool"))
	}

	params := web.Params(r)
	tickNr, ok := params["tick"]
	if !ok {
		return web.NewRequestError(errors.New("endpoint should have the tick number parameter in the request"), http.StatusBadRequest)
	}
	tickNumber, err := strconv.ParseInt(tickNr, 10, 32)
	if err != nil {
		return web.NewRequestError(errors.New("tick number should be a valid integer"), http.StatusBadRequest)
	}

	res, err := tick.GetTickTxs(ctx, qc, uint32(tickNumber))
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick transactions"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (h *tickHandler) GetTickData(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := h.pool.Get()
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting qubic conn from pool"))
	}

	params := web.Params(r)
	tickNr, ok := params["tick"]
	if !ok {
		return web.NewRequestError(errors.New("endpoint should have the tick number parameter in the request"), http.StatusBadRequest)
	}
	tickNumber, err := strconv.ParseInt(tickNr, 10, 32)
	if err != nil {
		return web.NewRequestError(errors.New("tick number should be a valid integer"), http.StatusBadRequest)
	}

	res, err := tick.GetTickData(ctx, qc, uint32(tickNumber))
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting tick data"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

//func (h *tickHandler) GetQuorum(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	params := web.Params(r)
//	tickNr, ok := params["tick"]
//	if !ok {
//		return web.NewRequestError(errors.New("endpoint should have the tick number parameter in the request"), http.StatusBadRequest)
//	}
//	tickNumber, err := strconv.ParseInt(tickNr, 10, 32)
//	if err != nil {
//		return web.NewRequestError(errors.New("tick number should be a valid integer"), http.StatusBadRequest)
//	}
//
//	res, err := h.opensearchClient.GetQuorum(ctx, int(tickNumber))
//	if err != nil {
//		return web.RespondError(ctx, w, errors.Wrap(err, "getting quorum data"))
//	}
//
//	return web.Respond(ctx, w, res, http.StatusOK)
//}
//
//func (h *tickHandler) GetComputors(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	params := web.Params(r)
//	epoch, ok := params["epoch"]
//	if !ok {
//		return web.NewRequestError(errors.New("endpoint should have the epoch number parameter in the request"), http.StatusBadRequest)
//	}
//	epochNr, err := strconv.ParseInt(epoch, 10, 32)
//	if err != nil {
//		return web.NewRequestError(errors.New("epoch number should be a valid integer"), http.StatusBadRequest)
//	}
//
//	res, err := h.opensearchClient.GetComputors(ctx, int(epochNr))
//	if err != nil {
//		return web.RespondError(ctx, w, errors.Wrap(err, "getting computors data"))
//	}
//
//	return web.Respond(ctx, w, res, http.StatusOK)
//}
