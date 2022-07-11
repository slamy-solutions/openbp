package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/olebedev/emitter"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	amqp "github.com/rabbitmq/amqp091-go"

	amqpTools "github.com/slamy-solutions/open-erp/modules/system/libs/go/rabbitmq"

	lambdaGRPC "github.com/slamy-solutions/open-erp/modules/native/services/lambda/entrypoint/src/grpc/native_lambda"
)

type LambdaEntrypointServer struct {
	lambdaGRPC.UnimplementedLambdaEntrypointServiceServer

	lambdaManagerClient lambdaGRPC.LambdaManagerServiceClient
	amqpChanel          *amqp.Channel
	responseEmitter     *emitter.Emitter
	replyQueue          string
}

const TASKS_EXCHANGE = "native_lambda_entrypoint_input"

func NewLambdaEntrypointServer(lambdaManagerClient lambdaGRPC.LambdaManagerServiceClient, amqpChanel *amqp.Channel, responseEmitter *emitter.Emitter, replyQueue string) (*LambdaEntrypointServer, error) {
	err := amqpChanel.ExchangeDeclare(TASKS_EXCHANGE, "direct", true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, err
	}

	return &LambdaEntrypointServer{
		lambdaManagerClient: lambdaManagerClient,
		amqpChanel:          amqpChanel,
		responseEmitter:     responseEmitter,
		replyQueue:          replyQueue,
	}, nil
}

func (s *LambdaEntrypointServer) Call(ctx context.Context, in *lambdaGRPC.CallLambdaRequest) (*lambdaGRPC.CallLambdaResponse, error) {
	lambdaResponse, err := s.lambdaManagerClient.Get(ctx, &lambdaGRPC.GetLambdaRequest{Namespace: in.Namespace, Uuid: in.Lambda})
	if err != nil {
		return nil, err
	}
	lambda := lambdaResponse.Lambda

	body := lambdaGRPC.AMQPLambdaTaskRequest{
		Lambda: lambda,
		Data:   in.Data,
	}
	bytesBody, err := proto.Marshal(&body)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	err = s.amqpChanel.Publish(TASKS_EXCHANGE, lambda.Runtime, false, false, amqp.Publishing{
		Headers:      amqpTools.InjectTelemetryAMQPHeaders(ctx),
		DeliveryMode: 2,
		Body:         bytesBody,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &lambdaGRPC.CallLambdaResponse{}, status.Errorf(grpccodes.OK, "")
}
func (s *LambdaEntrypointServer) Execute(ctx context.Context, in *lambdaGRPC.ExecuteLambdaRequest) (*lambdaGRPC.ExecuteLambdaResponse, error) {
	lambdaResponse, err := s.lambdaManagerClient.Get(ctx, &lambdaGRPC.GetLambdaRequest{Namespace: in.Namespace, Uuid: in.Lambda})
	if err != nil {
		return nil, err
	}
	lambda := lambdaResponse.Lambda

	requestUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	correlationId := requestUUID.String()

	body := lambdaGRPC.AMQPLambdaTaskRequest{
		Lambda: lambda,
		Data:   in.Data,
	}
	bytesBody, err := proto.Marshal(&body)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	err = s.amqpChanel.Publish(TASKS_EXCHANGE, lambda.Runtime, false, false, amqp.Publishing{
		Headers:       amqpTools.InjectTelemetryAMQPHeaders(ctx),
		CorrelationId: correlationId,
		ReplyTo:       s.replyQueue,
		DeliveryMode:  2,
		Body:          bytesBody,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	select {
	case response := <-s.responseEmitter.Once(correlationId):
		delivery, ok := response.Args[0].(amqp.Delivery)
		if !ok {
			return nil, status.Error(grpccodes.Internal, "Wrong delivery format")
		}
		var decodedResponse lambdaGRPC.AMQPLambdaTaskResponse
		err := proto.Unmarshal(delivery.Body, &decodedResponse)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Lambda responsed with bad data format. "+err.Error())
		}

		if decodedResponse.StatusCode != uint32(grpccodes.OK) {
			return nil, status.Error(grpccodes.Code(decodedResponse.StatusCode), "Bad response from the lambda: "+decodedResponse.Message)
		}

		return &lambdaGRPC.ExecuteLambdaResponse{Result: decodedResponse.Data}, status.Error(grpccodes.OK, "")
	case <-time.After(time.Duration(uint64(time.Millisecond) * in.Timeout)):
		return nil, status.Error(grpccodes.DeadlineExceeded, "Execution Timeout")
	}
}
