package nodes

import (
	"math/rand"
)

type Pool struct {
	nodeIPs []string
}

func NewPool(nodeIPs []string) *Pool {
	return &Pool{nodeIPs: nodeIPs}
}

func (p *Pool) GetRandomIP() string {
	index := rand.Intn(len(p.nodeIPs))

	return p.nodeIPs[index]
}
