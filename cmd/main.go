package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/feynmaz/pkg/logger"
	"github.com/feynmaz/promfeeder/config"
	"github.com/feynmaz/promfeeder/internal/server"
)

func main() {
	log := logger.New()

	cfg, err := config.GetDefault()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config")
	}
	log.SetLevel(cfg.LogLevel)

	prettyJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal config")
	}
	log.Debug().Msgf("Config: %s", string(prettyJSON))

	server := server.New(cfg, log)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := server.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	server.Shutdown()
}
