package logging

import "log/slog"

// Create with [New], [NewRoot]
type Logger struct {
	*slog.Logger
}

func (l Logger) Fatal(msg string, args ...any) {
	l.Logger.Error(msg, args...)
	panic(msg)
}

func (l Logger) With(args ...any) Logger {
	c := l.Logger.With(args...)
	return Logger{c}
}

func New(name string) Logger {
	return Logger{slog.Default().With("logger", name)}
}

func NewRoot() Logger {
	return Logger{slog.Default()}
}
