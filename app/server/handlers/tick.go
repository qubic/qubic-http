package handlers

import (
	"context"
	"github.com/0xluk/go-qubic"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"qubic-api-sidecar/foundation/web"
	"strconv"
)

type tickHandler struct {
	client *qubic.Client
}

func (h *tickHandler) GetTickInfo(w http.ResponseWriter, r *http.Request) {
	res, err := qubic.GetTickInfo(context.Background(), h.client.Qc)
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "getting tick info"), http.StatusInternalServerError)
		return
	}

	web.Respond(w, res)
}

func (h *tickHandler) GetTickTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tick := vars["tick"]
	tickNumber, err := strconv.ParseInt(tick, 10, 32)
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "parsing input"), http.StatusInternalServerError)
	}

	res, err := qubic.GetTickTransactions(context.Background(), h.client.Qc, uint32(tickNumber))
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "getting tick transactions"), http.StatusInternalServerError)
		return
	}

	web.Respond(w, res)
}

func (h *tickHandler) GetTickData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tick := vars["tick"]
	tickNumber, err := strconv.ParseInt(tick, 10, 32)
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "parsing input"), http.StatusInternalServerError)
	}

	res, err := qubic.GetTickData(context.Background(), h.client.Qc, uint32(tickNumber))
	if err != nil {
		web.RespondError(w, errors.Wrap(err, "getting tick data"), http.StatusInternalServerError)
		return
	}

	web.Respond(w, res)
}
