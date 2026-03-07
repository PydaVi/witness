package router

import (
	"fmt"
	"strings"

	"witness/internal/config"
)

// Router implementa roteamento simples por host + prefixo de path.
// A regra e: primeira rota que casar vence (deterministico e facil de entender).
type Router struct {
	routes []config.Route
}

// New cria um roteador a partir das rotas da config.
func New(routes []config.Route) *Router {
	// Copiamos para evitar que alteracoes externas mudem o comportamento em runtime.
	copied := make([]config.Route, len(routes))
	copy(copied, routes)

	return &Router{routes: copied}
}

// Match encontra a primeira rota que casa com host + path.
func (r *Router) Match(host, path string) (string, error) {
	normalizedHost := strings.ToLower(host)

	for _, route := range r.routes {
		if strings.ToLower(route.Host) != normalizedHost {
			continue
		}
		if !strings.HasPrefix(path, route.PathPrefix) {
			continue
		}

		return route.Upstream, nil
	}

	return "", fmt.Errorf("no route for host=%q path=%q", host, path)
}
