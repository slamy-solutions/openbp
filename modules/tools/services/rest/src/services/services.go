package services

import (
	"context"
	"os"
)

type ServicesConnectionHandler struct {
	System *SystemConnectionHandler
	Native *NativeConnectionHandler
}

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConnectToServices(ctx context.Context) (*ServicesConnectionHandler, error) {
	system, err := ConnectToSystemServices(ctx)
	if err != nil {
		return nil, err
	}

	native, err := ConnectToNativeServices(ctx)
	if err != nil {
		system.Shutdown(ctx)
		return nil, err
	}

	return &ServicesConnectionHandler{
		System: system,
		Native: native,
	}, nil
}

func (h *ServicesConnectionHandler) Shutdown(ctx context.Context) {
	h.Native.Shutdown(ctx)
	h.System.Shutdown(ctx)
}
