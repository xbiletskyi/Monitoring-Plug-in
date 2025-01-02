package zaplogger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger global effective zaplogger by uber zap
var Logger *zap.SugaredLogger

// InitLogger initializes predefined zap zaplogger
func InitLogger() {
	customLogger := buildCustomLogger()
	Logger = customLogger.Sugar()
}

// customTimeEncoder formats the time as "MM/DD HH:mm:ss"
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01/02 15:04:05"))
}

// customCallerEncoder formats the caller as "(file:line)"
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("(" + caller.TrimmedPath() + ")")
}

// buildCustomLogger creates a zaplogger with the custom format
func buildCustomLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "zaplogger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		EncodeTime:     customTimeEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
}
