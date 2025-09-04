package log

import (
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure sets up the global logger configuration.
func Configure(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// Add stack trace hook
	log.Logger = log.Logger.Hook(StackHook{})
}

type StackHook struct{}

func (h StackHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level >= zerolog.PanicLevel { // Add stack trace for Error and Fatal levels
		e.Str("stack", string(debug.Stack()))
	}
}
