package main

import (
	"context"
	"fmt"
	"github.com/zhayt/user-storage-service/config"
	_ "github.com/zhayt/user-storage-service/docs"
	"github.com/zhayt/user-storage-service/internal/service"
	"github.com/zhayt/user-storage-service/internal/storage"
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

//	@title			OneLab HomeWork API
//	@version		1.0
//	@description	API service for User Storage.
//	@description	Where they can create, retrieve, update, delete books.
//  @description	And can rent these books
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8000
//	@BasePath	/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	repo := storage.NewStorage(ctx, &wg, l, cfg)

	// service
	serv := service.NewService(l, repo)

	// middleware
	mid := middleware.NewJWTAuth(cfg)

	// handler
	hand := handler.NewHandler(l, serv, mid)

	// server
	server := http.NewServer(cfg, hand, mid)

	l.Info("Start server")
	server.Start()

	// grace full shutdown
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

	cancel()
	wg.Wait()

	return nil
}
