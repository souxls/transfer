package app

import (
	"transfer/internal/app/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func InitLogger() error {
	defer logger.Sync()
	logger.SetLevel(zapcore.Level(viper.GetViper().GetInt("log.Level")))
	return nil
}
