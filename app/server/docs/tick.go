package docs

import "github.com/qubic/qubic-http/business/data/tick"

// swagger-ui:route GET /v1/tick-info tick tick-info
// Get tick info.
// responses:
//   200: tickInfoResponse

// Tick Info
// swagger-ui:response tickInfoResponse
type tickInfoResponse struct {
	//in:body
	Body tick.GetTickInfoOutput
}

// swagger-ui:route GET /v1/tick-data/{tick} tick get-tick-data
// Get tick data.
// responses:
//   200: tickDataResponse

// Tick Data
// swagger-ui:response tickDataResponse
type tickDataResponse struct {
	//in:body
	Body tick.GetTickDataOutput
}

// swagger-ui:parameters get-tick-data
type tickDataRequest struct {
	// Tick number
	// in:path
	Tick string `json:"tick"`
}


// swagger-ui:route GET /v1/tick-transactions/{tick} tick get-tick-transactions
// Get tick transactions.
// responses:
//   200: tickTransactionsResponse
// Tick Transactions
// swagger-ui:response tickTransactionsResponse
type tickTransactionsResponse struct {
	//in:body
	Body tick.GetTickTransactionsOutput
}

// swagger-ui:parameters get-tick-transactions
type tickTransactionsRequest struct {
	// Tick number
	// in:path
	Tick string `json:"tick"`
}