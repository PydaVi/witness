package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config e o modelo raiz de configuracao.
// Mantemos tudo em structs para validar cedo e tornar a configuracao explicita.
type Config struct {
	Listener  ListenerConfig  `yaml:"listener"`
	Upstreams []Upstream      `yaml:"upstreams"`
	Routes    []Route         `yaml:"routes"`
}

// ListenerConfig descreve onde o proxy vai escutar conexoes TCP.
type ListenerConfig struct {
	Addr    string `yaml:"addr"`
	Backlog int    `yaml:"backlog"`
}

// Upstream representa um conjunto de targets que recebem trafego.
type Upstream struct {
	Name    string   `yaml:"name"`
	Targets []string `yaml:"targets"`
}

// Route descreve a regra de roteamento simples por host + prefixo de path.
type Route struct {
	Host       string `yaml:"host"`
	PathPrefix string `yaml:"path_prefix"`
	Upstream   string `yaml:"upstream"`
}

// Load carrega um arquivo YAML e valida o conteudo.
// A validacao e feita aqui para falhar cedo no boot, antes do listener iniciar.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate garante invariantes minimos para evitar comportamento implicito.
func (c *Config) Validate() error {
	if c.Listener.Addr == "" {
		return fmt.Errorf("listener.addr is required")
	}

	if len(c.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream is required")
	}

	upstreamByName := make(map[string]struct{}, len(c.Upstreams))
	for i, upstream := range c.Upstreams {
		if upstream.Name == "" {
			return fmt.Errorf("upstreams[%d].name is required", i)
		}
		if len(upstream.Targets) == 0 {
			return fmt.Errorf("upstreams[%d].targets must have at least one target", i)
		}
		if _, exists := upstreamByName[upstream.Name]; exists {
			return fmt.Errorf("upstream %q is duplicated", upstream.Name)
		}
		upstreamByName[upstream.Name] = struct{}{}
	}

	if len(c.Routes) == 0 {
		return fmt.Errorf("at least one route is required")
	}

	for i, route := range c.Routes {
		if route.Host == "" {
			return fmt.Errorf("routes[%d].host is required", i)
		}
		if route.PathPrefix == "" {
			return fmt.Errorf("routes[%d].path_prefix is required", i)
		}
		if route.Upstream == "" {
			return fmt.Errorf("routes[%d].upstream is required", i)
		}
		if _, exists := upstreamByName[route.Upstream]; !exists {
			return fmt.Errorf("routes[%d].upstream %q does not exist", i, route.Upstream)
		}
	}

	return nil
}
