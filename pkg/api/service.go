package api

import (
	"context"
	"errors"
	"time"

	"github.com/qrinef/process-app/pkg/api/entities"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

type Service struct {
	configSvc configService
}

func (s *Service) GetLimits() (uint64, time.Duration) {
	return s.configSvc.GetExternalBatchItemsLimit(), s.configSvc.GetExternalIntervalLimit()
}

func (s *Service) Process(_ context.Context, batch entities.Batch) error {
	if len(batch) > int(s.configSvc.GetExternalBatchItemsLimit()) {
		return ErrBlocked
	}

	return nil
}

func NewService(configSvc configService) *Service {
	return &Service{
		configSvc,
	}
}
