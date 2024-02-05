package client

import (
	"context"
	"errors"
	"time"

	commonApiEntities "github.com/qrinef/process-app/pkg/api/entities"
)

var (
	ErrMissingBatch     = errors.New("missing batch")
	ErrUnableAddNewItem = errors.New("unable to add new item")
)

type configService interface {
	GetClientBufferSize() uint64
	SetClientBufferSize(value uint64)
}

type externalService interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch commonApiEntities.Batch) error
}
