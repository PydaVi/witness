package main

import (
	"bufio"
	"flag"
	"log/slog"
	"net"
	"os"
	"time"

	"witness/internal/balancer"
	"witness/internal/config"
	"witness/internal/http1"
	"witness/internal/listener"
	"witness/internal/logging"
	"witness/internal/proxy"
	"witness/internal/router"
)

const defaultConfigPath = "configs/dev.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "path to config YAML")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	logger := logging.New()

	logger.Info("config loaded",
		"listener", cfg.Listener.Addr,
		"upstreams", len(cfg.Upstreams),
		"routes", len(cfg.Routes),
	)

	balancers := make(map[string]*balancer.RoundRobin, len(cfg.Upstreams))
	for _, upstream := range cfg.Upstreams {
		rr, err := balancer.NewRoundRobin(upstream.Targets)
		if err != nil {
			logger.Error("balancer error", "upstream", upstream.Name, "err", err)
			os.Exit(1)
		}
		balancers[upstream.Name] = rr
	}

	r := router.New(cfg.Routes)
	p := &proxy.Proxy{
		DialTimeout:  cfg.Timeouts.Connect.Duration,
		ReadTimeout:  cfg.Timeouts.Read.Duration,
		WriteTimeout: cfg.Timeouts.Write.Duration,
	}

	handler := func(conn net.Conn) {
		start := time.Now()
		defer func() {
			if err := conn.Close(); err != nil {
				logger.Error("close connection", "err", err)
			}
		}()

		remote := conn.RemoteAddr()
		reader := bufio.NewReader(conn)

		if cfg.Timeouts.Read.Duration > 0 {
			if err := conn.SetReadDeadline(time.Now().Add(cfg.Timeouts.Read.Duration)); err != nil {
				logger.Error("set client read deadline", "remote", remote.String(), "err", err)
				return
			}
		}

		req, err := http1.ParseRequest(reader)
		if err != nil {
			latency := time.Since(start).Milliseconds()
			logger.Error("request error",
				"remote", remote.String(),
				"latency_ms", latency,
				"err", err,
			)
			return
		}

		host := req.Headers["host"]
		upstream, err := r.Match(host, req.Path)
		if err != nil {
			latency := time.Since(start).Milliseconds()
			logger.Error("route not found",
				"remote", remote.String(),
				"host", host,
				"path", req.Path,
				"latency_ms", latency,
				"err", err,
			)
			return
		}

		rr, ok := balancers[upstream]
		if !ok {
			logger.Error("balancer not found", "remote", remote.String(), "upstream", upstream)
			return
		}

		target := rr.Next()

		if err := p.Forward(conn, target, req); err != nil {
			latency := time.Since(start).Milliseconds()
			logger.Error("proxy error",
				"remote", remote.String(),
				"method", req.Method,
				"path", req.Path,
				"host", host,
				"upstream", upstream,
				"target", target,
				"latency_ms", latency,
				"err", err,
			)
			return
		}

		latency := time.Since(start).Milliseconds()
		logger.Info("request",
			"remote", remote.String(),
			"method", req.Method,
			"path", req.Path,
			"host", host,
			"upstream", upstream,
			"target", target,
			"latency_ms", latency,
		)
	}

	if err := listener.ListenAndServe(cfg.Listener.Addr, cfg.Listener.Backlog, handler); err != nil {
		logger.Error("listener error", "err", err)
		os.Exit(1)
	}
}
