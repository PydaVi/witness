package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"witness/internal/balancer"
	"witness/internal/config"
	"witness/internal/http1"
	"witness/internal/listener"
	"witness/internal/proxy"
	"witness/internal/router"
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

	balancers := make(map[string]*balancer.RoundRobin, len(cfg.Upstreams))
	for _, upstream := range cfg.Upstreams {
		rr, err := balancer.NewRoundRobin(upstream.Targets)
		if err != nil {
			log.Printf("balancer error: upstream=%s err=%v", upstream.Name, err)
			os.Exit(1)
		}
		balancers[upstream.Name] = rr
	}

	r := router.New(cfg.Routes)
	p := &proxy.Proxy{DialTimeout: 2 * time.Second}

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
		upstream, err := r.Match(host, req.Path)
		if err != nil {
			log.Printf("route not found: remote=%s host=%s path=%s err=%v", remote.String(), host, req.Path, err)
			return
		}

		rr, ok := balancers[upstream]
		if !ok {
			log.Printf("balancer not found: remote=%s upstream=%s", remote.String(), upstream)
			return
		}

		target := rr.Next()

		if err := p.Forward(conn, target, req); err != nil {
			log.Printf("proxy error: remote=%s target=%s err=%v", remote.String(), target, err)
			return
		}

		log.Printf("request: remote=%s method=%s path=%s host=%s upstream=%s target=%s",
			remote.String(), req.Method, req.Path, host, upstream, target)
	}

	if err := listener.ListenAndServe(cfg.Listener.Addr, cfg.Listener.Backlog, handler); err != nil {
		log.Printf("listener error: %v", err)
		os.Exit(1)
	}
}
