package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"qubic-api-sidecar/foundation/web"

	"github.com/0xluk/go-qubic"
)

type identitiesHandler struct {
	client qubic.Client
}

func (h *identitiesHandler) One(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	identity := vars["identity"]

	res, err := h.client.GetBalance(context.Background(), identity)
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "parsing input"), http.StatusInternalServerError)
		return
	}

	web.Respond(w, res)
}