package handlers

import (
	_ "github.com/qubic/qubic-http/app/server/docs"
	 "github.com/qubic/qubic-http/business/mid"
	"github.com/qubic/qubic-http/external/opensearch"
	"github.com/qubic/qubic-http/foundation/nodes"
	"github.com/qubic/qubic-http/foundation/web"
	"log"
	"net/http"
	"os"
)

func New(shutdown chan os.Signal, log *log.Logger, pool *nodes.Pool, osclient *opensearch.Client) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	ih := identitiesHandler{pool: pool}
	app.Handle(http.MethodGet, "/v1/identities/:identity", ih.One)

	th := tickHandler{pool: pool, opensearchClient: osclient}
	app.Handle(http.MethodGet, "/v1/tick-info", th.GetTickInfo)
	app.Handle(http.MethodGet, "/v1/tick-transactions/:tick", th.GetTickTransactions)
	//app.Handle(http.MethodGet, "/v1/tick-data/:tick", th.GetTickData)

	app.Handle(http.MethodGet, "/v1/tick-data/:tick", th.GetTickDataV2)

	txH := txHandler{pool: pool, opensearchClient: osclient}
	app.Handle(http.MethodPost, "/v1/send-raw-tx", txH.SendRawTx)
	app.Handle(http.MethodPost, "/v1/get-tx-status", txH.GetTxStatus)
	app.Handle(http.MethodGet, "/v1/tx/:txID", txH.GetTx)
	app.Handle(http.MethodGet, "/v1/bx/:bxID", txH.GetBx)


	return app
}
