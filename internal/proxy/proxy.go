package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"witness/internal/http1"
)

// Proxy implementa o encaminhamento minimo de uma request para um backend.
// Nesta etapa nao suportamos body; apenas request line + headers.
type Proxy struct {
	DialTimeout time.Duration
}

// Forward conecta ao backend, envia a request e encaminha a resposta ao cliente.
func (p *Proxy) Forward(client net.Conn, target string, req *http1.Request) error {
	if hasBody(req.Headers) {
		return fmt.Errorf("request body not supported yet")
	}

	backend, err := net.DialTimeout("tcp", target, p.DialTimeout)
	if err != nil {
		return fmt.Errorf("connect to backend %s: %w", target, err)
	}
	defer func() {
		_ = backend.Close()
	}()

	backendWriter := bufio.NewWriter(backend)
	if err := writeRequest(backendWriter, req); err != nil {
		return fmt.Errorf("write request to backend %s: %w", target, err)
	}
	if err := backendWriter.Flush(); err != nil {
		return fmt.Errorf("flush request to backend %s: %w", target, err)
	}

	if _, err := io.Copy(client, backend); err != nil {
		return fmt.Errorf("copy response from backend %s: %w", target, err)
	}

	return nil
}

// hasBody detecta se a request tem body baseado em headers conhecidos.
// Nesta etapa, qualquer body e considerado erro explicito.
func hasBody(headers map[string]string) bool {
	if value, ok := headers["transfer-encoding"]; ok {
		return strings.TrimSpace(strings.ToLower(value)) != ""
	}

	if value, ok := headers["content-length"]; ok {
		value = strings.TrimSpace(value)
		if value == "" {
			return false
		}
		n, err := strconv.Atoi(value)
		if err != nil {
			return true
		}
		return n > 0
	}

	return false
}

// writeRequest reconstroi a request a partir do parser.
// Mantemos os headers exatamente como chegaram, exceto por garantir CRLF.
func writeRequest(w *bufio.Writer, req *http1.Request) error {
	if _, err := fmt.Fprintf(w, "%s %s %s\r\n", req.Method, req.Path, req.Proto); err != nil {
		return err
	}

	for key, value := range req.Headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", key, value); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("\r\n"); err != nil {
		return err
	}

	return nil
}
