package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func RespondError(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err.Error())
	w.WriteHeader(statusCode)
}

func Respond(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		RespondError(w, errors.Wrap(err, "parsing output"), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
