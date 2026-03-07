package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"

	"witness/internal/config"
	"witness/internal/http1"
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
		reader := bufio.NewReader(conn)

		req, err := http1.ParseRequest(reader)
		if err != nil {
			log.Printf("parse request error: remote=%s err=%v", remote.String(), err)
			return
		}

		host := req.Headers["host"]
		log.Printf("request: remote=%s method=%s path=%s host=%s", remote.String(), req.Method, req.Path, host)
	}

	if err := listener.ListenAndServe(cfg.Listener.Addr, cfg.Listener.Backlog, handler); err != nil {
		log.Printf("listener error: %v", err)
		os.Exit(1)
	}
}
