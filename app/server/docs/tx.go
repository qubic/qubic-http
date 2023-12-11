package docs

import (
	"github.com/qubic/qubic-http/business/data/tick"
	"github.com/qubic/qubic-http/business/data/tx"
)

// swagger-ui:route GET /v1/send-tx tx send-tx
// Send signed transaction.
// responses:
//   201: sendSignedTransactionResponse
// Send signed transaction
// swagger-ui:response sendSignedTransactionResponse
type sendSignedTransactionResponse struct {
	//in:body
	Body struct{}
}

// swagger-ui:parameters send-tx
type sendTxRequest struct {
	// in:body
	Body tx.SendSignedTxInput
}

// swagger-ui:route POST /v1/get-tx-status tx get-tx-status
// Get transaction status
// responses:
//   200: getTransactionStatusResponse
// Get transaction status
// swagger-ui:response getTransactionStatusResponse
type getTransactionStatusResponse struct {
	//in:body
	Body tick.GetTickDataOutput
}

// swagger-ui:parameters get-tx-status
type getTransactionStatusRequest struct {
	// in:body
	Body tx.GetTxStatusInput
}
