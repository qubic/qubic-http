package tx

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/0xluk/go-qubic/data/tx"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetTxStatusInput_ToQubicModel(t *testing.T) {
	var expectedDigest [32]byte
	_, err := rand.Read(expectedDigest[:])
	if err != nil {
		t.Fatalf("Got err when generating random digest. err: %s", err.Error())
	}
	
	var expectedSignature [64]byte
	_, err = rand.Read(expectedSignature[:])
	if err != nil {
		t.Fatalf("Got err when generating random signature. err: %s", err.Error())
	}
	expected := tx.RequestTxStatus{
		Tick:      1,
		Digest:    expectedDigest,
		Signature: expectedSignature,
	}
	
	input := GetTxStatusInput{
		Tick:      1,
		HexDigest: hex.EncodeToString(expectedDigest[:]),
		Signature: hex.EncodeToString(expectedSignature[:]),
	}

	got, err := input.toQubicModel()
	if err != nil {
		t.Fatalf("Got err when converting input to qubic. err :%s", err.Error())
	}

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatalf("expected is different from got. diff: %s", diff)
	}
}

func TestTxStatusOutput_FromQubicModel(t *testing.T) {
	var digest [32]byte
	_, err := rand.Read(digest[:])
	if err != nil {
		t.Fatalf("Got err when generating random digest. err :%s", err.Error())
	}
	var padding [2]byte
	_, err = rand.Read(padding[:])
	if err != nil {
		t.Fatalf("Got err when generating padding. err: %s", err.Error())
	}

	expected := GetTxStatusOutput{
		CurrentTickOfNode: 15,
		Tick:              20,
		MoneyFlew:         false,
		Executed:          true,
		HexPadding:        hex.EncodeToString(padding[:]),
		HexDigest:         hex.EncodeToString(digest[:]),
	}

	model := tx.ResponseTxStatus{
		CurrentTickOfNode: 15,
		TickOfTx:          20,
		MoneyFlew:         false,
		Executed:          true,
		Padding:           padding,
		Digest:            digest,
	}

	var got GetTxStatusOutput

	got = got.fromQubicModel(model)
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatalf("expected is different from got. diff: %s", diff)
	}
}

func TestSendSignedTxInput_ToQubicModel(t *testing.T) {
	input := SendSignedTxInput{
		HexRawTx:       "68656c6c6f20776f726c64", // "hello world"
		HexTxSignature: "746573745369676e6174757265", // "testSignature"
	}

	expected := tx.SignedTransaction{
		RawTx:     []byte("hello world"),
		Signature: [64]byte{116, 101, 115, 116, 83, 105, 103, 110, 97, 116, 117, 114, 101},
	}

	result, err := input.toQubicModel()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	diff := cmp.Diff(expected, result)
	if diff != "" {
		t.Errorf("Mismatch (-expected +got):\n%s", diff)
	}
}
