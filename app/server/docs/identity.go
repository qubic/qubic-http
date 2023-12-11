package docs

import (
	"github.com/qubic/qubic-http/business/data/identity"
)

// swagger-ui:route GET /v1/identity/{id} identity get-identity
// Get identity.
// responses:
//   200: identityResponse

// Identity
// swagger-ui:response identityResponse
type identityResponse struct {
	//in:body
	Body identity.GetIdentityOutput
}

// swagger-ui:parameters get-identity
type identityRequest struct {
	// 60 char ID
	// in:path
	ID string `json:"id"`
}
