package handlers

import (
	"context"
	"github.com/0xluk/go-qubic/foundation/tcp"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/business/data/identity"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type identitiesHandler struct {
	pool *nodes.Pool
}

func (h *identitiesHandler) One(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ip := h.pool.GetRandomIP()
	qc, err := tcp.NewQubicConnection(ctx, ip, "21841")
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
		return web.RespondError(ctx, w, errors.Wrap(err, "getting balance"))
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}