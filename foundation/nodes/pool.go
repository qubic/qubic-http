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

func (p *Pool) GetMaxTargetRandomIPs(target int) []string {
	if len(p.nodeIPs) <= target {
		return p.nodeIPs
	}

	ipsSet := make(map[string]struct{})

	for {
		ipsSet[p.nodeIPs[rand.Intn(len(p.nodeIPs))]] = struct{}{}
		if len(ipsSet) == target {
			break
		}
	}

	ips := make([]string, 0, target)
	for ip := range ipsSet {
		ips = append(ips, ip)
	}

	return ips
}
