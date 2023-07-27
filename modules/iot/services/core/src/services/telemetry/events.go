package telemetry

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	telemetryGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RaiseNewDeviceEventEvent(systemStub *system.SystemStub, event *telemetryGRPC.Event) error {
	subject := fmt.Sprintf("iot.core.telemetry.events.%s.%s", event.DeviceNamespace, event.DeviceUUID)
	messageData, err := proto.Marshal(event)
	if err != nil {
		return errors.New("error while marshalling event: " + err.Error())
	}
	err = systemStub.Nats.Publish(subject, messageData)
	if err != nil {
		return errors.New("error while publishing event: " + err.Error())
	}

	return RaiseNewBHeartBeatEvent(systemStub, event.DeviceNamespace, event.DeviceUUID, event.Timestamp.AsTime())
}

func RaiseNewLogEntryEvent(systemStub *system.SystemStub, logEntry *telemetryGRPC.LogEntry) error {
	subject := fmt.Sprintf("iot.core.telemetry.logs.%s.%s", logEntry.DeviceNamespace, logEntry.DeviceUUID)
	messageData, err := proto.Marshal(logEntry)
	if err != nil {
		return errors.New("error while marshalling log entry: " + err.Error())
	}
	err = systemStub.Nats.Publish(subject, messageData)
	if err != nil {
		return errors.New("error while publishing log entry: " + err.Error())
	}

	return RaiseNewBHeartBeatEvent(systemStub, logEntry.DeviceNamespace, logEntry.DeviceUUID, logEntry.Timestamp.AsTime())
}

func RaiseNewBasicMetricsEvent(systemStub *system.SystemStub, metrics *telemetryGRPC.BasicMetrics) error {
	subject := fmt.Sprintf("iot.core.telemetry.basic_metrics.%s.%s", metrics.DeviceNamespace, metrics.DeviceUUID)
	messageData, err := proto.Marshal(metrics)
	if err != nil {
		return errors.New("error while marshalling basic metrics: " + err.Error())
	}
	err = systemStub.Nats.Publish(subject, messageData)
	if err != nil {
		return errors.New("error while publishing basic metrics: " + err.Error())
	}

	return RaiseNewBHeartBeatEvent(systemStub, metrics.DeviceNamespace, metrics.DeviceUUID, metrics.Timestamp.AsTime())
}

func RaiseNewBHeartBeatEvent(systemStub *system.SystemStub, deviceNamespace string, deviceUUID string, timestamp time.Time) error {
	subject := fmt.Sprintf("iot.core.telemetry.heartbeat.%s.%s", deviceNamespace, deviceUUID)
	messageData, err := proto.Marshal(&telemetryGRPC.HeartBeat{
		DeviceNamespace: deviceNamespace,
		DeviceUUID:      deviceUUID,
		Timestamp:       timestamppb.New(timestamp),
	})
	if err != nil {
		return errors.New("error while marshalling heartbeat: " + err.Error())
	}
	err = systemStub.Nats.Publish(subject, messageData)
	if err != nil {
		return errors.New("error while publishing heartbeat: " + err.Error())
	}

	return nil
}
