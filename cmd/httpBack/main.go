package main

import (
	"budget/internal/app"
	"budget/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	app := app.NewApp(cfg)
	app.Run()
}
