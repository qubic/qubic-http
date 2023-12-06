package handlers

import (
	"github.com/0xluk/go-qubic"
	"github.com/gorilla/mux"
	"net/http"
)

func New(client *qubic.Client) http.Handler {
	router := mux.NewRouter()

	ih := identitiesHandler{client: client}
	router.HandleFunc("/identities/{identity}", ih.One).Methods(http.MethodGet)

	th := tickHandler{client: client}
	router.HandleFunc("/tick-info", th.GetTickInfo).Methods(http.MethodGet)
	router.HandleFunc("/tick-transactions/{tick}", th.GetTickTransactions).Methods(http.MethodGet)
	router.HandleFunc("/tick-data/{tick}", th.GetTickData).Methods(http.MethodGet)

	return router
}
