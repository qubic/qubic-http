package handlers

import (
	_ "github.com/qubic/qubic-http/app/server/docs"
	mid2 "github.com/qubic/qubic-http/business/mid"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"log"
	"net/http"
	"os"
)

func New(shutdown chan os.Signal, log *log.Logger, pool *nodes.Pool) http.Handler {
	app := web.NewApp(shutdown, mid2.Logger(log), mid2.Errors(log), mid2.Metrics(), mid2.Panics(log))

	ih := identitiesHandler{pool: pool}
	app.Handle(http.MethodGet, "/v1/identities/:identity", ih.One)

	th := tickHandler{pool: pool}
	app.Handle(http.MethodGet, "/v1/tick-info", th.GetTickInfo)
	app.Handle(http.MethodGet, "/v1/tick-transactions/:tick", th.GetTickTransactions)
	app.Handle(http.MethodGet, "/v1/tick-data/:tick", th.GetTickData)

	txH := txHandler{pool: pool}
	app.Handle(http.MethodPost, "/v1/send-tx", txH.SendSignedTx)
	app.Handle(http.MethodPost, "/v1/get-tx-status", txH.GetTxStatus)

	return app
}
