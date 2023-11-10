package services

import (
	"context"
	"errors"
	"log/slog"

	settings "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"
	repository "github.com/slamy-solutions/openbp/modules/crm/services/core/src/settings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SettingsService struct {
	settings.UnimplementedSettingsServiceServer

	settings repository.SettingsRepository
	logger   *slog.Logger
}

func NewSettingsServer(settings repository.SettingsRepository, logger *slog.Logger) *SettingsService {
	return &SettingsService{
		settings: settings,
		logger:   logger,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context, in *settings.GetSettingsRequest) (*settings.GetSettingsResponse, error) {
	settingsData, err := s.settings.Get(ctx, in.Namespace, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get settings"), err)
		s.logger.With(slog.String("route", "GetSettings"), slog.String("namespace", in.Namespace)).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get settings: %s", err.Error())
	}

	return &settings.GetSettingsResponse{
		Settings: settingsData.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *SettingsService) SetSettings(ctx context.Context, in *settings.SetSettingsRequest) (*settings.SetSettingsResponse, error) {
	settingsToSet := repository.SettingsFromGRPC(in.Settings)

	err := s.settings.Set(ctx, settingsToSet)
	if err != nil {
		err := errors.Join(errors.New("failed to set settings"), err)
		s.logger.With(slog.String("route", "SetSettings"), slog.String("namespace", in.Settings.Namespace)).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to set settings: %s", err.Error())
	}

	s.logger.With(slog.String("route", "SetSettings"), slog.String("namespace", in.Settings.Namespace), settingsToSet.ToSlogAttr("newSettings")).Info("settings was updated set")
	return &settings.SetSettingsResponse{}, status.Error(codes.OK, "")
}
