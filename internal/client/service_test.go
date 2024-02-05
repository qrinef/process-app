package client

import (
	"context"
	"testing"
	"time"

	"github.com/qrinef/process-app/internal/client/mock"
	"github.com/qrinef/process-app/internal/config"

	commonApi "github.com/qrinef/process-app/pkg/api"
	commonApiEntities "github.com/qrinef/process-app/pkg/api/entities"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type ServiceSuite struct {
	suite.Suite

	logger      *zap.Logger
	externalSvc *client.MockExternal

	configSvc configService
}

func (s *ServiceSuite) SetupTest() {
	s.logger = zaptest.NewLogger(s.T())
	s.externalSvc = client.NewMockExternal(s.T())
	s.configSvc = &config.Config{}
}

func (s *ServiceSuite) TestAddItemSuccess() {
	s.configSvc.SetClientBufferSize(10)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(30, time.Second)

	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	err := clientSvc.AddItem(commonApiEntities.Item{})
	s.NoError(err)
}

func (s *ServiceSuite) TestAddItemError() {
	s.configSvc.SetClientBufferSize(1)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(30, time.Second)

	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	err := clientSvc.AddItem(commonApiEntities.Item{})
	s.NoError(err)

	err = clientSvc.AddItem(commonApiEntities.Item{})
	s.ErrorIs(err, ErrUnableAddNewItem)
}

func (s *ServiceSuite) TestProcessSuccess() {
	s.configSvc.SetClientBufferSize(10)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(10, time.Second)

	s.externalSvc.EXPECT().
		Process(mock.Anything, mock.Anything).
		Return(nil)

	ctx := context.Background()
	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	err := clientSvc.AddItem(commonApiEntities.Item{})
	s.NoError(err)

	err = clientSvc.process(ctx, clientSvc.handler())
	s.NoError(err)
}

func (s *ServiceSuite) TestProcessError() {
	s.configSvc.SetClientBufferSize(10)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(1, time.Second)

	s.externalSvc.EXPECT().
		Process(mock.Anything, mock.Anything).
		Return(commonApi.ErrBlocked)

	ctx := context.Background()
	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	err := clientSvc.AddItem(commonApiEntities.Item{})
	s.NoError(err)

	err = clientSvc.process(ctx, clientSvc.handler())
	s.ErrorIs(err, commonApi.ErrBlocked)
}

func (s *ServiceSuite) TestProcessWaitItems() {
	s.configSvc.SetClientBufferSize(10)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(1, time.Second)

	ctx := context.Background()
	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	err := clientSvc.process(ctx, nil)
	s.ErrorIs(err, ErrMissingBatch)
}

func (s *ServiceSuite) TestServiceCancelCtx() {
	s.configSvc.SetClientBufferSize(10)

	s.externalSvc.EXPECT().
		GetLimits().
		Return(10, time.Second)

	ctx, cancelCtxFunc := context.WithCancel(context.Background())
	clientSvc := NewService(s.logger, s.configSvc, s.externalSvc)

	cancelCtxFunc()
	clientSvc.Start(ctx)
}

func TestExternalServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
