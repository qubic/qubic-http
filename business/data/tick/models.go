package tick

import (
	"encoding/hex"
	"github.com/0xluk/go-qubic/data/tick"
	"github.com/qubic/qubic-http/external/opensearch"
)

type GetTickDataOutput struct {
	ComputorIndex uint16 `json:"computor_index"`
	Epoch         uint16 `json:"epoch"`
	Tick          uint32 `json:"tick"`
	Millisecond   uint16 `json:"millisecond"`
	Second        uint8  `json:"second"`
	Minute        uint8  `json:"minute"`
	Hour          uint8  `json:"hour"`
	Day           uint8  `json:"day"`
	Month         uint8  `json:"month"`
	Year          uint8  `json:"year"`
	//HexUnionData       string   `json:"hex_union_data"`
	HexTimelock        string   `json:"hex_timelock"`
	TransactionDigests []string `json:"transaction_digests"`
	//ContractFees       []int64  `json:"contract_fees"`
	Signature    string                   `json:"signature"`
	PotentialBxs []opensearch.PotentialBx `json:"potentialBx"`
}

type PotentialBx struct {
	Index       int    `json:"index"`
	Destination string `json:"dest"`
	Amount      string `json:"amount"`
}

func (o *GetTickDataOutput) fromQubicModel(model tick.TickData) GetTickDataOutput {
	var contractFees []int64
	for _, b := range model.ContractFees {
		if b != 0 {
			contractFees = append(contractFees, b)
		}
	}

	return GetTickDataOutput{
		ComputorIndex: model.ComputorIndex,
		Epoch:         model.Epoch,
		Tick:          model.Tick,
		Millisecond:   model.Millisecond,
		Second:        model.Second,
		Minute:        model.Minute,
		Hour:          model.Hour,
		Day:           model.Day,
		Month:         model.Month,
		Year:          model.Year,
		//HexUnionData:       hex.EncodeToString(model.UnionData[:]),
		HexTimelock:        hex.EncodeToString(model.Timelock[:]),
		TransactionDigests: byteArraysToHexStrings(model.TransactionDigests[:]),
		//ContractFees:       contractFees,
		Signature: hex.EncodeToString(model.Signature[:]),
	}
}

func (o *GetTickDataOutput) fromOpensearchModel(model opensearch.TickDataResponse) GetTickDataOutput {
	timeArr := timeSliceTArr(model.Time)
	return GetTickDataOutput{
		ComputorIndex:      uint16(model.Computor),
		Epoch:              uint16(model.Epoch),
		Tick:               model.Tick,
		Millisecond:        uint16(timeArr[0]),
		Second:             uint8(timeArr[1]),
		Minute:             uint8(timeArr[2]),
		Hour:               uint8(timeArr[3]),
		Day:                uint8(timeArr[4]),
		Month:              uint8(timeArr[5]),
		Year:               uint8(timeArr[6]),
		HexTimelock:        model.Timelock,
		TransactionDigests: model.TransactionIDs,
		Signature:          model.Signature,
		PotentialBxs:       model.PotentialBxs,
	}
}

func byteArraysToHexStrings(arrays [][32]byte) []string {
	var zeroArray [32]byte
	result := make([]string, 0, len(arrays))
	for _, arr := range arrays {
		if arr == zeroArray {
			continue
		}

		result = append(result, hex.EncodeToString(arr[:]))
	}
	return result
}

type Transaction struct {
	SourcePublicKey      string `json:"source_public_key"`
	DestinationPublicKey string `json:"destination_public_key"`
	Amount               int64  `json:"amount"`
	Tick                 uint32 `json:"tick"`
	InputType            uint16 `json:"input_type"`
	InputSize            uint16 `json:"input_size"`
	Hash                 string `json:"hash"`
}

type GetTickTransactionsOutput []Transaction

func (o *GetTickTransactionsOutput) fromQubicModel(tickTxs []tick.Transaction) GetTickTransactionsOutput {
	txs := make([]Transaction, 0, len(tickTxs))

	for _, tx := range tickTxs {
		txs = append(txs, Transaction{
			SourcePublicKey:      hex.EncodeToString(tx.Data.Header.SourcePublicKey[:]),
			DestinationPublicKey: hex.EncodeToString(tx.Data.Header.DestinationPublicKey[:]),
			Amount:               tx.Data.Header.Amount,
			Tick:                 tx.Data.Header.Tick,
			InputType:            tx.Data.Header.InputType,
			InputSize:            tx.Data.Header.InputSize,
			Hash:                 byteArrayToString(tx.Hash),
		})
	}

	return txs
}

type GetTickInfoOutput struct {
	TickDuration            uint16 `json:"tick_duration"`
	Epoch                   uint16 `json:"epoch"`
	Tick                    uint32 `json:"tick"`
	NumberOfAlignedVotes    uint16 `json:"number_of_aligned_votes"`
	NumberOfMisalignedVotes uint16 `json:"number_of_misaligned_votes"`
}

func (o *GetTickInfoOutput) fromQubicModel(model tick.CurrentTickInfo) GetTickInfoOutput {
	return GetTickInfoOutput{
		TickDuration:            model.TickDuration,
		Epoch:                   model.Epoch,
		Tick:                    model.Tick,
		NumberOfAlignedVotes:    model.NumberOfAlignedVotes,
		NumberOfMisalignedVotes: model.NumberOfMisalignedVotes,
	}
}

func byteArrayToString(arr [60]byte) string {
	var zeroArray [60]byte

	if arr == zeroArray {
		return ""
	}

	return string(arr[:])
}

func timeSliceTArr(slice []int) [7]int {
	var arr [7]int

	for i, t := range slice {
		arr[i] = t
	}

	return arr
}
