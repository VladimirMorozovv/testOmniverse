package main

import (
	"context"
	"flag"
	"github.com/SedovSG/zaplog"

	"os"
	"os/signal"
	"service_api/internal/apiserver"
	"service_api/internal/config"
	"syscall"
)

func main() {
	configFile := flag.String("config", "configs/config.yml", "Path to config file.")
	flag.Parse()
	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		zaplog.Throw().Fatal(err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := apiserver.New(ctx, cfg)
	if err != nil {
		zaplog.Throw().Fatal(err.Error())
	}
	if err := s.Start(); err != nil {
		zaplog.Throw().Fatal(err.Error())
	}

	// Gracefully shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		zaplog.Throw().Info("shutdown application")
		cancel()
	}()
}
