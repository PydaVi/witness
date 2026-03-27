package logging

import (
	"log/slog"
	"os"
)

// New cria um logger estruturado basico.
// Mantemos o handler simples nesta fase para focar em formato consistente.
func New() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	return slog.New(handler)
}
