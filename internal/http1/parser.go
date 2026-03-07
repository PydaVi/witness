package http1

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	// maxLineBytes limita o tamanho de uma linha (request line ou header).
	// Esse limite protege contra abusos e facilita debug.
	maxLineBytes = 8 * 1024
	// maxHeaders limita a quantidade de headers por request.
	maxHeaders = 100
)

// Request representa o minimo que precisamos para roteamento na v0.1.
type Request struct {
	Method  string
	Path    string
	Proto   string
	Headers map[string]string
}

// ParseRequest le uma request HTTP/1.1 a partir de um bufio.Reader.
// Nesta etapa, parseamos apenas request line + headers.
func ParseRequest(r *bufio.Reader) (*Request, error) {
	line, err := readLine(r)
	if err != nil {
		return nil, fmt.Errorf("read request line: %w", err)
	}
	if line == "" {
		return nil, fmt.Errorf("empty request line")
	}

	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line %q", line)
	}

	method := parts[0]
	path := parts[1]
	proto := parts[2]
	if proto != "HTTP/1.1" {
		return nil, fmt.Errorf("unsupported protocol %q", proto)
	}

	req := &Request{
		Method:  method,
		Path:    path,
		Proto:   proto,
		Headers: make(map[string]string),
	}

	for i := 0; i < maxHeaders; i++ {
		headerLine, err := readLine(r)
		if err != nil {
			return nil, fmt.Errorf("read header line: %w", err)
		}
		if headerLine == "" {
			return req, nil
		}

		colon := strings.Index(headerLine, ":")
		if colon <= 0 {
			return nil, fmt.Errorf("invalid header line %q", headerLine)
		}

		key := strings.TrimSpace(headerLine[:colon])
		value := strings.TrimSpace(headerLine[colon+1:])
		if key == "" {
			return nil, fmt.Errorf("empty header name")
		}

		// Armazenamos headers em lowercase para facilitar lookup (ex: host).
		req.Headers[strings.ToLower(key)] = value
	}

	return nil, fmt.Errorf("too many headers (max=%d)", maxHeaders)
}

// readLine le uma linha terminada em \n, respeitando maxLineBytes.
// Implementamos manualmente para ter controle de limite e erros.
func readLine(r *bufio.Reader) (string, error) {
	var buf []byte

	for {
		frag, err := r.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			if len(buf)+len(frag) > maxLineBytes {
				return "", fmt.Errorf("line too long (max=%d)", maxLineBytes)
			}
			buf = append(buf, frag...)
			continue
		}
		if err != nil {
			return "", err
		}

		if len(buf)+len(frag) > maxLineBytes {
			return "", fmt.Errorf("line too long (max=%d)", maxLineBytes)
		}

		buf = append(buf, frag...)
		break
	}

	line := string(buf)
	line = strings.TrimRight(line, "\r\n")
	return line, nil
}
