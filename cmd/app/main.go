package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qrinef/process-app/internal/client"
	"github.com/qrinef/process-app/internal/config"

	commonApi "github.com/qrinef/process-app/pkg/api"
	commonApiEntities "github.com/qrinef/process-app/pkg/api/entities"
	commonConfig "github.com/qrinef/process-app/pkg/config"
	commonLogger "github.com/qrinef/process-app/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg := &config.Config{}
	err := commonConfig.Prepare(appCfg)
	if err != nil {
		log.Fatalf("unable to prepare application config: %v", err)
	}

	loggerSvc, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatalf("unable to create application logger: %v", err)
	}
	loggerEntry := loggerSvc.NewLoggerEntry("main")

	externalSvc := commonApi.NewService(appCfg)

	clientSvc := client.NewService(loggerEntry, appCfg, externalSvc)
	go clientSvc.Start(ctx)

	for i := 0; i < 15; i++ {
		errAddItem := clientSvc.AddItem(commonApiEntities.Item{})
		if errAddItem != nil {
			loggerEntry.Error("unable to add new item", zap.Error(errAddItem))
			break
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	loggerEntry.Warn("shutdown application")
	cancelCtxFunc()

	err = loggerEntry.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		log.Fatalf("unable to sync logger: %v", err)
	}

	log.Print("stopped")
}
