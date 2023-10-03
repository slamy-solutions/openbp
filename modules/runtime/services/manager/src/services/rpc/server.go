package rpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	grpcRPC "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/rpc"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	systemNATS "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManagerRPCServer struct {
	grpcRPC.UnimplementedRPCServiceServer

	systemStub *system.SystemStub
	logger     *slog.Logger
}

func NewManagerRPCServer(logger *slog.Logger, systemStub *system.SystemStub) *ManagerRPCServer {
	return &ManagerRPCServer{
		systemStub: systemStub,
		logger:     logger,
	}
}

func (s *ManagerRPCServer) Call(ctx context.Context, in *grpcRPC.CallRequest) (*grpcRPC.CallResponse, error) {
	if in.Timeout == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "timeout must be greater than 0")
	}

	logger := s.logger.With(
		slog.String("endpoint", "Call"),
		slog.String("namespace", in.Namespace),
		slog.String("runtimeName", in.RuntimeName),
		slog.String("methodName", in.MethodName),
	)

	requestMessage := grpcRPC.RPCRequestMesasge{
		Data: in.Payload,
	}
	requestBytes, err := proto.Marshal(&requestMessage)
	if err != nil {
		err := errors.Join(errors.New("failed to marshal request"), err)
		logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	msg := &nats.Msg{
		Subject: fmt.Sprintf("runtime.rpc.jsexecutor.inbound.%s.%s.%s", in.Namespace, in.RuntimeName, in.MethodName),
		Data:    requestBytes,
	}
	systemNATS.InjectTelemetryContext(ctx, msg)

	response, err := s.systemStub.Nats.RequestMsg(msg, time.Duration(in.Timeout)*time.Millisecond)
	if err != nil {
		if err == nats.ErrTimeout {
			return nil, status.Errorf(codes.DeadlineExceeded, "timeout")
		}

		err := errors.Join(errors.New("failed to call method"), err)
		logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var responseMessage grpcRPC.RPCResponseMessage
	err = proto.Unmarshal(response.Data, &responseMessage)
	if err != nil {
		err := errors.Join(errors.New("failed to unmarshal response"), err)
		logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &grpcRPC.CallResponse{
		Response:     responseMessage.Response,
		Error:        responseMessage.Error,
		ErrorMessage: responseMessage.ErrorMessage,
	}, nil
}
