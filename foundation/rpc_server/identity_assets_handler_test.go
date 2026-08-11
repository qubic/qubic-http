package rpc

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector/v2"
	"github.com/qubic/go-node-connector/v2/types"
	"github.com/qubic/qubic-http/protobuff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultTestBatchLimit = 15

func assetsResponse(t *testing.T, identities ...string) []qubic.AddressAssets {
	t.Helper()

	response := make([]qubic.AddressAssets, 0, len(identities))
	for _, identity := range identities {
		response = append(response, qubic.AddressAssets{
			Identity: identity,
			Owned: types.OwnedAssets{
				{Data: testOwnedAssetData(t), Info: types.AssetInfo{Tick: 42, UniverseIndex: 43}},
			},
			Possessed: types.PossessedAssets{},
		})
	}

	return response
}

// requireInvalidArgument asserts the request was rejected before any node connection was
// taken, so bad input never costs a pooled connection.
func requireInvalidArgument(t *testing.T, pool *fakeNodePool, err error) {
	t.Helper()

	st, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error, got %v", err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Zero(t, pool.gets, "no pool connection should be taken for an invalid request")
}

func TestServer_GetAssetsForIdentities_success(t *testing.T) {
	identities := []string{ownerIdentity, possessorIdentity}

	client := &fakeNodeClient{
		prefetch: func(_ context.Context, ids []string) ([]qubic.AddressAssets, error) {
			return assetsResponse(t, ids...), nil
		},
	}
	pool := &fakeNodePool{client: client}

	response, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{Identities: identities})
	require.NoError(t, err)

	assert.Equal(t, identities, client.prefetchedIdentities, "identities should be forwarded verbatim")

	require.Len(t, response.Assets, 2)
	assert.Equal(t, ownerIdentity, response.Assets[0].Identity)
	assert.Equal(t, possessorIdentity, response.Assets[1].Identity)
	require.Len(t, response.Assets[0].OwnedAssets, 1)
	assert.Equal(t, expectedOwnedAssetData(), response.Assets[0].OwnedAssets[0].Data)

	// one batch means exactly one connection, returned to the pool for reuse
	assert.Equal(t, 1, pool.gets)
	assert.Equal(t, 1, pool.puts)
	assert.Zero(t, pool.closes)
}

func TestServer_GetAssetsForIdentities_identityWithoutAssets(t *testing.T) {
	client := &fakeNodeClient{
		prefetch: func(_ context.Context, ids []string) ([]qubic.AddressAssets, error) {
			return []qubic.AddressAssets{{
				Identity:  ids[0],
				Owned:     types.OwnedAssets{},
				Possessed: types.PossessedAssets{},
			}}, nil
		},
	}
	pool := &fakeNodePool{client: client}

	response, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{Identities: []string{ownerIdentity}})
	require.NoError(t, err)

	// the identity is still represented, with empty lists rather than being omitted
	require.Len(t, response.Assets, 1)
	assert.Equal(t, ownerIdentity, response.Assets[0].Identity)
	assert.Empty(t, response.Assets[0].OwnedAssets)
	assert.Empty(t, response.Assets[0].PossessedAssets)
}

func TestServer_GetAssetsForIdentities_emptyIdentities(t *testing.T) {
	pool := &fakeNodePool{client: &fakeNodeClient{}}

	_, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{Identities: nil})

	requireInvalidArgument(t, pool, err)
}

func TestServer_GetAssetsForIdentities_tooManyIdentities(t *testing.T) {
	pool := &fakeNodePool{client: &fakeNodeClient{}}

	_, err := newTestServer(pool, 1).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{
			Identities: []string{ownerIdentity, possessorIdentity},
		})

	requireInvalidArgument(t, pool, err)
	assert.Contains(t, status.Convert(err).Message(), "maximum is 1")
}

func TestServer_GetAssetsForIdentities_invalidIdentity(t *testing.T) {
	pool := &fakeNodePool{client: &fakeNodeClient{}}

	// a single malformed entry rejects the whole batch, so nothing is fetched
	_, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{
			Identities: []string{ownerIdentity, "not-an-identity"},
		})

	requireInvalidArgument(t, pool, err)
	assert.Contains(t, status.Convert(err).Message(), "not-an-identity")
}

func TestServer_GetAssetsForIdentities_duplicateIdentities(t *testing.T) {
	pool := &fakeNodePool{client: &fakeNodeClient{}}

	_, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{
			Identities: []string{ownerIdentity, ownerIdentity},
		})

	requireInvalidArgument(t, pool, err)
	assert.Contains(t, status.Convert(err).Message(), "duplicate identity")
}

func TestServer_GetAssetsForIdentities_poolGetError(t *testing.T) {
	pool := &fakeNodePool{getErr: errors.New("pool exhausted")}

	_, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{Identities: []string{ownerIdentity}})

	assert.Equal(t, codes.Internal, status.Code(err))
	// nothing was handed out, so nothing may be returned or discarded
	assert.Zero(t, pool.puts)
	assert.Zero(t, pool.closes)
}

// TestServer_GetAssetsForIdentities_prefetchError pins the pool safety rule: a failed
// pipelined batch leaves the connection byte-unaligned, so it must be discarded rather
// than put back where another request would pick it up.
func TestServer_GetAssetsForIdentities_prefetchError(t *testing.T) {
	client := &fakeNodeClient{
		prefetch: func(context.Context, []string) ([]qubic.AddressAssets, error) {
			return nil, errors.New("connection reset")
		},
	}
	pool := &fakeNodePool{client: client}

	_, err := newTestServer(pool, defaultTestBatchLimit).
		GetAssetsForIdentities(context.Background(), &protobuff.GetAssetsForIdentitiesRequest{Identities: []string{ownerIdentity}})

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(t, 1, pool.closes, "a failed batch must close the connection")
	assert.Zero(t, pool.puts, "a failed batch must never return the connection to the pool")
}
