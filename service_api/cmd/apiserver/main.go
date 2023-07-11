package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"service_api/internal/apiserver"
	"service_api/internal/config"
	"service_api/internal/log"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := log.LoggerFromContext(ctx)
	ctx = log.ContextWithLogger(ctx, l)
	configFile := flag.String("config", "configs/config.yml", "Path to config file.")
	flag.Parse()
	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		l.Fatal("fail to initialize config", zap.Error(err))
	}
	s, err := apiserver.New(ctx, cfg)
	if err != nil {
		l.Fatal("fail to initialize service api", zap.Error(err))
	}
	if err = s.Start(ctx); err != nil {
		l.Fatal("fail to start service api", zap.Error(err))
	}

	// Gracefully shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		l.Info("shutdown application")
		cancel()
	}()
}
