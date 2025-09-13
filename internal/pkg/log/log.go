package log

import (
	"context"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure sets up the global logger configuration.
func Configure(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" || env == "test" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// Add stack trace hook
	log.Logger = log.Logger.Hook(StackHook{})
}

type StackHook struct{}

func (h StackHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level >= zerolog.PanicLevel {
		e.Str("stack", string(debug.Stack()))
	}
}

// WithContext returns a new zerolog.Logger bound to the provided context and an "usecase" field.
// It ensures:
//   - ctx is non-nil (falls back to context.Background()).
//   - usecase is non-empty (falls back to "unknown") for consistent log filtering.
func WithContext(ctx context.Context, usecase string) zerolog.Logger {
	if ctx == nil {
		ctx = context.Background()
	}
	if usecase == "" {
		usecase = "unknown"
	}
	return log.With().Ctx(ctx).Str("usecase", usecase).Logger()
}
