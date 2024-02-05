package logger

type LoggerConfig struct {
	// -------------------
	// Application configs
	// -------------------
	MinimalLogsLevel string `envconfig:"LOGGER_LEVEL" default:"debug"`
	StackTraceEnable bool   `envconfig:"LOGGER_STACKTRACE_ENABLE" default:"false"`
}

func (c *LoggerConfig) GetMinimalLogLevel() string {
	return c.MinimalLogsLevel
}

func (c *LoggerConfig) IsStacktraceEnabled() bool {
	return c.StackTraceEnable
}
