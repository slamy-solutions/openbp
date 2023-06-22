package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/ticket"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TicketInMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

type TicketServer struct {
	ticket.UnimplementedTicketServiceServer

	logger *logrus.Entry

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

func (t *TicketInMongo) ToGRPCTicket(namespaceName string) *ticket.Ticket {
	return &ticket.Ticket{
		Namespace:   namespaceName,
		Uuid:        t.UUID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		Created:     timestamppb.New(t.Created),
		Updated:     timestamppb.New(t.Updated),
		Version:     t.Version,
	}
}

func NewTicketServer(systemStub *system.SystemStub, nativeStub *native.NativeStub) *TicketServer {
	return &TicketServer{
		systemStub: systemStub,
		nativeStub: nativeStub,

		logger: logrus.StandardLogger().WithField("service", "ticket"),
	}
}

func TicketCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	if namespaceName == "" {
		return systemStub.DB.Database("openbp_global").Collection("crm_core_ticket")
	} else {
		db := systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespaceName))
		return db.Collection("crm_core_ticket")
	}
}

func (s *TicketServer) Create(ctx context.Context, in *ticket.CreateTicketRequest) (*ticket.CreateTicketResponse, error) {
	collection := TicketCollectionByNamespace(s.systemStub, in.Namespace)

	creationTime := time.Now().UTC()
	insertData := TicketInMongo{
		Title:       in.Title,
		Description: in.Description,
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}

	insertResult, err := collection.InsertOne(ctx, insertData)
	if err != nil {
		err = errors.New("error while inserting new ticket to the database: " + err.Error())
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	insertData.UUID = insertResult.InsertedID.(primitive.ObjectID)
	return &ticket.CreateTicketResponse{Ticket: insertData.ToGRPCTicket(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *TicketServer) Get(ctx context.Context, in *ticket.GetTicketRequest) (*ticket.GetTicketresponse, error) {
	collection := TicketCollectionByNamespace(s.systemStub, in.Namespace)

	ticketUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Ticket not found. Bad ticket UUID format.")
	}

	var foundedTicket TicketInMongo
	err = collection.FindOne(ctx, bson.M{"_id": ticketUUID}).Decode(&foundedTicket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Tiket with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Ticket not found. Probably namespace doesnt exist")
			}
		}

		err = errors.New("error while getting ticket from the database: " + err.Error())
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ticket.GetTicketresponse{Ticket: foundedTicket.ToGRPCTicket(in.Namespace)}, status.Error(codes.OK, "")
}
