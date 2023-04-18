package logger

import (
	"github.com/zhayt/user-storage-service/config"
	"go.uber.org/zap"
)

func Init(cfg *config.Config) (*zap.Logger, error) {
	switch cfg.Level {
	case "dev":
		return zap.NewDevelopment()
	default:
		return zap.NewProduction()
	}
}
