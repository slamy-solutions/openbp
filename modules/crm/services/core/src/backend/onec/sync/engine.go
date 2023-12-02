package sync

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"

	systemSync "sync"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
	onecPerformer "github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/performer"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/settings"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type SyncEngine interface {
	Start()
	Stop()

	SyncAllNow(ctx context.Context) error
	SyncNow(ctx context.Context, namespace string) error
	ListSyncEvents(namespace string, skip int, limit int) ([]SyncEvent, int, error)
}

type syncEngine struct {
	logger     *slog.Logger
	systemStub *system.SystemStub
	nativeStub *native.NativeStub

	settingsRepository settings.SettingsRepository

	syncContext      context.Context
	syncCancel       context.CancelFunc
	syncWorkerWaiter systemSync.WaitGroup
}

func NewSyncEngine(logger *slog.Logger, systemStub *system.SystemStub, nativeStub *native.NativeStub, settingsRepository settings.SettingsRepository) SyncEngine {
	return &syncEngine{
		logger:             logger,
		systemStub:         systemStub,
		nativeStub:         nativeStub,
		settingsRepository: settingsRepository,

		syncContext:      nil,
		syncCancel:       nil,
		syncWorkerWaiter: systemSync.WaitGroup{},
	}
}

func (e *syncEngine) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	e.syncContext = ctx
	e.syncCancel = cancel
	e.syncWorkerWaiter.Add(1)
	go e.syncWorker()
}

func (e *syncEngine) Stop() {
	e.syncCancel()
	e.syncWorkerWaiter.Wait()
}

func (e *syncEngine) syncWorker() {
	e.logger.Info("Sync worker started")
	defer e.syncWorkerWaiter.Done()
	for {
		select {
		case <-e.syncContext.Done():
			break
		case <-time.After(5 * time.Minute):
			e.SyncAllNow(context.Background())
		}
	}
}

func (e *syncEngine) SyncAllNow(ctx context.Context) error {
	namespacesStream, err := e.nativeStub.Services.Namespace.GetAll(ctx, &namespace.GetAllNamespacesRequest{
		UseCache: true,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to get namespaces"), err)
		e.logger.Error(err.Error())
		return err
	}

	var namespaces []*namespace.Namespace
	for {
		r, err := namespacesStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err := errors.Join(errors.New("failed to get namespaces"), err)
			e.logger.Error(err.Error())
			return err
		}

		namespaces = append(namespaces, r.Namespace)
	}

	for _, n := range namespaces {
		err := e.SyncNow(ctx, n.Name)
		if err != nil {
			e.logger.Error(err.Error())
		}
	}
	err = e.SyncNow(ctx, "")
	if err != nil {
		e.logger.Error(err.Error())
	}

	return nil
}

func (e *syncEngine) SyncNow(ctx context.Context, namespace string) error {
	logBuffer := new(strings.Builder)

	jsonLogHandler := slog.NewJSONHandler(logBuffer, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(jsonLogHandler)

	submitLog := func(finalError error) {
		logs := logBuffer.String()
		success := finalError == nil
		errorMessage := ""
		if !success {
			errorMessage = finalError.Error()
		}
		syncEvent := SyncEvent{
			Namespace: namespace,
			Timestamp: time.Now(),
			Status:    success,
			Error:     errorMessage,
			Log:       logs,
		}

		err := addSyncLog(ctx, e.systemStub, syncEvent)
		if err != nil {
			err := errors.Join(errors.New("failed to add sync log"), err)
			e.logger.Error(err.Error())
		}
	}

	settings, err := e.settingsRepository.Get(ctx, namespace, true)
	if err != nil {
		err := errors.Join(errors.New("failed to get settings"), err)
		e.logger.Error(err.Error())
		return err
	}

	if settings.BackendType == models.BackendType1C {
		c := connector.NewOneCConnector(settings.OneCData.RemoteURL, settings.OneCData.Token)

		err := onecPerformer.SyncWithOneCServer(ctx, logger, namespace, e.systemStub, e.nativeStub, c)
		if err != nil {
			submitLog(err)
			return err
		}
	}

	submitLog(nil)
	return nil
}

func (e *syncEngine) ListSyncEvents(namespace string, skip int, limit int) ([]SyncEvent, int, error) {
	return getSyncLog(e.syncContext, e.systemStub, namespace, skip, limit)
}
