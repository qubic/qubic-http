package identity

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/0xluk/go-qubic/data/identity"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetIdentityOutput_FromQubicModel(t *testing.T) {
	var pubkey [32]byte
	_, err := rand.Read(pubkey[:])
	if err != nil {
		t.Fatalf("Got err when generating random pubkey. err :%s", err.Error())
	}
	// Example Qubic model response
	qubicModel := identity.GetIdentityResponse{
		Entity: identity.Entity{
			PublicKey:                  pubkey, // replace with your public key
			IncomingAmount:             100,
			OutgoingAmount:             50,
			NumberOfIncomingTransfers:  5,
			NumberOfOutgoingTransfers:  3,
			LatestIncomingTransferTick: 123,
			LatestOutgoingTransferTick: 456,
		},
		Siblings: [identity.SpectrumDepth][32]byte{
			{1, 2, 3}, // replace with your sibling data
			{4, 5, 6}, // replace with your sibling data
			// Add more siblings as needed
		},
	}

	// Expected output after conversion
	expectedOutput := GetIdentityOutput{
		PublicKey:                  hex.EncodeToString(qubicModel.Entity.PublicKey[:]),
		IncomingAmount:             qubicModel.Entity.IncomingAmount,
		OutgoingAmount:             qubicModel.Entity.OutgoingAmount,
		NumberOfIncomingTransfers:  qubicModel.Entity.NumberOfIncomingTransfers,
		NumberOfOutgoingTransfers:  qubicModel.Entity.NumberOfOutgoingTransfers,
		LatestIncomingTransferTick: qubicModel.Entity.LatestIncomingTransferTick,
		LatestOutgoingTransferTick: qubicModel.Entity.LatestOutgoingTransferTick,
		Siblings: []string{
			hex.EncodeToString(qubicModel.Siblings[0][:]),
			hex.EncodeToString(qubicModel.Siblings[1][:]),
			// Add more siblings as needed
		},
	}

	// Perform the conversion
	var got GetIdentityOutput
	got = got.fromQubicModel(qubicModel)

	// Use cmp.Diff to compare the expected and actual outputs
	diff := cmp.Diff(expectedOutput, got)
	if diff != "" {
		t.Errorf("Mismatch (-expected +got):\n%s", diff)
	}
}
