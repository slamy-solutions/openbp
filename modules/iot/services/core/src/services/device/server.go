package device

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	deviceGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceServer struct {
	deviceGRPC.UnimplementedDeviceServiceServer

	logger  *logrus.Entry
	randGen *rand.Rand

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

func NewDeviceServer(ctx context.Context, logger *logrus.Entry, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*DeviceServer, error) {
	err := CreateIndexesForNamespace(ctx, systemStub, "")
	if err != nil {
		return nil, errors.New("failed to create indexes for global namespace: " + err.Error())
	}

	return &DeviceServer{
		systemStub: systemStub,
		nativeStub: nativeStub,

		logger:  logger,
		randGen: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func DeviceCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	if namespaceName == "" {
		return systemStub.DB.Database("openbp_global").Collection("iot_core_device")
	} else {
		db := systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespaceName))
		return db.Collection("iot_core_device")
	}
}

var managementIdLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *DeviceServer) randManagementId() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = managementIdLetters[rand.Intn(len(managementIdLetters))]
	}
	return string(b)
}

func (s *DeviceServer) Create(ctx context.Context, in *deviceGRPC.CreateRequest) (*deviceGRPC.CreateResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	// Create identity for the device.
	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       in.Namespace,
		Name:            "",
		InitiallyActive: true,
		Managed: &identity.CreateIdentityRequest_Service{
			Service: &identity.ServiceManagedData{
				Service:      "iot_core_device",
				Reason:       "Manage identity for device",
				ManagementId: s.randManagementId(),
			},
		},
	})
	if err != nil {
		err = errors.New("error while creating identity for new device: " + err.Error())
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	creationTime := time.Now().UTC()
	insertData := DeviceInMongo{
		Name:        in.Name,
		Description: in.Description,
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}

	insertResult, err := collection.UpdateOne(ctx, bson.M{"name": in.Name}, bson.M{"$setOnInsert": insertData}, options.MergeUpdateOptions().SetUpsert(true))
	if err != nil {
		_, deleteErr := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
			Namespace: in.Namespace,
			Uuid:      identityCreateResponse.Identity.Uuid,
		})
		if deleteErr != nil {
			s.logger.Warn("error while deleteting identity after the device failed to create: " + deleteErr.Error())
		} else {
			s.logger.Info("identity was successfully deleted after the device failed to create")
		}

		err = errors.New("error while inserting new device to the database: " + err.Error())
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if insertResult.UpsertedCount == 0 {
		return nil, status.Error(codes.AlreadyExists, "Device with same name already exist")
	}

	insertData.UUID = insertResult.UpsertedID.(primitive.ObjectID)
	s.logger.WithFields(logrus.Fields{
		"iot_core_device_uuid": insertData.UUID.Hex(),
		"namespace":            in.Namespace,
	}).Info("Successfully created new device")
	return &deviceGRPC.CreateResponse{Device: insertData.ToGRPCDevice(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *DeviceServer) Get(ctx context.Context, in *deviceGRPC.GetRequest) (*deviceGRPC.GetResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	deviceUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Device not found. Bad device UUID format.")
	}

	var foundedDevice DeviceInMongo
	err = collection.FindOne(ctx, bson.M{"_id": deviceUUID}).Decode(&foundedDevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Device with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Device not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while getting device from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.Uuid,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deviceGRPC.GetResponse{Device: foundedDevice.ToGRPCDevice(in.Namespace)}, status.Error(codes.OK, "")
}

func (s *DeviceServer) GetByIdentity(ctx context.Context, in *deviceGRPC.GetByIdentityRequest) (*deviceGRPC.GetByIdentityResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	var foundedDevice DeviceInMongo
	err := collection.FindOne(ctx, bson.M{"identity": in.Identity}).Decode(&foundedDevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Device with this identity not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Device not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while getting device by identity from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_identity": in.Identity,
			"namespace":                in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deviceGRPC.GetByIdentityResponse{Device: foundedDevice.ToGRPCDevice(in.Namespace)}, status.Error(codes.OK, "")
}

func (s *DeviceServer) Exists(ctx context.Context, in *deviceGRPC.ExistsRequest) (*deviceGRPC.ExistsResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	deviceUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &deviceGRPC.ExistsResponse{Exists: false}, status.Error(codes.OK, "Device doesnt exist. Bad device UUID format.")
	}

	count, err := collection.CountDocuments(ctx, bson.M{"_id": deviceUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &deviceGRPC.ExistsResponse{Exists: false}, status.Error(codes.OK, "Device doesnt exist. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while checking if device exists in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.Uuid,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deviceGRPC.ExistsResponse{Exists: count > 0}, status.Error(codes.OK, "")
}

func (s *DeviceServer) Count(ctx context.Context, in *deviceGRPC.CountRequest) (*deviceGRPC.CountResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &deviceGRPC.CountResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist.")
			}
		}

		err = errors.New("error while counting devices in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deviceGRPC.CountResponse{Count: uint64(count)}, status.Error(codes.OK, "")
}

func (s *DeviceServer) List(in *deviceGRPC.ListRequest, out deviceGRPC.DeviceService_ListServer) error {
	ctx := out.Context()
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	findOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	if in.Limit != 0 {
		findOptions.SetLimit(int64(in.Limit))
	}
	if in.Skip != 0 {
		findOptions.SetSkip(int64(in.Skip))
	}
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		err = errors.New("error while opening find stream to list devices in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var foundedDevice DeviceInMongo
		err := cursor.Decode(&foundedDevice)

		if err != nil {
			err = errors.New("error while listing: failed to decode device from database find stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&deviceGRPC.ListResponse{
			Device: foundedDevice.ToGRPCDevice(in.Namespace),
		})
		if err != nil {
			err = errors.New("error while listing: failed to send device information to outgoing stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	return status.Error(codes.OK, "")
}

func (s *DeviceServer) Update(ctx context.Context, in *deviceGRPC.UpdateRequest) (*deviceGRPC.UpdateResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	deviceUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Device doesnt exist. Bad device UUID format.")
	}

	var deviceAfterUpdate DeviceInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": deviceUUID}, bson.M{
		"$set": bson.M{
			"description": in.Description,
		},
		"$inc": bson.M{
			"version": 1,
		},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&deviceAfterUpdate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Device with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Device not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while updating device in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.Uuid,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deviceGRPC.UpdateResponse{Device: deviceAfterUpdate.ToGRPCDevice(in.Namespace)}, status.Error(codes.OK, "")
}

func (s *DeviceServer) Delete(ctx context.Context, in *deviceGRPC.DeleteRequest) (*deviceGRPC.DeleteResponse, error) {
	collection := DeviceCollectionByNamespace(s.systemStub, in.Namespace)

	deviceUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &deviceGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Device doesnt exist. Bad device UUID format.")
	}

	// Find device to get its Identity because it must also be deleted.
	var foundedDevice DeviceInMongo
	err = collection.FindOne(ctx, bson.M{"_id": deviceUUID}).Decode(&foundedDevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &deviceGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Device with this UUID doesnt exist")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &deviceGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Device doesnt exist. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while getting device before the deleletion from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.Uuid,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Actually deleting device and its identity
	deviceDeleteResponse, err := collection.DeleteOne(ctx, bson.M{"_id": deviceUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &deviceGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Device doesnt exist. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while deleting device from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.Uuid,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if deviceDeleteResponse.DeletedCount > 0 {
		_, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
			Namespace: in.Namespace,
			Uuid:      foundedDevice.Identity,
		})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"iot_core_device_uuid": in.Uuid,
				"namespace":            in.Namespace,
			}).Warn("error while deleting device: error while deleting identity for the device: " + err.Error())
		}
	}

	return &deviceGRPC.DeleteResponse{Existed: deviceDeleteResponse.DeletedCount > 0}, status.Error(codes.OK, "")
}
