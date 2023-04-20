package main

import (
	"fmt"
	"github.com/zhayt/transaction-service/config"
	_ "github.com/zhayt/transaction-service/docs"
	"github.com/zhayt/transaction-service/logger"
	"github.com/zhayt/transaction-service/service"
	"github.com/zhayt/transaction-service/storage"
	"github.com/zhayt/transaction-service/transport/http"
	"github.com/zhayt/transaction-service/transport/http/handler"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//	@title			OneLab HomeWork API
//	@version		1.0
//	@description	Transaction microservice.
//	@termsOfService	http://swagger.io/terms/

// @host		localhost:8081
// @BasePath	/local/db
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
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
	repo, err := storage.NewStorage(cfg, l)
	if err != nil {
		return err
	}

	// service

	serv := service.NewService(repo, l)
	// handler

	handl := handler.NewHandler(serv, l)

	// server
	server := http.NewServer(handl, cfg)

	l.Info("Start server")
	server.StartHTTPServer()

	// grace full shutdown
	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err = <-server.Notify:
		l.Info("server closing", zap.Error(err))
	}

	if err = server.GracefullyShutdown(); err != nil {
		return fmt.Errorf("error while shutting down server: %s", err)
	}

	return nil
}
