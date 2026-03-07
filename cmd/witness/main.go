package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"witness/internal/config"
)

const defaultConfigPath = "configs/dev.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "path to config YAML")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("config error: %v", err)
		os.Exit(1)
	}

	// Neste ponto so validamos a configuracao.
	// O listener e o proxy virão nas proximas etapas.
	fmt.Printf("config loaded: listener=%s upstreams=%d routes=%d\n",
		cfg.Listener.Addr,
		len(cfg.Upstreams),
		len(cfg.Routes),
	)
}
