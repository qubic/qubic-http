package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func Decode(r *http.Request, val interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "reading request body")
	}

	err = json.Unmarshal(b, val)
	if err != nil {
		return errors.Wrap(err, "unmarshalling request body data")
	}

	return nil
}
