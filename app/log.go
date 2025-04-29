package app

import (
	"github.com/Qitmeer/llama.go/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"log/slog"
	"os"
)

func initLog(cfg *config.Config) error {
	logOutput := os.Stdout
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(logOutput)
	if usecolor {
		output = colorable.NewColorable(logOutput)
	}
	verbosity := parseLevel(cfg.LogLevel)
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(output, verbosity, usecolor)))
	return nil
}

func parseLevel(level string) slog.Level {
	switch level {
	case "trace":
		return slog.LevelDebug
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
