package handlers

import (
	"context"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/business/data/identity"
	"github.com/qubic/qubic-http/foundation/web"
	"log"
	"net/http"
)

type identitiesHandler struct {
	pool *qubic.Pool
}

func (h *identitiesHandler) One(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var err error
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
