package main

import (
	"log/slog"
	"os"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/components"
	"github.com/lmittmann/tint"
)

func main() {
	raptor.New(
		components.New(),
		config.Routes(),
		raptor.WithLogHandler(func(level *slog.LevelVar) slog.Handler {
			return tint.NewHandler(os.Stderr, &tint.Options{Level: level})
		}),
	).Run()
}
