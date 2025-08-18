package log

// Logger interface defines the methods for logging.
type Logger interface {
	Usecase(usecase string) Logger
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}
