package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvironmentNameTag = "environment"
)

type Service struct {
	configSvc configService

	defaultLogger *zap.Logger
}

func (s *Service) NewLoggerEntry(named string) *zap.Logger {
	var cores = []zapcore.Core{
		s.defaultLogger.Core(),
	}

	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	l = l.Named(named).With(zap.String(EnvironmentNameTag, s.configSvc.GetEnvironmentName()))

	return l
}

func NewService(configSvc configService) (*Service, error) {
	logsLevel := new(zapcore.Level)
	err := logsLevel.Set(configSvc.GetMinimalLogLevel())
	if err != nil {
		return nil, err
	}

	lCfg := zap.NewProductionConfig()
	lCfg.OutputPaths = []string{"stdout"}
	lCfg.ErrorOutputPaths = []string{"stdout"}
	lCfg.Level = zap.NewAtomicLevelAt(*logsLevel)
	lCfg.DisableStacktrace = !configSvc.IsStacktraceEnabled()

	if configSvc.IsDebug() {
		lCfg.Level.SetLevel(zap.DebugLevel)
	}

	defaultLogger, err := lCfg.Build()
	if err != nil {
		return nil, err
	}

	return &Service{
		configSvc:     configSvc,
		defaultLogger: defaultLogger,
	}, nil
}
