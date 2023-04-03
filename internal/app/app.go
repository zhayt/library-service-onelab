package app

import (
	"fmt"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/service"
	"github.com/zhayt/user-storage-service/internal/storage"
	"github.com/zhayt/user-storage-service/internal/transport/http"
	"github.com/zhayt/user-storage-service/internal/transport/http/handler"
	logger2 "github.com/zhayt/user-storage-service/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run() error {
	logger := logger2.NewLogger()
	cfg, err := config.New()
	if err != nil {
		return err
	}

	repo := storage.NewStorage()

	serv := service.NewUserService(repo)

	hand := handler.NewHandler(logger, serv)

	server := http.NewServer(cfg, hand.InitRoute())

	go func() {
		logger.LogInfo.Printf("Start server at port %v -> http://localhost%v", cfg.HTTP.Port, cfg.HTTP.Port)
		server.Start()
	}()

	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		logger.LogInfo.Printf("signal accepted: %v", s.String())
	case err = <-server.Notify():
		logger.LogInfo.Printf("server closing: %v", err)
	}

	if err = server.Shutdown(); err != nil {
		return fmt.Errorf("error while shutting down server: %s", err)
	}

	return nil
}
