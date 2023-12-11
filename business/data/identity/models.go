package identity

import (
	"encoding/hex"
	"github.com/0xluk/go-qubic/data/identity"
)

type GetIdentityOutput struct {
	PublicKey                  string   `json:"public_key"`
	IncomingAmount             int64    `json:"incoming_amount"`
	OutgoingAmount             int64    `json:"outgoing_amount"`
	NumberOfIncomingTransfers  uint32   `json:"number_of_incoming_transfers"`
	NumberOfOutgoingTransfers  uint32   `json:"number_of_outgoing_transfers"`
	LatestIncomingTransferTick uint32   `json:"latest_incoming_transfer_tick"`
	LatestOutgoingTransferTick uint32   `json:"latest_outgoing_transfer_tick"`
	Siblings                   []string `json:"siblings"`
}

func (o *GetIdentityOutput) fromQubicModel(model identity.GetIdentityResponse) GetIdentityOutput {
	return GetIdentityOutput{
		PublicKey:                  hex.EncodeToString(model.Entity.PublicKey[:]),
		IncomingAmount:             model.Entity.IncomingAmount,
		OutgoingAmount:             model.Entity.OutgoingAmount,
		NumberOfIncomingTransfers:  model.Entity.NumberOfIncomingTransfers,
		NumberOfOutgoingTransfers:  model.Entity.NumberOfOutgoingTransfers,
		LatestIncomingTransferTick: model.Entity.LatestIncomingTransferTick,
		LatestOutgoingTransferTick: model.Entity.LatestOutgoingTransferTick,
		Siblings:                   byteSlicesToStrings(model.Siblings),
	}
}

func byteSlicesToStrings(slices [identity.SpectrumDepth][32]byte) []string {
	var zeroArray [32]byte

	result := make([]string, 0, identity.SpectrumDepth)
	for _, slice := range slices {
		if slice == zeroArray {
			continue
		}
		result = append(result, hex.EncodeToString(slice[:]))
	}
	return result
}
