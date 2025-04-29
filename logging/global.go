package logging

var globalLogLevel string
var globalLogLevelChanger func(level string) error = func(level string) error { return nil }

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

// program using concrete logging framework is expected to call this
func SetGlobalLevelChanger(changer func(level string) error) {
	globalLogLevelChanger = changer
}
