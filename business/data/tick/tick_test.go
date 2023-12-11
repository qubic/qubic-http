package tick

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/0xluk/go-qubic/data/tick"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func generateRandomByteArray(size int) [][32]byte {
	var result [][32]byte
	for i := 0; i < size; i++ {
		var byteArray [32]byte
		_, err := rand.Read(byteArray[:])
		if err != nil {
			panic(err)
		}
		result = append(result, byteArray)
	}
	return result
}

func TestGetTickDataOutput_FromQubicModel(t *testing.T) {
	var union [256]byte
	_, err := rand.Read(union[:])
	if err != nil {
		t.Fatalf("Got err when generating random union. err :%s", err.Error())
	}

	var txDigests [1024][32]byte
	copy(txDigests[:], generateRandomByteArray(tick.NUMBER_OF_TRANSACTIONS_PER_TICK))

	var sig [64]byte
	_, err = rand.Read(sig[:])
	if err != nil {
		t.Fatalf("Got err when generating random sig. err :%s", err.Error())
	}

	// Sample Qubic model data with random byte arrays
	qubicModel := tick.TickData{
		ComputorIndex:      42,
		Epoch:              2023,
		Tick:               123456,
		Millisecond:        789,
		Second:             45,
		Minute:             30,
		Hour:               12,
		Day:                15,
		Month:              11,
		Year:               22,
		UnionData:          union,
		Timelock:           generateRandomByteArray(32)[0],
		TransactionDigests: txDigests,
		ContractFees:       [1024]int64{100, 200, 300},
		Signature:          sig,
	}

	// Expected output
	expectedOutput := GetTickDataOutput{
		ComputorIndex:      42,
		Epoch:              2023,
		Tick:               123456,
		Millisecond:        789,
		Second:             45,
		Minute:             30,
		Hour:               12,
		Day:                15,
		Month:              11,
		Year:               22,
		HexUnionData:       hex.EncodeToString(union[:]),
		HexTimelock:        hex.EncodeToString(qubicModel.Timelock[:]),
		TransactionDigests: byteArraysToHexStrings(qubicModel.TransactionDigests[:]),
		ContractFees:       []int64{100, 200, 300},
		Signature:          hex.EncodeToString(qubicModel.Signature[:]),
	}

	// Convert Qubic model to GetTickDataOutput
	actualOutput := new(GetTickDataOutput).fromQubicModel(qubicModel)

	// Compare actual and expected outputs
	if diff := cmp.Diff(actualOutput, expectedOutput); diff != "" {
		t.Errorf("Mismatch (-actual +expected):\n%s", diff)
	}
}

func TestGetTickTransactionsOutput_FromQubicModel(t *testing.T) {
	var sourcePubKey [32]byte
	_, err := rand.Read(sourcePubKey[:])
	if err != nil {
		t.Fatalf("Got err when generating random source pubkey. err :%s", err.Error())
	}

	var destPubKey [32]byte
	_, err = rand.Read(destPubKey[:])
	if err != nil {
		t.Fatalf("Got err when generating random dest pubkey. err :%s", err.Error())
	}

	qubicModel := []tick.Transaction{
		{
			SourcePublicKey:      sourcePubKey,
			DestinationPublicKey: destPubKey,
			Amount:               15,
			Tick:                 10,
			InputType:            7,
			InputSize:            9,
		},
		{
			SourcePublicKey:      destPubKey,
			DestinationPublicKey: sourcePubKey,
			Amount:               6,
			Tick:                 1,
			InputType:            9,
			InputSize:            7,
		},
	}

	expected := GetTickTransactionsOutput{
		{
			SourcePublicKey:      hex.EncodeToString(sourcePubKey[:]),
			DestinationPublicKey: hex.EncodeToString(destPubKey[:]),
			Amount:               15,
			Tick:                 10,
			InputType:            7,
			InputSize:            9,
		},
		{
			SourcePublicKey:      hex.EncodeToString(destPubKey[:]),
			DestinationPublicKey: hex.EncodeToString(sourcePubKey[:]),
			Amount:               6,
			Tick:                 1,
			InputType:            9,
			InputSize:            7,
		},
	}

	var got GetTickTransactionsOutput
	got = got.fromQubicModel(qubicModel)
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("Mismatch (-actual +expected):\n%s", diff)
	}
}
