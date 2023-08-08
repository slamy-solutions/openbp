package fleet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	fleetGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	deviceServer "github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FleetServer struct {
	fleetGRPC.UnimplementedFleetServiceServer

	logger *logrus.Entry

	systemStub *system.SystemStub
	nativeStub *native.NativeStub

	device *deviceServer.DeviceServer
}

func NewFleetServer(ctx context.Context, logger *logrus.Entry, systemStub *system.SystemStub, nativeStub *native.NativeStub, device *deviceServer.DeviceServer) (*FleetServer, error) {
	err := CreateIndexesForNamespace(ctx, systemStub, "")
	if err != nil {
		return nil, errors.New("failed to create indexes for global namespace: " + err.Error())
	}

	return &FleetServer{
		systemStub: systemStub,
		nativeStub: nativeStub,

		logger: logger,

		device: device,
	}, nil
}

func FleetCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	if namespaceName == "" {
		return systemStub.DB.Database("openbp_global").Collection("iot_core_fleet")
	} else {
		db := systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespaceName))
		return db.Collection("iot_core_fleet")
	}
}

func FleetDevicesCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	if namespaceName == "" {
		return systemStub.DB.Database("openbp_global").Collection("iot_core_fleetdevices")
	} else {
		db := systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespaceName))
		return db.Collection("iot_core_fleetdevices")
	}
}

func (s *FleetServer) Create(ctx context.Context, in *fleetGRPC.CreateRequest) (*fleetGRPC.CreateResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	creationTime := time.Now().UTC()
	insertData := FleetInMongo{
		Name:        in.Name,
		Description: in.Description,
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}

	updateResult, err := collection.UpdateOne(ctx, bson.M{"name": in.Name}, bson.M{"$setOnInsert": insertData}, options.MergeUpdateOptions().SetUpsert(true))
	if err != nil {
		err = errors.New("error while inserting new fleet to the database: " + err.Error())
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if updateResult.UpsertedCount == 0 {
		return nil, status.Error(codes.AlreadyExists, "fleet with same name already exist")
	}

	insertData.UUID = updateResult.UpsertedID.(primitive.ObjectID)
	s.logger.WithFields(logrus.Fields{
		"iot_core_fleet_uuid": insertData.UUID.Hex(),
		"namespace":           in.Namespace,
	}).Info("Successfully created new fleet")
	return &fleetGRPC.CreateResponse{Fleet: insertData.ToGRPCFleet(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *FleetServer) Get(ctx context.Context, in *fleetGRPC.GetRequest) (*fleetGRPC.GetResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Fleet not found. Bad fleet UUID format.")
	}

	var foundedFleet FleetInMongo
	err = collection.FindOne(ctx, bson.M{"_id": fleetUUID}).Decode(&foundedFleet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Fleet with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Fleet not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while getting fleet from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.GetResponse{Fleet: foundedFleet.ToGRPCFleet(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *FleetServer) Exists(ctx context.Context, in *fleetGRPC.ExistsRequest) (*fleetGRPC.ExistsResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &fleetGRPC.ExistsResponse{Exists: false}, status.Error(codes.OK, "Fleet doesnt exist. Bad fleet UUID format.")
	}

	count, err := collection.CountDocuments(ctx, bson.M{"_id": fleetUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.ExistsResponse{Exists: false}, status.Error(codes.OK, "Fleet doesnt exist. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while checking if fleet exists in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.ExistsResponse{Exists: count > 0}, status.Error(codes.OK, "")
}
func (s *FleetServer) Update(ctx context.Context, in *fleetGRPC.UpdateRequest) (*fleetGRPC.UpdateResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Fleet doesnt exist. Bad fleet UUID format.")
	}

	var fleetAfterUpdate FleetInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": fleetUUID}, bson.M{
		"$set": bson.M{
			"description": in.Description,
		},
		"$inc": bson.M{
			"version": 1,
		},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&fleetAfterUpdate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Fleet with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Fleet not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while updating fleet in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.UpdateResponse{Fleet: fleetAfterUpdate.ToGRPCFleet(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *FleetServer) Delete(ctx context.Context, in *fleetGRPC.DeleteRequest) (*fleetGRPC.DeleteResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &fleetGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Fleet doesnt exist. Bad fleet UUID format.")
	}

	deleteResponse, err := collection.DeleteOne(ctx, bson.M{"_id": fleetUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Fleet doesnt exist. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while deleting fleet from the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.DeleteResponse{Existed: deleteResponse.DeletedCount > 0}, status.Error(codes.OK, "")
}
func (s *FleetServer) Count(ctx context.Context, in *fleetGRPC.CountRequest) (*fleetGRPC.CountResponse, error) {
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.CountResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist.")
			}
		}

		err = errors.New("error while counting fleets in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.CountResponse{Count: uint64(count)}, status.Error(codes.OK, "")
}
func (s *FleetServer) List(in *fleetGRPC.ListRequest, out fleetGRPC.FleetService_ListServer) error {
	ctx := out.Context()
	collection := FleetCollectionByNamespace(s.systemStub, in.Namespace)

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

		err = errors.New("error while opening find stream to list fleets in the database: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	for cursor.Next(ctx) {
		var foundedFleet FleetInMongo
		err := cursor.Decode(&foundedFleet)

		if err != nil {
			err = errors.New("error while listing: failed to decode fleet from database find stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&fleetGRPC.ListResponse{
			Fleet: foundedFleet.ToGRPCFleet(in.Namespace),
		})
		if err != nil {
			err = errors.New("error while listing: failed to send fleet information to outgoing stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	return status.Error(codes.OK, "")
}

func makeListDecivesWithoutFleetAggregation(skip int64, limit int64) bson.A {
	aggregation := bson.A{
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "created", Value: 1}}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "iot_core_fleet"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "deviceUUID"},
				{Key: "as", Value: "fleetInfo"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$fleetInfo"}, {Key: "preserveNullAndEmptyArrays", Value: true}}},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "fleetInfo", Value: nil},
			}},
		},
	}
	if skip != 0 {
		aggregation = append(aggregation, bson.D{
			{Key: "$skip", Value: skip},
		})
	}
	if limit != 0 {
		aggregation = append(aggregation, bson.D{
			{Key: "$limit", Value: limit},
		})
	}
	return aggregation
}

func makeListDevicesInFleetAggregation(fleetUUID primitive.ObjectID, skip int64, limit int64) bson.A {
	aggregation := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "fleetUUID", Value: fleetUUID},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "added", Value: 1}}},
		},
	}

	if skip != 0 {
		aggregation = append(aggregation, bson.D{
			{Key: "$skip", Value: skip},
		})
	}
	if limit != 0 {
		aggregation = append(aggregation, bson.D{
			{Key: "$limit", Value: limit},
		})
	}

	aggregation = append(aggregation, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "iot_core_device"},
			{Key: "localField", Value: "deviceUUID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "deviceInfo"},
		}},
	})

	return aggregation
}

func (s *FleetServer) CountDevicesWithoutFleet(ctx context.Context, namespace string) (*fleetGRPC.CountDevicesResponse, error) {
	aggregation := bson.A{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "iot_core_fleet"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "deviceUUID"},
				{Key: "as", Value: "fleetInfo"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$fleetInfo"}, {Key: "preserveNullAndEmptyArrays", Value: true}}},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "fleetInfo", Value: nil},
			}},
		},
		bson.D{
			{Key: "$count", Value: "count"},
		},
	}

	collection := deviceServer.DeviceCollectionByNamespace(s.systemStub, namespace)
	cursor, err := collection.Aggregate(ctx, aggregation)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.CountDevicesResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		err = errors.New("error while opening find stream to count devices that are not in a fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	cursor.Next(ctx)

	var response struct{ count int }
	err = cursor.Decode(&response)
	if err != nil {
		err = errors.New("error while counting devices without fleet: error while decoding database stream: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.CountDevicesResponse{Count: uint64(response.count)}, status.Error(codes.OK, "")
}

func (s *FleetServer) CountDevices(ctx context.Context, in *fleetGRPC.CountDevicesRequest) (*fleetGRPC.CountDevicesResponse, error) {
	collection := FleetDevicesCollectionByNamespace(s.systemStub, in.Namespace)

	if in.Uuid == "" {
		return s.CountDevicesWithoutFleet(ctx, in.Namespace)
	}

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &fleetGRPC.CountDevicesResponse{Count: 0}, status.Error(codes.OK, "Fleet doesnt exist. Bad fleet UUID format.")
	}

	count, err := collection.CountDocuments(ctx, bson.M{"fleetUUID": fleetUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.CountDevicesResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist.")
			}
		}

		err = errors.New("error while counting devices in the fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fleetGRPC.CountDevicesResponse{Count: uint64(count)}, status.Error(codes.OK, "")
}

func (s *FleetServer) ListDevicesWithoutFleet(ctx context.Context, in *fleetGRPC.ListDevicesRequest, out fleetGRPC.FleetService_ListDevicesServer) error {
	collection := deviceServer.DeviceCollectionByNamespace(s.systemStub, in.Namespace)
	aggregation := makeListDecivesWithoutFleetAggregation(int64(in.Skip), int64(in.Limit))

	cursor, err := collection.Aggregate(ctx, aggregation, options.Aggregate())
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		err = errors.New("error while opening find stream to list devices that are not in a fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"namespace": in.Namespace,
		}).Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	var foundedDevice deviceServer.DeviceInMongo
	for cursor.Next(ctx) {
		err := cursor.Decode(&foundedDevice)
		if err != nil {
			err = errors.New("error while listing devices without fleet: error while decoding database stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&fleetGRPC.ListDevicesResponse{
			Device: &fleetGRPC.ListDevicesResponse_FleetDevice{
				Device: foundedDevice.ToGRPCDevice(in.Namespace),
				Added:  nil,
			},
		})
		if err != nil {
			err = errors.New("error while listing devices without fleet: error while sending the response: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"iot_core_device_uuid": foundedDevice.UUID.Hex(),
				"namespace":            in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	return status.Error(codes.OK, "")
}

func (s *FleetServer) ListDevices(in *fleetGRPC.ListDevicesRequest, out fleetGRPC.FleetService_ListDevicesServer) error {
	ctx := out.Context()
	collection := FleetDevicesCollectionByNamespace(s.systemStub, in.Namespace)

	// Listing the devices that are not in the namespace
	if in.Uuid == "" {
		return s.ListDevicesWithoutFleet(ctx, in, out)
	}

	fleetUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return status.Error(codes.OK, "Fleet doesnt exist. Bad fleet UUID format.")
	}

	aggregation := makeListDevicesInFleetAggregation(fleetUUID, int64(in.Skip), int64(in.Limit))
	cursor, err := collection.Aggregate(ctx, aggregation, options.Aggregate())
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		err = errors.New("error while opening find stream to list devices in the fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	sendedDevices := 0
	for cursor.Next(ctx) {
		var entry struct {
			DeviceUUID primitive.ObjectID `bson:"deviceUUID"`
			FleetUUID  primitive.ObjectID `bson:"fleetUUID"`
			Added      time.Time          `bson:"added"`

			DeviceInfo []*deviceServer.DeviceInMongo `bson:"deviceInfo"`
		}

		err := cursor.Decode(&entry)
		if err != nil {
			err = errors.New("error while listing devices in the fleet: error while decoding database stream: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"iot_core_fleet_uuid": in.Uuid,
				"namespace":           in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		// Fix possible race condition while adding device to the fleet.
		// Its not possible to reliable validate if device exist during adding it to the fleet.
		// This scenario is extremelly rare
		if len(entry.DeviceInfo) == 0 {
			cursor.Close(ctx)
			_, err := collection.DeleteOne(ctx, bson.M{"fleetUUID": fleetUUID, "deviceUUID": entry.DeviceUUID})
			if err != nil {
				err = errors.New("error while listing devices in the fleet: error on clearing operation: " + err.Error())
				s.logger.WithFields(logrus.Fields{
					"iot_core_fleet_uuid":  in.Uuid,
					"iot_core_device_uuid": entry.DeviceUUID.Hex(),
					"namespace":            in.Namespace,
				}).Error(err.Error())
				return status.Error(codes.Internal, err.Error())
			}
			return s.ListDevices(&fleetGRPC.ListDevicesRequest{
				Namespace: in.Namespace,
				Uuid:      in.Uuid,
				Skip:      in.Skip + uint64(sendedDevices),
				Limit:     in.Limit,
			}, out)
		}

		err = out.Send(&fleetGRPC.ListDevicesResponse{
			Device: &fleetGRPC.ListDevicesResponse_FleetDevice{
				Device: entry.DeviceInfo[0].ToGRPCDevice(in.Namespace),
				Added:  timestamppb.New(entry.Added),
			},
		})
		if err != nil {
			err = errors.New("error while listing devices in the fleet: error while sending the response: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"iot_core_fleet_uuid":  in.Uuid,
				"iot_core_device_uuid": entry.DeviceUUID.Hex(),
				"namespace":            in.Namespace,
			}).Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		sendedDevices += 1
	}

	if cursor.Err() != nil {
		err = errors.New("error while listing devices in the fleet: error while getting next value from database cursor: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.Uuid,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.OK, "")
}
func (s *FleetServer) AddDevice(ctx context.Context, in *fleetGRPC.AddDeviceRequest) (*fleetGRPC.AddDeviceResponse, error) {
	fleetDevicesCollection := FleetDevicesCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.FleetUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Fleet doesnt exist. Bad fleet UUID format.")
	}
	deviceUUID, err := primitive.ObjectIDFromHex(in.DeviceUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Device doesnt exist. Bad device UUID format.")
	}

	// Those operations (check if fleet and device exist) are not atomic
	// Race conditions may occure between check and ading device to the fleet
	// This is expected. They will be removed during the first list operation.

	fleetExistsResponse, err := s.Exists(ctx, &fleetGRPC.ExistsRequest{Namespace: in.Namespace, Uuid: in.FleetUUID})
	if err != nil {
		err = errors.New("error while checking if fleet exists: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_fleet_uuid": in.FleetUUID,
			"namespace":           in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !fleetExistsResponse.Exists {
		return nil, status.Error(codes.NotFound, "Fleet doesnt exist.")
	}

	deviceExistsResponse, err := s.device.Exists(ctx, &device.ExistsRequest{Namespace: in.Namespace, Uuid: in.DeviceUUID})
	if err != nil {
		err = errors.New("error while checking if device exists: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.DeviceUUID,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !deviceExistsResponse.Exists {
		return nil, status.Error(codes.NotFound, "Device doesnt exist.")
	}

	_, err = fleetDevicesCollection.InsertOne(ctx, bson.M{
		"fleetUUID":  fleetUUID,
		"deviceUUID": deviceUUID,
		"added":      time.Now().UTC(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(codes.AlreadyExists, "Device is already assigned to a fleet")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Fleet and device not found. Namespace doesnt exist.")
			}
		}

		err = errors.New("error while adding device to the fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.DeviceUUID,
			"iot_core_fleet_uuid":  in.FleetUUID,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.WithFields(logrus.Fields{
		"iot_core_device_uuid": in.DeviceUUID,
		"iot_core_fleet_uuid":  in.FleetUUID,
		"namespace":            in.Namespace,
	}).Info("Successfully added device to the fleet")
	return &fleetGRPC.AddDeviceResponse{}, status.Error(codes.OK, "")
}
func (s *FleetServer) RemoveDevice(ctx context.Context, in *fleetGRPC.RemoveDeviceRequest) (*fleetGRPC.RemoveDeviceResponse, error) {
	fleetDevicesCollection := FleetDevicesCollectionByNamespace(s.systemStub, in.Namespace)

	fleetUUID, err := primitive.ObjectIDFromHex(in.FleetUUID)
	if err != nil {
		return &fleetGRPC.RemoveDeviceResponse{}, status.Error(codes.OK, "Fleet doesnt exist. Bad fleet UUID format.")
	}
	deviceUUID, err := primitive.ObjectIDFromHex(in.DeviceUUID)
	if err != nil {
		return &fleetGRPC.RemoveDeviceResponse{}, status.Error(codes.OK, "Device doesnt exist. Bad device UUID format.")
	}

	deleteResult, err := fleetDevicesCollection.DeleteOne(ctx, bson.M{"fleetUUID": fleetUUID, "deviceUUID": deviceUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &fleetGRPC.RemoveDeviceResponse{}, status.Error(codes.OK, "Namespace doesnt exist.")
			}
		}

		err = errors.New("error while removing device from the fleet: " + err.Error())
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.DeviceUUID,
			"iot_core_fleet_uuid":  in.FleetUUID,
			"namespace":            in.Namespace,
		}).Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if deleteResult.DeletedCount > 0 {
		s.logger.WithFields(logrus.Fields{
			"iot_core_device_uuid": in.DeviceUUID,
			"iot_core_fleet_uuid":  in.FleetUUID,
			"namespace":            in.Namespace,
		}).Info("Successfully removed device from the fleet")
	}
	return &fleetGRPC.RemoveDeviceResponse{}, status.Error(codes.OK, "")
}
