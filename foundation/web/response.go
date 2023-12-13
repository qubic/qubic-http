package web

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	ctx, span := otel.Tracer("qubic-http").Start(ctx, "foundation.web.respond")
	defer span.End()

	// Set the status code for the request logger middleware.
	// If the context is missing this value, request the service
	// to be shutdown gracefully.
	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return NewShutdownError("web value missing from context")
	}
	v.StatusCode = statusCode

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	var res interface{}
	var code int

	switch v := errors.Cause(err).(type) {
	case *FieldsError:
		res = FieldErrorResponse{
			Error:  v.Err.Error(),
			Fields: v.Fields,
		}
		code = http.StatusBadRequest
	case *RequestError:
		res = ErrorResponse{Error: v.Err.Error()}
		code = v.Status
	default:
		res = ErrorResponse{Error: err.Error()}
		code = http.StatusInternalServerError
	}

	if err := Respond(ctx, w, res, code); err != nil {
		return err
	}
	return nil
}
