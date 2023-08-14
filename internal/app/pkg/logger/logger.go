package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

type Logger struct {
	l *zap.Logger
	// https://pkg.go.dev/go.uber.org/zap#example-AtomicLevel
	al *zap.AtomicLevel
}

func New(level Level, opts ...Option) *Logger {

	al := zap.NewAtomicLevelAt(level)
	cfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		NameKey:        "transfer",
		CallerKey:      "caller",
		LevelKey:       "level",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), al,
	)

	return &Logger{l: zap.New(core, opts...), al: &al}
}

// SetLevel 动态更改日志级别
// 对于使用 NewTee 创建的 Logger 无效，因为 NewTee 本意是根据不同日志级别
// 创建的多个 zap.Core，不应该通过 SetLevel 将多个 zap.Core 日志级别统一
func (l *Logger) SetLevel(level Level) {
	if l.al != nil {
		l.al.SetLevel(level)
	}
}

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.l.Sugar().Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.l.Sugar().Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.l.Sugar().Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.l.Sugar().Errorf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.l.Sugar().Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.l.Sugar().Fatalf(template, args...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

var std = New(InfoLevel, AddCaller(), AddCallerSkip(2))

func Default() *Logger         { return std }
func ReplaceDefault(l *Logger) { std = l }

func SetLevel(level Level) { std.SetLevel(level) }

func Debug(msg string, fields ...Field) { std.Debug(msg, fields...) }
func Info(msg string, fields ...Field)  { std.Info(msg, fields...) }
func Warn(msg string, fields ...Field)  { std.Warn(msg, fields...) }
func Error(msg string, fields ...Field) { std.Error(msg, fields...) }
func Panic(msg string, fields ...Field) { std.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field) { std.Fatal(msg, fields...) }

func Debugf(template string, args ...interface{}) { std.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { std.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { std.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { std.Errorf(template, args...) }
func Panicf(template string, args ...interface{}) { std.Panicf(template, args...) }
func Fatalf(template string, args ...interface{}) { std.Fatalf(template, args...) }

func Sync() error { return std.Sync() }
