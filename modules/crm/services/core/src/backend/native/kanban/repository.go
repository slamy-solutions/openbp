package kanban

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/client"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/department"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/performer"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/project"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KanbanRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
	nativeStub *native.NativeStub

	clientRepository     *client.ClientRepository
	performerRepository  *performer.PerformerRepository
	departmentRepository *department.DepartmentRepository
	projectRepository    *project.ProjectRepository
}

func NewKanbanRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub, clientRepository *client.ClientRepository, performerRepository *performer.PerformerRepository, departmentRepository *department.DepartmentRepository, projectRepository *project.ProjectRepository) KanbanRepository {
	return KanbanRepository{
		logger:               logger,
		namespace:            namespace,
		systemStub:           systemStub,
		nativeStub:           nativeStub,
		clientRepository:     clientRepository,
		performerRepository:  performerRepository,
		departmentRepository: departmentRepository,
		projectRepository:    projectRepository,
	}
}

func (r *KanbanRepository) CreateStage(ctx context.Context, name string, departmentUUID string, arrangementIndex uint32) (*models.TicketStage, error) {
	defaprtmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrDepartmentUUIDInvalid
	}

	collection := GetStageCollection(r.systemStub, r.namespace)

	stage := stageInMongo{
		Name:             name,
		DepartmentUUID:   defaprtmentUUIDObject,
		ArrangementIndex: arrangementIndex,
	}
	insertResponse, err := collection.InsertOne(ctx, stage)
	if err != nil {
		err = errors.Join(errors.New("failed to create stage in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	stage.UUID = insertResponse.InsertedID.(primitive.ObjectID)
	return stage.ToBackendModel(r.namespace), nil
}

func (r *KanbanRepository) GetStage(ctx context.Context, uuid string, useCache bool) (*models.TicketStage, error) {
	stageUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrTicketStageUUIDInvalid
	}

	collection := GetStageCollection(r.systemStub, r.namespace)

	var stage stageInMongo
	err = collection.FindOne(ctx, bson.M{
		"_id": stageUUID,
	}).Decode(&stage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketStageNotFound
		}

		err = errors.Join(errors.New("failed to get stage from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return stage.ToBackendModel(r.namespace), nil
}

func (r *KanbanRepository) GetStages(ctx context.Context, departmentUUID string, useCache bool) ([]models.TicketStage, error) {
	defaprtmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrDepartmentUUIDInvalid
	}

	collection := GetStageCollection(r.systemStub, r.namespace)
	cursor, err := collection.Find(ctx, bson.M{
		"departmentUUID": defaprtmentUUIDObject,
	})
	if err != nil {
		err = errors.Join(errors.New("failed to get stages from the database. failed to open cursor"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var stages []stageInMongo
	err = cursor.All(ctx, &stages)
	if err != nil {
		err = errors.Join(errors.New("failed to get stages from the database. failed to decode stage"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	stagesResponse := make([]models.TicketStage, len(stages))
	for i, stage := range stages {
		stagesResponse[i] = *stage.ToBackendModel(r.namespace)
	}

	return stagesResponse, nil
}
func (r *KanbanRepository) UpdateStage(ctx context.Context, uuid string, name string, arrangementIndex uint32) (*models.TicketStage, error) {
	stageUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrTicketStageUUIDInvalid
	}

	collection := GetStageCollection(r.systemStub, r.namespace)

	var stage stageInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{
		"_id": stageUUID,
	}, bson.M{
		"$set": bson.M{
			"name":             name,
			"arrangementIndex": arrangementIndex,
		},
	}).Decode(&stage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketStageNotFound
		}

		err = errors.Join(errors.New("failed to update stage in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return stage.ToBackendModel(r.namespace), nil
}
func (r *KanbanRepository) DeleteStage(ctx context.Context, uuid string) (*models.TicketStage, error) {
	stageUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrTicketStageUUIDInvalid
	}

	collection := GetStageCollection(r.systemStub, r.namespace)

	var stage stageInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{
		"_id": stageUUID,
	}).Decode(&stage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketStageNotFound
		}

		err = errors.Join(errors.New("failed to delete stage in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return stage.ToBackendModel(r.namespace), nil
}

func (r *KanbanRepository) CreateTicket(ctx context.Context, ticket *models.TicketCreationInfo) (*models.Ticket, error) {
	clientUUID, err := primitive.ObjectIDFromHex(ticket.ClientUUID)
	if err != nil {
		return nil, errors.Join(models.ErrTicketCreationInfoInvalid, models.ErrClientUUIDInvalid)
	}

	contactPersonUUID, err := primitive.ObjectIDFromHex(ticket.ContactPersonUUID)
	if err != nil {
		return nil, errors.Join(models.ErrTicketCreationInfoInvalid, models.ErrClientContactPersonUUIDInvalid)
	}

	departmentUUID, err := primitive.ObjectIDFromHex(ticket.DepartmentUUID)
	if err != nil {
		return nil, errors.Join(models.ErrTicketCreationInfoInvalid, models.ErrDepartmentUUIDInvalid)
	}

	performerUUID, err := primitive.ObjectIDFromHex(ticket.PerformerUUID)
	if err != nil {
		return nil, errors.Join(models.ErrTicketCreationInfoInvalid, models.ErrPerformerUUIDInvalid)
	}

	projectUUID, err := primitive.ObjectIDFromHex(ticket.ProjectUUID)
	if err != nil {
		return nil, errors.Join(models.ErrTicketCreationInfoInvalid, models.ErrProjectUUIDInvalid)
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	creationTime := time.Now().UTC()

	ticketInMongo := ticketInMongo{
		Name:              ticket.Name,
		Description:       ticket.Description,
		Files:             ticket.Files,
		ClientUUID:        clientUUID,
		ContactPersonUUID: contactPersonUUID,
		DepartmentUUID:    departmentUUID,
		PerformerUUID:     performerUUID,
		ProjectUUID:       projectUUID,
		Planning: ticketPlanningInfoInMongo{
			ExpectedStartDate: nil,
		},
		Tracking: ticketTrackingInfoInMongo{
			StoryPointsPlan: ticket.TrackingStoryPointsPlan,
			StoryPointsFact: 0,
			TrackedTime:     0,
		},
		Priority: ticket.Priority,

		Feed: []ticketFeedEntryInMongo{},

		StageUUID: primitive.NilObjectID,

		CloseDate: nil,

		Created: creationTime,
		Updated: creationTime,
		Version: 1,
	}
	insertResponse, err := collection.InsertOne(ctx, ticketInMongo)
	if err != nil {
		err = errors.Join(errors.New("failed to create ticket in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketInMongo.UUID = insertResponse.InsertedID.(primitive.ObjectID)
	ticketModel, err := ticketInMongo.ToBackendModelWithFetch(ctx, r.namespace, false, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to create ticket in the database: failed to get ticket information after it creation"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil
}
func (r *KanbanRepository) GetTicket(ctx context.Context, uuid string, useCache bool) (*models.Ticket, error) {
	ticketUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrTicketNotFound
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	var ticket ticketInMongo
	err = collection.FindOne(ctx, bson.M{
		"_id": ticketUUID,
	}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketNotFound
		}

		err = errors.Join(errors.New("failed to get ticket from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketModel, err := ticket.ToBackendModelWithFetch(ctx, r.namespace, useCache, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to get ticket from the database: failed to fetch ticket information"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil

}
func (r *KanbanRepository) GetTickets(ctx context.Context, useCache bool, filter models.TicketsFilter) ([]models.Ticket, error) {

	aggregation := []bson.M{}

	if filter.DepartmentUUID != nil || filter.PerformerUUID != nil {
		f := bson.M{}
		if filter.DepartmentUUID != nil {
			f["departmentUUID"] = filter.DepartmentUUID
		}
		if filter.PerformerUUID != nil {
			f["performerUUID"] = filter.PerformerUUID
		}
		aggregation = append(aggregation, bson.M{
			"$match": f,
		})
	}

	// Add client information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         client.GetClientCollection(r.systemStub, r.namespace).Name(),
			"localField":   "clientUUID",
			"foreignField": "_id",
			"as":           "client",
		},
		"$unwind": bson.M{
			"path":                       "$client",
			"preserveNullAndEmptyArrays": true,
		},
	})

	// Add contact person information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         client.GetClientContactPersonCollection(r.systemStub, r.namespace).Name(),
			"localField":   "contactPersonUUID",
			"foreignField": "_id",
			"as":           "contactPerson",
		},
		"$unwind": bson.M{
			"path":                       "$contactPerson",
			"preserveNullAndEmptyArrays": true,
		},
	})

	// Add department information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         department.GetDepartmentCollection(r.systemStub, r.namespace).Name(),
			"localField":   "departmentUUID",
			"foreignField": "_id",
			"as":           "department",
		},
		"$unwind": bson.M{
			"path":                       "$department",
			"preserveNullAndEmptyArrays": true,
		},
	})

	// Add performer information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         performer.GetPerformerCollection(r.systemStub, r.namespace).Name(),
			"localField":   "performerUUID",
			"foreignField": "_id",
			"as":           "performer",
		},
		"$unwind": bson.M{
			"path":                       "$performer",
			"preserveNullAndEmptyArrays": true,
		},
	})

	// Add project information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         project.GetProjectCollection(r.systemStub, r.namespace).Name(),
			"localField":   "projectUUID",
			"foreignField": "_id",
			"as":           "project",
		},
		"$unwind": bson.M{
			"path":                       "$project",
			"preserveNullAndEmptyArrays": true,
		},
	})

	// Add stage information
	aggregation = append(aggregation, bson.M{
		"$lookup": bson.M{
			"from":         GetStageCollection(r.systemStub, r.namespace).Name(),
			"localField":   "stageUUID",
			"foreignField": "_id",
			"as":           "stage",
		},
		"$unwind": bson.M{
			"path":                       "$stage",
			"preserveNullAndEmptyArrays": true,
		},
	})

	collection := GetTicketCollection(r.systemStub, r.namespace)
	cursor, err := collection.Aggregate(ctx, aggregation)
	if err != nil {
		err = errors.Join(errors.New("failed to get tickets from the database. failed to open cursor"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []ticketInMongo
	err = cursor.All(ctx, &tickets)
	if err != nil {
		err = errors.Join(errors.New("failed to get tickets from the database. failed to decode ticket"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	//TODO: everything to one big pipeline?

	// Find and map all the client contact persons
	clients := make([]string, 0, 10)
	for _, ticket := range tickets {
		clients = append(clients, ticket.ClientUUID.Hex())
	}
	clients = slices.Compact(clients)
	clientContactPersonsMap, err := r.clientRepository.GetContactPersonsForClients(ctx, clients)
	if err != nil {
		err = errors.Join(errors.New("failed to get tickets from the database. failed to get client contact persons"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	// Find and map all the performers names and avatars
	performerUsers := make([]string, 0, 10)
	for _, ticket := range tickets {
		if ticket.Performer != nil {
			performerUsers = append(performerUsers, ticket.Performer.UserUUID.Hex())
		}
	}
	performerUsers = slices.Compact(performerUsers)

	users := make([]*user.User, 0, len(performerUsers))
	//TODO: move to goroutine
	for _, userUUID := range performerUsers {
		userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
			Namespace: r.namespace,
			Uuid:      userUUID,
			UseCache:  useCache,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				continue
			}

			err = errors.Join(errors.New("failed to get user"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		users = append(users, userResponse.User)
	}
	usersMap := make(map[string]*user.User, len(users))
	for _, user := range users {
		usersMap[user.Uuid] = user
	}

	// Map all the tickets
	ticketsResponse := make([]models.Ticket, len(tickets))
	for i, ticket := range tickets {
		var performerName string = ""
		var performerAvatar string = ""
		if ticket.Performer != nil {
			if user, ok := usersMap[ticket.Performer.UserUUID.Hex()]; ok {
				performerName = user.FullName
				performerAvatar = user.Avatar
			}
		}

		ticketsResponse[i] = ticket.ToBackendModel(r.namespace, clientContactPersonsMap[ticket.ClientUUID.Hex()], performerName, performerAvatar)
	}

	return ticketsResponse, nil
}
func (r *KanbanRepository) DeleteTicket(ctx context.Context, uuid string) (*models.Ticket, error) {
	ticketUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrTicketNotFound
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	var ticket ticketInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{
		"_id": ticketUUID,
	}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketNotFound
		}

		err = errors.Join(errors.New("failed to delete ticket in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketModel, err := ticket.ToBackendModelWithFetch(ctx, r.namespace, false, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to delete ticket in the database: failed to fetch ticket information"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil
}

func (r *KanbanRepository) UpdateTicketStage(ctx context.Context, ticketUUID string, ticketStageUUID string) (*models.Ticket, error) {
	ticketUUIDObject, err := primitive.ObjectIDFromHex(ticketUUID)
	if err != nil {
		return nil, models.ErrTicketNotFound
	}

	ticketStageUUIDObject, err := primitive.ObjectIDFromHex(ticketStageUUID)
	if err != nil {
		return nil, models.ErrTicketStageUUIDInvalid
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	var ticket ticketInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{
		"_id": ticketUUIDObject,
	}, bson.M{
		"$set": bson.M{
			"stageUUID": ticketStageUUIDObject,
		},
		"$inc":         bson.M{"version": 1},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		//TODO: push infomration to feed
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketNotFound
		}

		err = errors.Join(errors.New("failed to update ticket in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketModel, err := ticket.ToBackendModelWithFetch(ctx, r.namespace, false, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to update ticket in the database: failed to fetch ticket information"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil

}
func (r *KanbanRepository) UpdateTicketPriority(ctx context.Context, ticketUUID string, priority uint32) (*models.Ticket, error) {
	ticketUUIDObject, err := primitive.ObjectIDFromHex(ticketUUID)
	if err != nil {
		return nil, models.ErrTicketNotFound
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	var ticket ticketInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{
		"_id": ticketUUIDObject,
	}, bson.M{
		"$set": bson.M{
			"priority": priority,
		},
		"$inc":         bson.M{"version": 1},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketNotFound
		}

		err = errors.Join(errors.New("failed to update ticket in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketModel, err := ticket.ToBackendModelWithFetch(ctx, r.namespace, false, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to update ticket in the database: failed to fetch ticket information"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil
}
func (r *KanbanRepository) CloseTicket(ctx context.Context, ticketUUID string) (*models.Ticket, error) {
	ticketUUIDObject, err := primitive.ObjectIDFromHex(ticketUUID)
	if err != nil {
		return nil, models.ErrTicketNotFound
	}

	collection := GetTicketCollection(r.systemStub, r.namespace)

	var ticket ticketInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{
		"_id": ticketUUIDObject,
	}, bson.M{
		"$inc": bson.M{"version": 1},
		"$currentDate": bson.M{
			"updated":   bson.M{"$type": "timestamp"},
			"closeDate": bson.M{"$type": "timestamp"},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&ticket)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrTicketNotFound
		}

		err = errors.Join(errors.New("failed to close ticket in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	ticketModel, err := ticket.ToBackendModelWithFetch(ctx, r.namespace, false, r.clientRepository, r.departmentRepository, r.performerRepository, r.projectRepository, r)
	if err != nil {
		err = errors.Join(errors.New("failed to close ticket in the database: failed to fetch ticket information"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return ticketModel, nil
}
