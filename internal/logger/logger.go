package logger

import (
	"fmt"
	"os"

	"github.com/major1ink/simple-notification-telegram/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(
	fileName string,
) (*zap.Logger, error) {
	cfgLogger := zap.NewDevelopmentConfig()

	cfgLogger.DisableCaller = false
	cfgLogger.DisableStacktrace = true

	cfgLogger.EncoderConfig.EncodeTime = customTimeEncoder
	cfgLogger.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	cfgLogger.Level.SetLevel(getLevel(config.AppConfig().Logger.GetLevel()))

	if fileName == "error.log" {
		cfgLogger.Level.SetLevel(zapcore.ErrorLevel)
		cfgLogger.DisableStacktrace = false
	}

	if config.AppConfig().Logger.GetLogMode() == "stdout" && fileName != "error.log" {
		logger, err := cfgLogger.Build()
		if err != nil {
			return nil, err
		}
		return logger, nil
	}

	logFile := fmt.Sprintf("%s/%s", config.AppConfig().Logger.GetLogDir(), fileName)
	err := os.MkdirAll(config.AppConfig().Logger.GetLogDir(), os.ModePerm)
	if err != nil {
		return nil, err
	}

	if config.AppConfig().Logger.GetRewriteLog() {
		if _, err := os.Stat(logFile); err == nil {
			err := os.Remove(logFile)
			if err != nil {
				return nil, err
			}
		}
	}
	cfgLogger.OutputPaths = []string{logFile}

	logger, err := cfgLogger.Build(zap.AddCaller())
	if err != nil {
		return nil, err
	}
	return logger, nil
}
