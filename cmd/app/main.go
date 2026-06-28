package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Sayfargo/url-shortener/internal/app"
	core_config "github.com/Sayfargo/url-shortener/internal/core/config"
)

func main() {

	cfg := core_config.Mustload()

	app := app.MustNew(cfg)

	go app.GRPCServer.Run()

	sigCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	<-sigCtx.Done()

	app.GRPCServer.Stop()

}
