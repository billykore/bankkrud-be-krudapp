package log

import (
	"fmt"
	"os"

	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/log"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

// ZapLogger provides logging functions.
type ZapLogger struct {
	logger  *zap.Logger
	usecase string
}

// NewZap return the singleton instance of Logger.
func NewZap() *ZapLogger {
	logger := newZap()
	return &ZapLogger{logger: logger}
}

func newZap() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.InfoLevel)
	logger := zap.New(core, zap.AddStacktrace(zap.PanicLevel))
	return logger
}

// Logger returns the underlying zap.Logger instance for the Logger.
func (zl *ZapLogger) Logger() *zap.Logger {
	return zl.logger
}

func (zl *ZapLogger) Usecase(usecase string) log.Logger {
	zl.usecase = usecase
	return zl
}

// Errorf log.
func (zl *ZapLogger) Errorf(format string, args ...any) {
	zl.logger.Error(
		fmt.Sprintf(format, args...),
		zap.String("usecase", zl.usecase),
	)
}

// Fatalf logs an error, then shutdown the domain.
func (zl *ZapLogger) Fatalf(format string, args ...any) {
	zl.logger.Fatal(
		fmt.Sprintf(format, args...),
		zap.String("usecase", zl.usecase),
	)
}

// Infof log.
func (zl *ZapLogger) Infof(format string, args ...any) {
	zl.logger.Info(
		fmt.Sprintf(format, args...),
		zap.String("usecase", zl.usecase),
	)
}
