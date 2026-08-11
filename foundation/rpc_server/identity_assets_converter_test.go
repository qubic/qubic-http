package rpc

import (
	"testing"

	qubic "github.com/qubic/go-node-connector/v2"
	"github.com/qubic/go-node-connector/v2/types"
	"github.com/qubic/qubic-http/protobuff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	issuerIdentity    = "TESTIOGXQKYYZEQXOXFSWWAJNYLCDBWFAPNBLNBUZFHDVFMYPJZXGMEEJEGI"
	ownerIdentity     = "AFZPUAIYVPNUYGJRQVLUKOPPVLHAZQTGLYAAUUNBXFTVTAMSBKQBLEIEPCVJ"
	possessorIdentity = "TLIDHAWGQTCELDIPZFHAKTNFFWMEAZKZALRUFSJUDALPRMVBHXCQNVYEEBDH"
)

func pubKeyOf(t *testing.T, identity string) [32]byte {
	t.Helper()

	id := types.Identity(identity)
	pubKey, err := id.ToPubKey(false)
	require.NoError(t, err)

	return pubKey
}

func testIssuedAssetData(t *testing.T) types.IssuedAssetData {
	t.Helper()

	return types.IssuedAssetData{
		PublicKey:             pubKeyOf(t, issuerIdentity),
		Type:                  1,
		Name:                  [7]int8{82, 65, 78, 68, 79, 77, 0},
		NumberOfDecimalPlaces: 2,
		UnitOfMeasurement:     [7]int8{1, 2, 3, 4, 5, 6, 7},
	}
}

func expectedIssuedAssetData() *protobuff.IssuedAssetData {
	return &protobuff.IssuedAssetData{
		IssuerIdentity:        issuerIdentity,
		Type:                  1,
		Name:                  "RANDOM",
		NumberOfDecimalPlaces: 2,
		UnitOfMeasurement:     []int32{1, 2, 3, 4, 5, 6, 7},
	}
}

// testOwnedAssetData uses values that differ from the possession record wrapping it in
// TestAssetConverter_convertPossessedAsset, so that mixing the two up is detectable.
func testOwnedAssetData(t *testing.T) types.OwnedAssetData {
	t.Helper()

	return types.OwnedAssetData{
		PublicKey:             pubKeyOf(t, ownerIdentity),
		Type:                  2,
		Padding:               [1]int8{7},
		ManagingContractIndex: 1,
		IssuanceIndex:         99,
		NumberOfUnits:         42,
		IssuedAsset:           testIssuedAssetData(t),
	}
}

func expectedOwnedAssetData() *protobuff.OwnedAssetData {
	return &protobuff.OwnedAssetData{
		OwnerIdentity:         ownerIdentity,
		Type:                  2,
		Padding:               7,
		ManagingContractIndex: 1,
		IssuanceIndex:         99,
		NumberOfUnits:         42,
		IssuedAsset:           expectedIssuedAssetData(),
	}
}

func TestAssetConverter_convertIssuedAsset(t *testing.T) {
	source := types.IssuedAsset{
		Data: testIssuedAssetData(t),
		Info: types.AssetInfo{Tick: 42, UniverseIndex: 43},
	}

	converted, err := convertIssuedAsset(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.IssuedAsset{
		Data: expectedIssuedAssetData(),
		Info: &protobuff.AssetInfo{Tick: 42, UniverseIndex: 43},
	}, converted)
}

func TestAssetConverter_convertOwnedAsset(t *testing.T) {
	source := types.OwnedAsset{
		Data: testOwnedAssetData(t),
		Info: types.AssetInfo{Tick: 42, UniverseIndex: 43},
	}

	converted, err := convertOwnedAsset(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.OwnedAsset{
		Data: expectedOwnedAssetData(),
		Info: &protobuff.AssetInfo{Tick: 42, UniverseIndex: 43},
	}, converted)
}

// TestAssetConverter_convertPossessedAsset pins that the nested ownership record is
// converted from its own data. Every field of the possession record deliberately differs
// from the ownership record it wraps, so copying one into the other would fail here.
func TestAssetConverter_convertPossessedAsset(t *testing.T) {
	source := types.PossessedAsset{
		Data: types.PossessedAssetData{
			PublicKey:             pubKeyOf(t, possessorIdentity),
			Type:                  3,
			Padding:               [1]int8{5},
			ManagingContractIndex: 8,
			IssuanceIndex:         11,
			NumberOfUnits:         7,
			OwnedAsset:            testOwnedAssetData(t),
		},
		Info: types.AssetInfo{Tick: 42, UniverseIndex: 43},
	}

	converted, err := convertPossessedAsset(source)
	assert.NoError(t, err)

	assert.Equal(t, &protobuff.PossessedAsset{
		Data: &protobuff.PossessedAssetData{
			PossessorIdentity:     possessorIdentity,
			Type:                  3,
			Padding:               5,
			ManagingContractIndex: 8,
			IssuanceIndex:         11,
			NumberOfUnits:         7,
			OwnedAsset:            expectedOwnedAssetData(),
		},
		Info: &protobuff.AssetInfo{Tick: 42, UniverseIndex: 43},
	}, converted)
}

func TestAssetConverter_convertOwnedAssets_empty(t *testing.T) {
	converted, err := convertOwnedAssets(types.OwnedAssets{})
	assert.NoError(t, err)

	assert.NotNil(t, converted)
	assert.Empty(t, converted)
}

func TestAssetConverter_convertPossessedAssets_empty(t *testing.T) {
	converted, err := convertPossessedAssets(types.PossessedAssets{})
	assert.NoError(t, err)

	assert.NotNil(t, converted)
	assert.Empty(t, converted)
}

func TestAssetConverter_convertOwnedAssets_preservesOrder(t *testing.T) {
	first := testOwnedAssetData(t)
	first.NumberOfUnits = 1
	second := testOwnedAssetData(t)
	second.NumberOfUnits = 2

	converted, err := convertOwnedAssets(types.OwnedAssets{
		{Data: first, Info: types.AssetInfo{Tick: 1}},
		{Data: second, Info: types.AssetInfo{Tick: 2}},
	})
	assert.NoError(t, err)

	require.Len(t, converted, 2)
	assert.Equal(t, int64(1), converted[0].Data.NumberOfUnits)
	assert.Equal(t, int64(2), converted[1].Data.NumberOfUnits)
}

func TestAssetConverter_convertAddressAssets(t *testing.T) {
	source := []qubic.AddressAssets{
		{
			Identity: ownerIdentity,
			Owned: types.OwnedAssets{
				{Data: testOwnedAssetData(t), Info: types.AssetInfo{Tick: 42, UniverseIndex: 43}},
			},
			Possessed: types.PossessedAssets{},
		},
		{
			Identity:  possessorIdentity,
			Owned:     types.OwnedAssets{},
			Possessed: types.PossessedAssets{},
		},
	}

	converted, err := convertAddressAssets(source)
	assert.NoError(t, err)

	require.Len(t, converted, 2)

	assert.Equal(t, ownerIdentity, converted[0].Identity)
	require.Len(t, converted[0].OwnedAssets, 1)
	assert.Equal(t, expectedOwnedAssetData(), converted[0].OwnedAssets[0].Data)
	assert.Empty(t, converted[0].PossessedAssets)

	// an identity holding nothing still gets an entry, with empty lists rather than
	// being dropped from the response
	assert.Equal(t, possessorIdentity, converted[1].Identity)
	assert.NotNil(t, converted[1].OwnedAssets)
	assert.Empty(t, converted[1].OwnedAssets)
	assert.NotNil(t, converted[1].PossessedAssets)
	assert.Empty(t, converted[1].PossessedAssets)
}

func TestAssetConverter_convertAddressAssets_empty(t *testing.T) {
	converted, err := convertAddressAssets(nil)
	assert.NoError(t, err)

	assert.NotNil(t, converted)
	assert.Empty(t, converted)
}
