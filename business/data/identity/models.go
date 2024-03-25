package identity

import (
	"encoding/hex"
	"github.com/qubic/go-node-connector/types"
)

type GetIdentityOutput struct {
	PublicKey                  string   `json:"public_key"`
	Tick                       uint32   `json:"tick"`
	Balance                    int64    `json:"balance"` // incoming_amount - outgoing_amount
	IncomingAmount             int64    `json:"incoming_amount"`
	OutgoingAmount             int64    `json:"outgoing_amount"`
	NumberOfIncomingTransfers  uint32   `json:"number_of_incoming_transfers"`
	NumberOfOutgoingTransfers  uint32   `json:"number_of_outgoing_transfers"`
	LatestIncomingTransferTick uint32   `json:"latest_incoming_transfer_tick"`
	LatestOutgoingTransferTick uint32   `json:"latest_outgoing_transfer_tick"`
	Siblings                   []string `json:"siblings"`
}

func (o *GetIdentityOutput) fromQubicModel(model types.AddressInfo) GetIdentityOutput {
	return GetIdentityOutput{
		PublicKey:                  hex.EncodeToString(model.AddressData.PublicKey[:]),
		Tick:                       model.Tick,
		Balance:                    model.AddressData.IncomingAmount - model.AddressData.OutgoingAmount,
		IncomingAmount:             model.AddressData.IncomingAmount,
		OutgoingAmount:             model.AddressData.OutgoingAmount,
		NumberOfIncomingTransfers:  model.AddressData.NumberOfIncomingTransfers,
		NumberOfOutgoingTransfers:  model.AddressData.NumberOfOutgoingTransfers,
		LatestIncomingTransferTick: model.AddressData.LatestIncomingTransferTick,
		LatestOutgoingTransferTick: model.AddressData.LatestOutgoingTransferTick,
		Siblings:                   byteSlicesToStrings(model.Siblings),
	}
}

func byteSlicesToStrings(slices [types.SpectrumDepth][32]byte) []string {
	var zeroArray [32]byte

	result := make([]string, 0, types.SpectrumDepth)
	for _, slice := range slices {
		if slice == zeroArray {
			continue
		}
		result = append(result, hex.EncodeToString(slice[:]))
	}
	return result
}
