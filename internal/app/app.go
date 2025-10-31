package app

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/major1ink/simple-notification-telegram/internal/config"
	"github.com/major1ink/simple-notification-telegram/internal/logger"
	"github.com/major1ink/simple-notification-telegram/pkg/closer"
)

type App struct {
	diContainer *diContainer
	logger      *zap.Logger
	closer      *closer.Closer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		if err := a.runAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("assembled consumer crashed: %v", err)
		}
	}()

	select {
	case <-a.closer.Done():
		a.logger.Info("Shutdown signal received")
		return a.gracefulShutdown(ctx)
	case err := <-errCh:
		a.logger.Error("Component crashed, shutting down", zap.Error(err))
		return a.gracefulShutdown(ctx)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initCloser,
		a.initLogger,
		a.initDI,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	return config.Load(configPath)
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	a.diContainer.SetLogger(a.logger)
	a.diContainer.SetCloser(a.closer)
	return nil
}

func (a *App) initLogger(ctx context.Context) error {

	l, err := logger.NewLog("simple-notification-telegram.log")
	if err != nil {
		return err
	}
	a.logger = l
	a.closer.SetLogger(a.logger)

	a.closer.AddNamed("logger", func(ctx context.Context) error {
		if err := a.logger.Sync(); err != nil && err.Error() != "sync /dev/stderr: invalid argument" {
			return err
		}
		return nil
	})

	return nil

}

func (a *App) initCloser(_ context.Context) error {
	a.closer = closer.NewWithLogger(zap.NewNop(), syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func (a *App) runAssembledConsumer(ctx context.Context) error {
	a.logger.Info("🚀 Kafka assembled consumer running")

	err := a.diContainer.AssembleConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) gracefulShutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := a.closer.CloseAll(ctx); err != nil {
		a.logger.Error("❌ Ошибка при завершении работы", zap.Error(err))
		return err
	}
	return nil
}
