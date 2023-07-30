package balena

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevicesServer struct {
	balena.UnimplementedBalenaDevicesServiceServer

	systemStub *system.SystemStub
	logger     *logrus.Entry
}

func NewDevicesServer(logger *logrus.Entry, systemStub *system.SystemStub) *DevicesServer {
	return &DevicesServer{
		logger:     logger,
		systemStub: systemStub,
	}
}

func (s *DevicesServer) Bind(ctx context.Context, in *balena.BindDeviceRequest) (*balena.BindDeviceResponse, error) {
	balenaDeviceUUID, err := primitive.ObjectIDFromHex(in.BalenaDeviceUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "device not found. Balena device UUID has bad format")
	}

	balenaUUID, err := primitive.ObjectIDFromHex(in.DeviceUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "device not found. Device UUID has bad format")
	}

	collection := getBalenaDeviceCollection(s.systemStub)
	var updatedDevice BalenaDeviceInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": balenaDeviceUUID},
		bson.M{
			"$set": bson.M{
				"deviceNamespace": in.DeviceNamespace,
				"deviceUUID":      balenaUUID,
			},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedDevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		err = errors.Join(errors.New("error while binding device in the database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.BindDeviceResponse{
		Device: updatedDevice.ToGRPCDevice(),
	}, status.Error(codes.OK, "")
}
func (s *DevicesServer) UnBind(ctx context.Context, in *balena.UnBindDeviceRequest) (*balena.UnBindDeviceResponse, error) {
	balenaDeviceUUID, err := primitive.ObjectIDFromHex(in.BalenaDeviceUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "device not found. Balena device UUID has bad format")
	}

	collection := getBalenaDeviceCollection(s.systemStub)
	var updatedDevice BalenaDeviceInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": balenaDeviceUUID},
		bson.M{
			"$set": bson.M{
				"deviceNamespace": nil,
				"deviceUUID":      nil,
			},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedDevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		err = errors.Join(errors.New("error while unbinding device in the database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.UnBindDeviceResponse{
		Device: updatedDevice.ToGRPCDevice(),
	}, status.Error(codes.OK, "")
}
func (s *DevicesServer) ListInNamespace(in *balena.ListDevicesInNamespaceRequest, out balena.BalenaDevicesService_ListInNamespaceServer) error {
	ctx := out.Context()

	filter := bson.M{"balenaServerNamespace": in.BalenaServersNamespace}
	switch in.BindingFilter {
	case balena.ListDevicesInNamespaceRequest_ONLY_BINDED:
		filter["bindedDeviceUUID"] = bson.M{"$ne": nil}
	case balena.ListDevicesInNamespaceRequest_ONLY_UNBINDED:
		filter["bindedDeviceUUID"] = nil
	default:
	}

	collection := getBalenaDeviceCollection(s.systemStub)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		err = errors.Join(errors.New("error while finding devices in namespace. Failed to initialize database cursor"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var device BalenaDeviceInMongo
		if err = cur.Decode(&device); err != nil {
			err = errors.Join(errors.New("error while finding devices in namespace. Failed to decode device from database"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&balena.ListDevicesInNamespaceResponse{
			Device: device.ToGRPCDevice(),
		})
		if err != nil {
			err = errors.Join(errors.New("error while finding devices in namespace. Failed to send back founded device as result"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	if err = cur.Err(); err != nil {
		err = errors.Join(errors.New("error while finding devices in namespace. Database cursor error"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.OK, "")
}
func (s *DevicesServer) CountInNamespace(ctx context.Context, in *balena.CountDevicesInNamespaceRequest) (*balena.CountDevicesInNamespaceResponse, error) {
	filter := bson.M{"balenaServerNamespace": in.BalenaServersNamespace}
	switch in.BindingFilter {
	case balena.CountDevicesInNamespaceRequest_ONLY_BINDED:
		filter["bindedDeviceUUID"] = bson.M{"$ne": nil}
	case balena.CountDevicesInNamespaceRequest_ONLY_UNBINDED:
		filter["bindedDeviceUUID"] = nil
	default:
	}

	collection := getBalenaDeviceCollection(s.systemStub)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		err = errors.Join(errors.New("error while counting devices in namespace: database error"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.CountDevicesInNamespaceResponse{Count: uint64(count)}, status.Error(codes.OK, "")
}
