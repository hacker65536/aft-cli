package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLog *zap.Logger

func InitializeLogger(lvl zapcore.Level) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(lvl)
	config.DisableStacktrace = true
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig
	var err error
	ZapLog, err = config.Build()
	if err != nil {
		panic(err)
	}
}
