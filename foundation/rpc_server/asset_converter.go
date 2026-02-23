package rpc

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/qubic/go-node-connector/types"
	"github.com/qubic/qubic-http/protobuff"
)

func convertAssetIssuance(source types.AssetIssuance) (*protobuff.AssetIssuance, error) {
	var issuerIdentity types.Identity
	issuerIdentity, err := issuerIdentity.FromPubKey(source.Asset.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for public key")
	}

	issuedAsset := protobuff.AssetIssuanceData{
		IssuerIdentity:        issuerIdentity.String(),
		Type:                  uint32(source.Asset.Type),
		Name:                  int8ArrayToString(source.Asset.Name[:]),
		NumberOfDecimalPlaces: int32(source.Asset.NumberOfDecimalPlaces),
		UnitOfMeasurement:     int8ArrayToInt32Array(source.Asset.UnitOfMeasurement[:]),
	}

	asset := protobuff.AssetIssuance{
		Data:          &issuedAsset,
		Tick:          source.Tick,
		UniverseIndex: source.UniverseIndex,
	}

	return &asset, nil
}

func convertAssetOwnership(source types.AssetOwnership) (*protobuff.AssetOwnership, error) {
	var owner types.Identity
	owner, err := owner.FromPubKey(source.Asset.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for public key")
	}

	ownedAsset := protobuff.AssetOwnershipData{
		OwnerIdentity:         owner.String(),
		Type:                  uint32(source.Asset.Type),
		ManagingContractIndex: uint32(source.Asset.ManagingContractIndex),
		IssuanceIndex:         source.Asset.IssuanceIndex,
		NumberOfUnits:         source.Asset.NumberOfUnits,
	}

	assetOwnership := protobuff.AssetOwnership{
		Data:          &ownedAsset,
		Tick:          source.Tick,
		UniverseIndex: source.UniverseIndex,
	}

	return &assetOwnership, nil
}

func convertAssetPossession(source types.AssetPossession) (*protobuff.AssetPossession, error) {

	var possessor types.Identity
	possessor, err := possessor.FromPubKey(source.Asset.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for public key")
	}

	possessedAsset := protobuff.AssetPossessionData{
		PossessorIdentity:     possessor.String(),
		Type:                  uint32(source.Asset.Type),
		ManagingContractIndex: uint32(source.Asset.ManagingContractIndex),
		OwnershipIndex:        source.Asset.OwnershipIndex,
		NumberOfUnits:         source.Asset.NumberOfUnits,
	}

	assetPossession := protobuff.AssetPossession{
		Data:          &possessedAsset,
		Tick:          source.Tick,
		UniverseIndex: source.UniverseIndex,
	}

	return &assetPossession, nil
}

func convertContractIpo(source types.ContractIpo) (*protobuff.IpoBidData, error) {

	ipoBidData := protobuff.IpoBidData{
		ContractIndex: source.ContractIndex,
		TickNumber:    source.TickNumber,
		Bids:          make(map[int32]*protobuff.IpoBid),
	}

	for index := 0; index < types.NumberOfComputors; index++ {
		if source.Prices[index] == 0 {
			continue
		}

		identity, err := new(types.Identity).FromPubKey(source.PubKeys[index], false)
		if err != nil {
			return nil, fmt.Errorf("failed to get identity for bid: %w", err)
		}

		ipoBidData.Bids[int32(index)] = &protobuff.IpoBid{
			Identity: identity.String(),
			Amount:   source.Prices[index],
		}
	}

	return &ipoBidData, nil
}
