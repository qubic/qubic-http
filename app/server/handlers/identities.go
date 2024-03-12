package handlers

import (
	"context"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/business/data/identity"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type identitiesHandler struct {
	pool *qubic.Pool
}

func (h *identitiesHandler) One(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qc, err := h.pool.Get()
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "getting qubic conn from pool"))
	}
	if err != nil {
		return web.RespondError(ctx, w, errors.Wrap(err, "creating qubic conn"))
	}
	params := web.Params(r)
	id, ok := params["identity"]
	if !ok {
		return web.NewRequestError(errors.New("request should have the id of the address in the endpoint"), http.StatusBadRequest)
	}

	res, err := identity.GetIdentity(ctx, qc, id)
	if err != nil {
		qc.Close()
		return web.RespondError(ctx, w, errors.Wrap(err, "getting balance"))
	}

	h.pool.Put(qc)
	return web.Respond(ctx, w, res, http.StatusOK)
}
