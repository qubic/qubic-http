package tx

//type GetTxStatusInput struct {
//	Tick      uint32 `json:"tick" validate:"required"`
//	HexDigest string `json:"digest" validate:"min=64,max=64"`
//	Signature string `json:"signature"`
//}
//
//func (i *GetTxStatusInput) toQubicModel() (types.RequestTxStatus, error) {
//	digest, err := hex.DecodeString(i.HexDigest)
//	if err != nil {
//		return tx.RequestTxStatus{}, errors.Wrap(err, "hex decoding digest")
//	}
//	if len(digest) != 32 {
//		return tx.RequestTxStatus{}, errors.Errorf("Hex digest input expected 32 chars. Got: %d", len(digest))
//	}
//
//	var qubicDigest [32]byte
//	copy(qubicDigest[:], digest)
//
//	signature, err := hex.DecodeString(i.Signature)
//	if err != nil {
//		return tx.RequestTxStatus{}, errors.Wrap(err, "hex decoding sig")
//	}
//	//if len(signature) != 64 {
//	//	return tx.RequestTxStatus{}, errors.Errorf("Hex signature input expected 64 chars. Got: %d", len(signature))
//	//}
//
//	var qubicSignature [64]byte
//	copy(qubicSignature[:], signature)
//
//	return tx.RequestTxStatus{
//		Tick:      i.Tick,
//		Digest:    qubicDigest,
//		Signature: qubicSignature,
//	}, nil
//}
//
//type GetTxStatusOutput struct {
//	CurrentTickOfNode uint32 `json:"current_tick_of_node"`
//	Tick              uint32 `json:"tick"`
//	MoneyFlew         bool   `json:"money_flew"`
//	Executed          bool   `json:"executed"`
//	Notfound          bool   `json:"not_found"`
//	HexDigest         string `json:"hex_digest"`
//}
//
//func (o GetTxStatusOutput) fromQubicModel(model tx.ResponseTxStatus) GetTxStatusOutput {
//	return GetTxStatusOutput{
//		CurrentTickOfNode: model.CurrentTickOfNode,
//		Tick:              model.TickOfTx,
//		MoneyFlew:         model.MoneyFlew,
//		Executed:          model.Executed,
//		Notfound:          model.NotFound,
//		HexDigest:         hex.EncodeToString(model.Digest[:]),
//	}
//}

//func (input *SendSignedTxInput) toQubicModel() (types., error) {
//	rawTx, err := hex.DecodeString(input.SignedTx)
//	if err != nil {
//		return tx.SignedTransaction{}, errors.New("failed to decode SignedTx")
//	}
//
//	signatureBytes, err := hex.DecodeString(input.HexTxSignature)
//	if err != nil {
//		return tx.SignedTransaction{}, errors.New("failed to decode HexTxSignature")
//	}
//
//	var signature [64]byte
//	copy(signature[:], signatureBytes)
//
//	return tx.SignedTransaction{
//		SignedTx:     rawTx,
//		Signature: signature,
//	}, nil
//}

type SendSignedTxInput struct {
	SignedTx string `json:"signed_tx" validate:"required"`
}
