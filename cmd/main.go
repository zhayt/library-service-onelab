package main

import (
	"fmt"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/service"
	"github.com/zhayt/user-storage-service/internal/storage"
	"github.com/zhayt/user-storage-service/internal/storage/postgres"
	"github.com/zhayt/user-storage-service/internal/transport/http"
	"github.com/zhayt/user-storage-service/internal/transport/http/handler"
	"github.com/zhayt/user-storage-service/internal/transport/http/middleware"
	"github.com/zhayt/user-storage-service/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		return
	}
}

func run() error {
	// config
	var once sync.Once

	once.Do(config.PrepareEnv)

	cfg, err := config.New()
	if err != nil {
		return err
	}

	// logger
	l, err := logger.Init(cfg)
	if err != nil {
		return fmt.Errorf("cannot init logger: %w", err)
	}
	defer func(l *zap.Logger) {
		err = l.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(l)

	// storage
	db, err := postgres.NewConnectionPool("pgx", cfg.DBConnectionURL)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := storage.NewStorage(l, db)

	// service
	serv := service.NewService(l, repo)

	// middleware
	mid := middleware.NewJWTAuth(cfg)

	// handler
	hand := handler.NewHandler(l, serv, mid)

	// server
	server := http.NewServer(cfg, hand, mid)

	l.Info("app started")

	go func() {
		l.Info("Start server", zap.String("host", cfg.AppHost), zap.String("port", cfg.AppPort))
		server.Start()
	}()

	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err = <-server.Notify():
		l.Info("server closing", zap.Error(err))
	}

	if err = server.Shutdown(); err != nil {
		return fmt.Errorf("error while shutting down server: %s", err)
	}

	return nil
}
