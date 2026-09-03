package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger(level string) {
	// 解析日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建 core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// 创建 logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = zapLogger.Sugar()
}

// 便捷方法
func Info(msg string, keysAndValues ...interface{}) {
	args := make([]interface{}, 0, 1+len(keysAndValues))
	args = append(args, msg)
	args = append(args, keysAndValues...)
	Log.Info(args...)
}
func Debug(msg string, keysAndValues ...interface{}) {
	args := make([]interface{}, 0, 1+len(keysAndValues))
	args = append(args, msg)
	args = append(args, keysAndValues...)
	Log.Debug(args...)
}
func Warn(msg string, keysAndValues ...interface{}) {
	args := make([]interface{}, 0, 1+len(keysAndValues))
	args = append(args, msg)
	args = append(args, keysAndValues...)
	Log.Warn(args...)
}
func Error(msg string, keysAndValues ...interface{}) {
	args := make([]interface{}, 0, 1+len(keysAndValues))
	args = append(args, msg)
	args = append(args, keysAndValues...)
	Log.Error(args...)
}
func Fatal(msg string, keysAndValues ...interface{}) {
	args := make([]interface{}, 0, 1+len(keysAndValues))
	args = append(args, msg)
	args = append(args, keysAndValues...)
	Log.Fatal(args...)
}
func Infof(template string, args ...interface{})      { Log.Infof(template, args...) }
func Errorf(template string, args ...interface{})     { Log.Errorf(template, args...) }
func Warnf(template string, args ...interface{})      { Log.Warnf(template, args...) }
func Debugf(template string, args ...interface{})     { Log.Debugf(template, args...) }
