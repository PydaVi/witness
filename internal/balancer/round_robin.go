package balancer

import (
	"fmt"
	"sync"
)

// RoundRobin implementa distribuicao ciclica entre targets.
// Este e o padrao "round-robin com mutex" para garantir acesso seguro em concorrencia.
type RoundRobin struct {
	mu      sync.Mutex
	targets []string
	index   int
}

// NewRoundRobin cria um balanceador round-robin com lista fixa de targets.
func NewRoundRobin(targets []string) (*RoundRobin, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("round-robin requires at least one target")
	}

	copied := make([]string, len(targets))
	copy(copied, targets)

	return &RoundRobin{targets: copied}, nil
}

// Next devolve o proximo target e avanca o indice de forma circular.
func (rr *RoundRobin) Next() string {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	target := rr.targets[rr.index]
	rr.index = (rr.index + 1) % len(rr.targets)
	return target
}
