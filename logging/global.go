package logging

import "log/slog"

var globalLogLevel string
var globalLogLevelChanger = DefaultGlobalLevelChanger

func SetGlobalLevel(level string) error {
	if err := globalLogLevelChanger(level); err != nil {
		return err
	}
	globalLogLevel = level
	return nil
}

func GetGlobalLevel() string {
	return globalLogLevel
}

func DefaultGlobalLevelChanger(level string) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return err
	}
	slog.SetLogLoggerLevel(l)
	return nil
}

// program using concrete logging framework is expected to call this
func SetGlobalLevelChanger(changer func(level string) error) {
	globalLogLevelChanger = changer
}
