package rpc

import (
	"context"
	"io"
	"log"

	qubic "github.com/qubic/go-node-connector/v2"
)

// fakeNodeClient embeds nodeClient so only the methods a test actually needs have to be
// implemented. Any other call panics, which surfaces unexpected node usage loudly.
type fakeNodeClient struct {
	nodeClient
	prefetchedIdentities []string
	prefetch             func(ctx context.Context, ids []string) ([]qubic.AddressAssets, error)
}

func (f *fakeNodeClient) PrefetchOwnedAndPossessedAssets(ctx context.Context, ids []string) ([]qubic.AddressAssets, error) {
	f.prefetchedIdentities = ids
	return f.prefetch(ctx, ids)
}

// fakeNodePool records how connections were handed out and returned so tests can assert
// that a failed request closes the connection instead of putting it back in the pool.
type fakeNodePool struct {
	client nodeClient
	getErr error

	gets   int
	puts   int
	closes int
}

func (p *fakeNodePool) Get() (nodeClient, error) {
	p.gets++
	if p.getErr != nil {
		return nil, p.getErr
	}
	return p.client, nil
}

func (p *fakeNodePool) Put(nodeClient) error {
	p.puts++
	return nil
}

func (p *fakeNodePool) Close(nodeClient) error {
	p.closes++
	return nil
}

func newTestServer(pool nodePool, maxBatchIdentities int) *Server {
	return &Server{
		logger:             log.New(io.Discard, "", 0),
		qPool:              pool,
		maxBatchIdentities: maxBatchIdentities,
	}
}
