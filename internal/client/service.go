package client

import (
	"context"
	"errors"
	"sync"
	"time"

	commonApi "github.com/qrinef/process-app/pkg/api"
	commonApiEntities "github.com/qrinef/process-app/pkg/api/entities"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	once   sync.Once

	configSvc   configService
	externalSvc externalService

	itemChannel     chan commonApiEntities.Item
	batchItemsLimit int
	intervalLimit   time.Duration
}

func (s *Service) Start(ctx context.Context) {
	ticker := time.NewTicker(s.intervalLimit)

	defer s.once.Do(func() {
		close(s.itemChannel)
		s.logger.Info("batch process aborted")
	})

	for {
		select {
		case <-ticker.C:
			err := s.process(ctx, s.handler())
			if err != nil && errors.Is(err, commonApi.ErrBlocked) {
				s.logger.Error("unable to batch process", zap.Error(err))

				/*
				 * Implementation of returning batch items back to the channel
				 * if the process ended unsuccessfully
				 */
			}
			if err != nil && errors.Is(err, ErrMissingBatch) {
				s.logger.Info("waiting new items...")
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *Service) AddItem(item commonApiEntities.Item) error {
	select {
	case s.itemChannel <- item:
		return nil
	default:
		return ErrUnableAddNewItem
	}
}

func (s *Service) handler() *commonApiEntities.Batch {
	filledBatch := make(commonApiEntities.Batch, 0, s.batchItemsLimit)

	for len(filledBatch) < s.batchItemsLimit {
		select {
		case num, ok := <-s.itemChannel:
			if !ok && len(filledBatch) != 0 {
				return &filledBatch
			}
			if !ok {
				return nil
			}

			filledBatch = append(filledBatch, num)

			if len(filledBatch) == s.batchItemsLimit {
				return &filledBatch
			}
		default:
			if len(filledBatch) != 0 {
				return &filledBatch
			}

			return nil
		}
	}

	return &filledBatch
}

func (s *Service) process(ctx context.Context, batch *commonApiEntities.Batch) error {
	if batch == nil {
		return ErrMissingBatch
	}

	err := s.externalSvc.Process(ctx, *batch)
	if err != nil {
		return err
	}

	s.logger.Info("batch processed successfully", zap.Int("items", len(*batch)))

	return nil
}

func NewService(logger *zap.Logger,
	configSvc configService,
	externalSvc externalService,
) *Service {
	service := &Service{
		logger: logger.Named("client.service"),

		configSvc:   configSvc,
		externalSvc: externalSvc,

		itemChannel: make(chan commonApiEntities.Item, configSvc.GetClientBufferSize()),
	}

	batchItemsLimit, intervalLimit := service.externalSvc.GetLimits()

	service.batchItemsLimit = int(batchItemsLimit)
	service.intervalLimit = intervalLimit

	return service
}
