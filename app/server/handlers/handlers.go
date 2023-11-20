package handlers

import (
	"github.com/0xluk/go-qubic"
	"github.com/gorilla/mux"
	"net/http"
)

func New(client qubic.Client) http.Handler {
	router := mux.NewRouter()

	ih := identitiesHandler{client: client}
	router.HandleFunc("/identities/{identity}", ih.One).Methods(http.MethodGet)

	return router
}
