package handlers

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/external/opensearch"
	"github.com/qubic/qubic-http/foundation/web"
	"net/http"
)

type statusHandler struct {
	opensearchClient *opensearch.Client
}

func (h *statusHandler) GetStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status, err := h.opensearchClient.GetStatus(ctx)
	if err != nil {
		return errors.Wrap(err, "getting status from opensearch")
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
