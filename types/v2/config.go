package types

import (
	"log/slog"
)

type Config struct {
	Logger *slog.Logger
}

var (
	DefaultConfig = Config{
		Logger: slog.Default(),
	}
)
