package logger

type configService interface {
	GetEnvironmentName() string
	IsDebug() bool

	GetMinimalLogLevel() string
	IsStacktraceEnabled() bool
}
