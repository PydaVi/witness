package main

import (
	"flag"
	"log"
	"net"
	"os"

	"witness/internal/config"
	"witness/internal/listener"
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

	log.Printf("config loaded: listener=%s upstreams=%d routes=%d",
		cfg.Listener.Addr, len(cfg.Upstreams), len(cfg.Routes))

	// Nesta etapa, apenas aceitamos conexoes e logamos o remoto.
	// A partir da proxima etapa, vamos parsear HTTP e encaminhar.
	handler := func(conn net.Conn) {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("close connection: %v", err)
			}
		}()

		remote := conn.RemoteAddr()
		log.Printf("accepted connection: remote=%s", remote.String())
	}

	if err := listener.ListenAndServe(cfg.Listener.Addr, cfg.Listener.Backlog, handler); err != nil {
		log.Printf("listener error: %v", err)
		os.Exit(1)
	}
}
