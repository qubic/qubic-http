package rpc

import (
	"fmt"

	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector/v2"
	"github.com/qubic/go-node-connector/v2/types"
	"github.com/qubic/qubic-http/protobuff"
)

func convertAssetInfo(source types.AssetInfo) *protobuff.AssetInfo {
	return &protobuff.AssetInfo{
		Tick:          source.Tick,
		UniverseIndex: source.UniverseIndex,
	}
}

func convertIssuedAssetData(source types.IssuedAssetData) (*protobuff.IssuedAssetData, error) {
	var issuer types.Identity
	issuer, err := issuer.FromPubKey(source.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for issued asset public key")
	}

	return &protobuff.IssuedAssetData{
		IssuerIdentity:        issuer.String(),
		Type:                  uint32(source.Type),
		Name:                  int8ArrayToString(source.Name[:]),
		NumberOfDecimalPlaces: int32(source.NumberOfDecimalPlaces),
		UnitOfMeasurement:     int8ArrayToInt32Array(source.UnitOfMeasurement[:]),
	}, nil
}

func convertOwnedAssetData(source types.OwnedAssetData) (*protobuff.OwnedAssetData, error) {
	var owner types.Identity
	owner, err := owner.FromPubKey(source.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for owned asset public key")
	}

	issuedAsset, err := convertIssuedAssetData(source.IssuedAsset)
	if err != nil {
		return nil, errors.Wrap(err, "converting issued asset")
	}

	return &protobuff.OwnedAssetData{
		OwnerIdentity:         owner.String(),
		Type:                  uint32(source.Type),
		Padding:               int32(source.Padding[0]),
		ManagingContractIndex: uint32(source.ManagingContractIndex),
		IssuanceIndex:         source.IssuanceIndex,
		NumberOfUnits:         source.NumberOfUnits,
		IssuedAsset:           issuedAsset,
	}, nil
}

func convertPossessedAssetData(source types.PossessedAssetData) (*protobuff.PossessedAssetData, error) {
	var possessor types.Identity
	possessor, err := possessor.FromPubKey(source.PublicKey, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity for possessed asset public key")
	}

	// the nested ownership record is converted from its own data, not from the
	// possession record that contains it
	ownedAsset, err := convertOwnedAssetData(source.OwnedAsset)
	if err != nil {
		return nil, errors.Wrap(err, "converting owned asset")
	}

	return &protobuff.PossessedAssetData{
		PossessorIdentity:     possessor.String(),
		Type:                  uint32(source.Type),
		Padding:               int32(source.Padding[0]),
		ManagingContractIndex: uint32(source.ManagingContractIndex),
		IssuanceIndex:         source.IssuanceIndex,
		NumberOfUnits:         source.NumberOfUnits,
		OwnedAsset:            ownedAsset,
	}, nil
}

func convertIssuedAsset(source types.IssuedAsset) (*protobuff.IssuedAsset, error) {
	data, err := convertIssuedAssetData(source.Data)
	if err != nil {
		return nil, err
	}

	return &protobuff.IssuedAsset{Data: data, Info: convertAssetInfo(source.Info)}, nil
}

func convertOwnedAsset(source types.OwnedAsset) (*protobuff.OwnedAsset, error) {
	data, err := convertOwnedAssetData(source.Data)
	if err != nil {
		return nil, err
	}

	return &protobuff.OwnedAsset{Data: data, Info: convertAssetInfo(source.Info)}, nil
}

func convertPossessedAsset(source types.PossessedAsset) (*protobuff.PossessedAsset, error) {
	data, err := convertPossessedAssetData(source.Data)
	if err != nil {
		return nil, err
	}

	return &protobuff.PossessedAsset{Data: data, Info: convertAssetInfo(source.Info)}, nil
}

func convertIssuedAssets(source types.IssuedAssets) ([]*protobuff.IssuedAsset, error) {
	converted := make([]*protobuff.IssuedAsset, 0, len(source))
	for _, asset := range source {
		issuedAsset, err := convertIssuedAsset(asset)
		if err != nil {
			return nil, err
		}
		converted = append(converted, issuedAsset)
	}

	return converted, nil
}

func convertOwnedAssets(source types.OwnedAssets) ([]*protobuff.OwnedAsset, error) {
	converted := make([]*protobuff.OwnedAsset, 0, len(source))
	for _, asset := range source {
		ownedAsset, err := convertOwnedAsset(asset)
		if err != nil {
			return nil, err
		}
		converted = append(converted, ownedAsset)
	}

	return converted, nil
}

func convertPossessedAssets(source types.PossessedAssets) ([]*protobuff.PossessedAsset, error) {
	converted := make([]*protobuff.PossessedAsset, 0, len(source))
	for _, asset := range source {
		possessedAsset, err := convertPossessedAsset(asset)
		if err != nil {
			return nil, err
		}
		converted = append(converted, possessedAsset)
	}

	return converted, nil
}

// convertAddressAssets maps the batched node response into the proto response, keeping
// the order the identities were requested in.
func convertAddressAssets(source []qubic.AddressAssets) ([]*protobuff.IdentityAssets, error) {
	converted := make([]*protobuff.IdentityAssets, 0, len(source))
	for _, addressAssets := range source {
		owned, err := convertOwnedAssets(addressAssets.Owned)
		if err != nil {
			return nil, errors.Wrapf(err, "converting owned assets for identity %s", addressAssets.Identity)
		}

		possessed, err := convertPossessedAssets(addressAssets.Possessed)
		if err != nil {
			return nil, errors.Wrapf(err, "converting possessed assets for identity %s", addressAssets.Identity)
		}

		converted = append(converted, &protobuff.IdentityAssets{
			Identity:        addressAssets.Identity,
			OwnedAssets:     owned,
			PossessedAssets: possessed,
		})
	}

	return converted, nil
}

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
