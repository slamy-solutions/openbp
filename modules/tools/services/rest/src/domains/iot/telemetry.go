package iot

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type TelemetryRouter struct {
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewTelemetryRouter(logger *logrus.Entry, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub) *TelemetryRouter {
	return &TelemetryRouter{
		logger:     logger,
		systemStub: systemStub,
		nativeStub: nativeStub,
		iotStub:    iotStub,
	}
}

type TelemetryListenRequest struct {
	Namespace          string `form:"namespace" binding:"lte=64,regexp=^[A-Za-z0-9]+$"`
	Devices            string `form:"devices" binding:"required,lte=2048,regexp=^[A-Za-z0-9,]+$"`
	ListenBasicMetrics bool   `form:"listenBasicMetrics"`
	ListenEvents       bool   `form:"listenEvents"`
	ListenLogs         bool   `form:"listenLogs"`
}

var telemetryWebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *TelemetryRouter) Listen(ctx *gin.Context) {
	var requestData TelemetryListenRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	resources := strings.Split(requestData.Devices, ",")
	for i := range resources {
		resources[i] = "iot.core.device." + resources[i]
	}
	if len(resources) > 100 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "You can listen to at most 100 devices at same time"})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"device.namespace": requestData.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            resources,
			Actions:              []string{"iot.core.telemetry.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	devices := strings.Split(requestData.Devices, ",")
	subjects := make([]string, 0, len(devices))
	if requestData.ListenBasicMetrics {
		for _, device := range devices {
			subjects = append(subjects, "iot.core.telemetry.basic_metrics."+requestData.Namespace+"."+device)
		}
	}
	if requestData.ListenEvents {
		for _, device := range devices {
			subjects = append(subjects, "iot.core.telemetry.events."+requestData.Namespace+"."+device)
		}
	}
	if requestData.ListenLogs {
		for _, device := range devices {
			subjects = append(subjects, "iot.core.telemetry.logs."+requestData.Namespace+"."+device)
		}
	}

	js, _ := r.systemStub.Nats.JetStream()
	streamInfo, err := js.AddStream(&nats.StreamConfig{
		Subjects:    subjects,
		MaxMsgs:     100,
		Storage:     nats.MemoryStorage,
		Discard:     nats.DiscardOld,
		Description: "tools_rest listen for device telemetry updates for connected client",
	})
	if err != nil {
		err := errors.New("failed to open NATS stream to listen for telemetry events: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer js.DeleteStream(streamInfo.Config.Name)

	subscription, err := js.SubscribeSync("*", nats.BindStream(streamInfo.Config.Name))
	if err != nil {
		err := errors.New("failed to subscribe to NATS stream to listen for telemetry events: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer subscription.Unsubscribe()

	ws, err := telemetryWebsocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Warn("failed to upgrade websocket connection while streaming device events: " + err.Error())
		return
	}
	defer ws.Close()

	for {
		msg, err := subscription.NextMsg(time.Minute * 15)
		if err != nil {
			if err == nats.ErrTimeout {
				return
			}

			err := errors.New("failed to open NATS stream to listen for telemetry events: " + err.Error())
			logger.Error(err.Error())

			ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""), time.Now().Add(time.Second*10))
			return
		}

		err = msg.Ack()
		if err != nil {
			logger.Error("failed to ack NATS device telemetry message : " + err.Error())

			ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""), time.Now().Add(time.Second*10))
			return
		}

		messageToSend := struct {
			MessageType string      `json:"messageType"`
			Message     interface{} `json:"message"`
		}{}
		if strings.HasPrefix(msg.Subject, "iot.core.telemetry.basic_metrics") {
			var basicMetrics telemetry.BasicMetrics
			if err := proto.Unmarshal(msg.Data, &basicMetrics); err != nil {
				logger.Error("failed to unmarshall devices basic metrics telemetry object: " + err.Error())
				continue
			}
			messageToSend.MessageType = "basicMetrics"
			messageToSend.Message = formatedTelemetryBasicMetricsFromGRPC(&basicMetrics)
		} else if strings.HasPrefix(msg.Subject, "iot.core.telemetry.logs") {
			var logEntry telemetry.LogEntry
			if err := proto.Unmarshal(msg.Data, &logEntry); err != nil {
				logger.Error("failed to unmarshall devices log entry telemetry object: " + err.Error())
				continue
			}
			messageToSend.MessageType = "log"
			messageToSend.Message = formatedTelemetryLogEntryFromGRPC(&logEntry)
		} else if strings.HasPrefix(msg.Subject, "iot.core.telemetry.events") {
			var event telemetry.Event
			if err := proto.Unmarshal(msg.Data, &event); err != nil {
				logger.Error("failed to unmarshall devices event telemetry object: " + err.Error())
				continue
			}
			messageToSend.MessageType = "event"
			messageToSend.Message = formatedTelemetryEventFromGRPC(&event)
		} else {
			logger.Warn("bad message subject while receiving devices telemetry: " + msg.Subject)
			continue
		}

		err = ws.WriteJSON(messageToSend)
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*10))
					return
				}
			}

			logger.Error("websocket error while sending telemetry object: " + err.Error())
			ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""), time.Now().Add(time.Second*10))
			return
		}
	}
}
