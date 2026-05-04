package slog

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {

	newLog := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		),
	)

	return newLog
}
