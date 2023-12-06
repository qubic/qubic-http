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
	client *qubic.Client
}

func (h *identitiesHandler) One(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	identity := vars["identity"]

	res, err := qubic.GetBalance(context.Background(), h.client.Qc, identity)
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "getting balance"), http.StatusInternalServerError)
		return
	}

	web.Respond(w, res)
}