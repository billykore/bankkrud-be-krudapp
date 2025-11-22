package log

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure sets up the global logger configuration.
func Configure(env string) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" || env == "test" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// Add stack trace hook
	log.Logger = log.Logger.Hook(StackHook{})
	// Add output hook
	logfile, err := os.OpenFile("./.logs/app.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open log file")
	}
	log.Logger = log.Output(logfile)
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
