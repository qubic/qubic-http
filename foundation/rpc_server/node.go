package rpc

import (
	"context"

	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector/v2"
	"github.com/qubic/go-node-connector/v2/types"
)

// nodeClient is the subset of *qubic.Client used by the handlers. It exists so that
// handlers can be unit tested against a fake instead of a live node connection.
type nodeClient interface {
	GetIdentity(ctx context.Context, id string) (types.AddressInfo, error)
	GetTickInfo(ctx context.Context) (types.TickInfo, error)
	QuerySmartContract(ctx context.Context, rcf qubic.RequestContractFunction, requestData []byte) (types.SmartContractData, error)
	SendRawTransaction(ctx context.Context, rawTx []byte) error
	GetIssuedAssets(ctx context.Context, id string) (types.IssuedAssets, error)
	GetOwnedAssets(ctx context.Context, id string) (types.OwnedAssets, error)
	GetPossessedAssets(ctx context.Context, id string) (types.PossessedAssets, error)
	PrefetchOwnedAndPossessedAssets(ctx context.Context, ids []string) ([]qubic.AddressAssets, error)
	GetAssetIssuancesByUniverseIndex(ctx context.Context, index uint32) (types.AssetIssuances, error)
	GetAssetOwnershipsByUniverseIndex(ctx context.Context, index uint32) (types.AssetOwnerships, error)
	GetAssetPossessionsByUniverseIndex(ctx context.Context, index uint32) (types.AssetPossessions, error)
	GetAssetIssuancesByFilter(ctx context.Context, issuerIdentity, assetName string) (types.AssetIssuances, error)
	GetAssetOwnershipsByFilter(ctx context.Context, issuerIdentity, assetName, ownerIdentity string, ownerContract uint16) (types.AssetOwnerships, error)
	GetAssetPossessionsByFilter(ctx context.Context, issuerIdentity, assetName, ownerIdentity, possessorIdentity string, ownerContract, possessorContract uint16) (types.AssetPossessions, error)
	GetActiveIpos(ctx context.Context) (types.Ipos, error)
	GetContractIpo(ctx context.Context, contractIndex uint32) (types.ContractIpo, error)
}

// compile time check that the real client still satisfies the seam
var _ nodeClient = (*qubic.Client)(nil)

// nodePool mirrors the Get/Put/Close contract of *qubic.Pool. Put returns a client for
// reuse, Close discards it. A client whose request failed must always be closed and
// never put back, since its connection may be left unusable.
type nodePool interface {
	Get() (nodeClient, error)
	Put(client nodeClient) error
	Close(client nodeClient) error
}

// qubicNodePool adapts *qubic.Pool, whose Get returns the concrete *qubic.Client, to
// the nodePool interface.
type qubicNodePool struct {
	pool *qubic.Pool
}

func newQubicNodePool(pool *qubic.Pool) *qubicNodePool {
	return &qubicNodePool{pool: pool}
}

func (p *qubicNodePool) Get() (nodeClient, error) {
	client, err := p.pool.Get()
	if err != nil {
		// returning client directly would produce a non nil interface holding a nil
		// *qubic.Client, which callers could not nil check
		return nil, err
	}
	return client, nil
}

func (p *qubicNodePool) Put(client nodeClient) error {
	c, ok := client.(*qubic.Client)
	if !ok {
		return errors.Errorf("putting back unexpected client type %T", client)
	}
	return p.pool.Put(c)
}

func (p *qubicNodePool) Close(client nodeClient) error {
	c, ok := client.(*qubic.Client)
	if !ok {
		return errors.Errorf("closing unexpected client type %T", client)
	}
	return p.pool.Close(c)
}
