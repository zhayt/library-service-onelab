package logger

import (
	"github.com/zhayt/transaction-service/config"
	"go.uber.org/zap"
)

func Init(cfg *config.Config) (*zap.Logger, error) {
	if cfg.AppMode == "dev" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}
