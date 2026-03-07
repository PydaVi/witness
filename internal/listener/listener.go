package listener

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

// Handler representa a funcao que processa uma conexao aceita.
// Na v0.1 ela apenas loga e fecha, mas o tipo evita acoplamento com o accept loop.
type Handler func(conn net.Conn)

// ListenAndServe cria um listener TCP com backlog configuravel e entra no accept loop.
// Usamos syscall direto para controlar o backlog, mantendo o comportamento transparente e didatico.
func ListenAndServe(addr string, backlog int, handler Handler) error {
	ln, err := newTCPListener(addr, backlog)
	if err != nil {
		return err
	}
	defer func() {
		if err := ln.Close(); err != nil {
			log.Printf("close listener: %v", err)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// Erros temporarios podem acontecer (ex: falta de recursos momentanea).
			// Nesse caso, esperamos um pouco e seguimos aceitando.
			if isTemporary(err) {
				log.Printf("accept temporary error: %v", err)
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return fmt.Errorf("accept connection: %w", err)
		}

		go handler(conn)
	}
}

// newTCPListener cria um socket de escuta via syscall para permitir backlog explicito.
func newTCPListener(addr string, backlog int) (net.Listener, error) {
	resolved, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("resolve tcp addr %q: %w", addr, err)
	}

	family := socketFamily(resolved)
	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, fmt.Errorf("create socket: %w", err)
	}

	// Se algo falhar depois daqui, garantimos o fechamento do fd.
	closeOnError := true
	defer func() {
		if closeOnError {
			_ = syscall.Close(fd)
		}
	}()

	// SO_REUSEADDR facilita reiniciar o processo sem esperar o TIME_WAIT expirar.
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return nil, fmt.Errorf("set SO_REUSEADDR: %w", err)
	}

	if err := syscall.Bind(fd, toSockaddr(resolved, family)); err != nil {
		return nil, fmt.Errorf("bind %q: %w", addr, err)
	}

	if backlog <= 0 {
		backlog = syscall.SOMAXCONN
	}

	if err := syscall.Listen(fd, backlog); err != nil {
		return nil, fmt.Errorf("listen backlog=%d: %w", backlog, err)
	}

	file := os.NewFile(uintptr(fd), "listener")
	if file == nil {
		return nil, fmt.Errorf("create os.File from fd")
	}

	ln, err := net.FileListener(file)
	if err != nil {
		return nil, fmt.Errorf("convert fd to net.Listener: %w", err)
	}

	// net.FileListener duplica o fd internamente; podemos fechar o original.
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close original fd: %w", err)
	}

	closeOnError = false
	return ln, nil
}

// socketFamily escolhe IPv4 ou IPv6 com base no endereco resolvido.
func socketFamily(addr *net.TCPAddr) int {
	if addr == nil || addr.IP == nil {
		return syscall.AF_INET
	}
	if addr.IP.To4() != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

// toSockaddr converte net.TCPAddr para o tipo esperado pelo syscall.Bind.
func toSockaddr(addr *net.TCPAddr, family int) syscall.Sockaddr {
	if family == syscall.AF_INET6 {
		sa := &syscall.SockaddrInet6{Port: addr.Port}
		if addr.IP != nil {
			copy(sa.Addr[:], addr.IP.To16())
		}
		return sa
	}

	sa := &syscall.SockaddrInet4{Port: addr.Port}
	if addr.IP != nil {
		copy(sa.Addr[:], addr.IP.To4())
	}
	return sa
}

func isTemporary(err error) bool {
	nerr, ok := err.(net.Error)
	if !ok {
		return false
	}
	return nerr.Temporary()
}
