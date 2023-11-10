package settings

import (
	"log/slog"

	settingsGRPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"
	backendModels "github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
)

type OneCData struct {
	RemoteURL string `bson:"remoteURL" json:"remoteURL"`
	Token     string `bson:"token" json:"token"`
}

type Settings struct {
	Namespace   string                    `bson:"namespace" json:"namespace"`
	BackendType backendModels.BackendType `bson:"backendType" json:"backendType"`

	OneCData *OneCData `bson:"oneCData" json:"oneCData"`
}

func (s *Settings) ToGRPC() *settingsGRPC.Settings {
	result := &settingsGRPC.Settings{
		Namespace:   s.Namespace,
		BackendType: backendModels.BackendTypeToGRPC(s.BackendType),
	}

	if s.OneCData != nil {
		result.Backend = &settingsGRPC.Settings_OneC{
			OneC: &settingsGRPC.OneCBackendSettings{
				RemoteURL: s.OneCData.RemoteURL,
				Token:     s.OneCData.Token,
			},
		}
	}

	if s.BackendType == backendModels.BackendTypeNative {
		result.Backend = &settingsGRPC.Settings_Native{
			Native: &settingsGRPC.NativeBackendSettings{},
		}
	}

	return result
}

func SettingsFromGRPC(settings *settingsGRPC.Settings) *Settings {
	result := &Settings{
		Namespace:   settings.Namespace,
		BackendType: backendModels.BackendTypeFromGRPC(settings.BackendType),
	}

	if settings.Backend != nil {
		switch backend := settings.Backend.(type) {
		case *settingsGRPC.Settings_OneC:
			result.OneCData = &OneCData{
				RemoteURL: backend.OneC.RemoteURL,
				Token:     backend.OneC.Token,
			}
		case *settingsGRPC.Settings_Native:
		}
	}

	return result
}

func (s *Settings) ToSlogAttr(groupKey string) slog.Attr {
	var data []any
	if groupKey == "" {
		groupKey = "settings"
	}

	data = append(data, slog.String("namespace", s.Namespace))
	data = append(data, slog.String("backendType", string(s.BackendType)))

	if s.OneCData != nil {
		data = append(data, slog.Group("oneCData", slog.String("remoteURL", s.OneCData.RemoteURL)))
	}

	return slog.Group(groupKey, data...)
}
