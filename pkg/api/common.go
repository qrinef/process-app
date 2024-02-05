package api

import "time"

type configService interface {
	GetExternalBatchItemsLimit() uint64
	GetExternalIntervalLimit() time.Duration
}
