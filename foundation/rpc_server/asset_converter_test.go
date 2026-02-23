package rpc

import (
	"github.com/qubic/go-node-connector/types"
	"github.com/qubic/qubic-http/protobuff"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssetConverter_convertAssetIssuance(t *testing.T) {

	source := types.AssetIssuance{
		Asset: types.AssetIssuanceData{
			PublicKey:             [32]byte{},
			Type:                  1,
			Name:                  [7]int8{82, 65, 78, 68, 79, 77, 0},
			NumberOfDecimalPlaces: 0,
			UnitOfMeasurement:     [7]int8{1, 2, 3, 4, 5, 6, 7},
		},
		Tick:          42,
		UniverseIndex: 43,
	}

	converted, err := convertAssetIssuance(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.AssetIssuance{
		Data: &protobuff.AssetIssuanceData{
			IssuerIdentity:        "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
			Type:                  1,
			Name:                  "RANDOM",
			NumberOfDecimalPlaces: 0,
			UnitOfMeasurement:     []int32{1, 2, 3, 4, 5, 6, 7},
		},
		Tick:          42,
		UniverseIndex: 43,
	}, converted)

}

func TestAssetConverter_convertAssetOwnership(t *testing.T) {

	id := types.Identity("TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI")
	pubKey, err := id.ToPubKey(false)
	assert.NoError(t, err)

	source := types.AssetOwnership{
		Asset: types.AssetOwnershipData{
			PublicKey:             pubKey,
			Type:                  2,
			Padding:               [1]int8{},
			ManagingContractIndex: 1,
			IssuanceIndex:         123,
			NumberOfUnits:         456,
		},
		Tick:          42,
		UniverseIndex: 43,
	}

	converted, err := convertAssetOwnership(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.AssetOwnership{
		Data: &protobuff.AssetOwnershipData{
			OwnerIdentity:         "TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI",
			Type:                  2,
			ManagingContractIndex: 1,
			IssuanceIndex:         123,
			NumberOfUnits:         456,
		},
		Tick:          42,
		UniverseIndex: 43,
	}, converted)

}

func TestAssetConverter_convertAssetPossession(t *testing.T) {

	id := types.Identity("TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI")
	pubKey, err := id.ToPubKey(false)
	assert.NoError(t, err)

	source := types.AssetPossession{
		Asset: types.AssetPossessionData{
			PublicKey:             pubKey,
			Type:                  3,
			Padding:               [1]int8{},
			ManagingContractIndex: 1,
			OwnershipIndex:        123,
			NumberOfUnits:         456,
		},
		Tick:          42,
		UniverseIndex: 43,
	}

	converted, err := convertAssetPossession(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.AssetPossession{
		Data: &protobuff.AssetPossessionData{
			PossessorIdentity:     "TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI",
			Type:                  3,
			ManagingContractIndex: 1,
			OwnershipIndex:        123,
			NumberOfUnits:         456,
		},
		Tick:          42,
		UniverseIndex: 43,
	}, converted)

}

func TestAssetConverter_convertContractIpo(t *testing.T) {

	id := types.Identity("TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI")
	pubKey, err := id.ToPubKey(false)
	assert.NoError(t, err)

	source := types.ContractIpo{
		ContractIndex: 5,
		TickNumber:    100,
	}

	// set one bid with a known identity and price
	source.PubKeys[0] = pubKey
	source.Prices[0] = 1000

	// set another bid at a different index
	source.PubKeys[10] = pubKey
	source.Prices[10] = 2000

	// index 1 has zero price, should be skipped
	source.PubKeys[1] = pubKey
	source.Prices[1] = 0

	converted, err := convertContractIpo(source)
	assert.NoError(t, err)

	assert.Equal(t, uint32(5), converted.ContractIndex)
	assert.Equal(t, uint32(100), converted.TickNumber)

	// only 2 bids should be present (zero-price skipped)
	assert.Len(t, converted.Bids, 2)

	assert.Equal(t, &protobuff.IpoBid{
		Identity: "TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI",
		Amount:   1000,
	}, converted.Bids[0])

	assert.Equal(t, &protobuff.IpoBid{
		Identity: "TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI",
		Amount:   2000,
	}, converted.Bids[10])

	// zero-price bid should not be in the map
	assert.Nil(t, converted.Bids[1])
}
