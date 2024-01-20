package kanban

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/client"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/department"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/performer"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/project"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stageInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	ArrangementIndex int64              `bson:"arrangementIndex"`
	DepartmentUUID   primitive.ObjectID `bson:"departmentUUID"`
}

func (s *stageInMongo) ToBackendModel(namespace string) *models.TicketStage {
	return &models.TicketStage{
		Namespace:        namespace,
		UUID:             s.UUID.Hex(),
		Name:             s.Name,
		ArrangementIndex: s.ArrangementIndex,
		DepartmentUUID:   s.DepartmentUUID.Hex(),
	}
}

type ticketPlanningInfoInMongo struct {
	ExpectedStartDate *time.Time `bson:"expectedStartDate,omitempty"`
}

func (p *ticketPlanningInfoInMongo) ToBackendModel() models.TicketPlanningInfo {
	return models.TicketPlanningInfo{
		ExpectedStartDate: p.ExpectedStartDate,
	}
}

type ticketTrackingInfoInMongo struct {
	StoryPointsPlan uint32 `bson:"storyPointsPlan"`
	StoryPointsFact uint32 `bson:"storyPointsFact"`

	TrackedTime time.Duration `bson:"trackedTime"`
}

func (t *ticketTrackingInfoInMongo) ToBackendModel() models.TicketTrackingInfo {
	return models.TicketTrackingInfo{
		StoryPointsPlan: t.StoryPointsPlan,
		StoryPointsFact: t.StoryPointsFact,

		TrackedTime: t.TrackedTime,
	}
}

type ticketFeedEntryInMongo struct {
	Type      models.TicketFeedEntryType `bson:"type"`
	Files     []string                   `bson:"files"`
	Timestamp time.Time                  `bson:"timestamp"`
}

func (t *ticketFeedEntryInMongo) ToBackendModel() models.TicketFeedEntry {
	return models.TicketFeedEntry{
		Type:      t.Type,
		Files:     t.Files,
		Timestamp: t.Timestamp,
	}
}

type ticketInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`

	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	Files       []string `bson:"files"`

	Priority int32 `bson:"priority"`

	ClientUUID        primitive.ObjectID `bson:"clientUUID"`
	ContactPersonUUID primitive.ObjectID `bson:"contactPersonUUID"`
	DepartmentUUID    primitive.ObjectID `bson:"departmentUUID"`
	PerformerUUID     primitive.ObjectID `bson:"performerUUID"`
	ProjectUUID       primitive.ObjectID `bson:"projectUUID"`

	/* Dynamically loaded fields */
	Client        *client.ClientInMongo         `bson:"client,omitempty"`
	ContactPerson *client.ContactPersonInMongo  `bson:"contactPerson,omitempty"`
	Department    *department.DepartmentInMongo `bson:"department,omitempty"`
	Performer     *performer.PerformerInMongo   `bson:"performer,omitempty"`
	Project       *project.ProjectInMongo       `bson:"project,omitempty"`
	Stage         *stageInMongo                 `bson:"stage,omitempty"`
	/* --- */

	Feed []ticketFeedEntryInMongo `bson:"feed"`

	StageUUID primitive.ObjectID `bson:"ticketStageUUID,omitempty"`

	Planning ticketPlanningInfoInMongo `bson:"planning"`
	Tracking ticketTrackingInfoInMongo `bson:"tracking"`

	CloseDate *time.Time `bson:"closeDate,omitempty"`
	Created   time.Time  `bson:"created"`
	Updated   time.Time  `bson:"updated"`
	Version   int        `bson:"version"`
}

func (t *ticketInMongo) ToBackendModel(namespace string, clientContactPersons []models.ContactPerson, performerName string, performerAvatar string) models.Ticket {
	feed := make([]models.TicketFeedEntry, len(t.Feed))
	for i, entry := range t.Feed {
		feed[i] = entry.ToBackendModel()
	}

	var clientModel *models.Client
	if t.Client != nil {
		clientModel = t.Client.ToBackendModel(namespace, clientContactPersons)
	}

	var contactPersonModel *models.ContactPerson
	if t.ContactPerson != nil {
		contactPersonModel = t.ContactPerson.ToBackendModel(namespace)
	}

	var departmentModel *models.Department
	if t.Department != nil {
		departmentModel = t.Department.ToBackendModel(namespace)
	}

	var performerModel *models.Performer
	if t.Performer != nil {
		performerModel = t.Performer.ToBackendModel(namespace, performerName, performerAvatar)
	}

	var projectModel *models.Project
	if t.Project != nil {
		projectModel = t.Project.ToBackendModel(namespace)
	}

	var stageModel *models.TicketStage
	if t.Stage != nil {
		stageModel = t.Stage.ToBackendModel(namespace)
	}

	return models.Ticket{
		Namespace: namespace,
		UUID:      t.UUID.Hex(),

		Name:        t.Name,
		Description: t.Description,
		Files:       t.Files,

		ClientUUID:        t.ClientUUID.Hex(),
		Client:            clientModel,
		ContactPersonUUID: t.ContactPersonUUID.Hex(),
		ContactPerson:     contactPersonModel,
		DepartmentUUID:    t.DepartmentUUID.Hex(),
		Department:        departmentModel,
		PerformerUUID:     t.PerformerUUID.Hex(),
		Performer:         performerModel,
		ProjectUUID:       t.ProjectUUID.Hex(),
		Project:           projectModel,
		StageUUID:         t.StageUUID.Hex(),
		Stage:             stageModel,

		Planning: t.Planning.ToBackendModel(),
		Tracking: t.Tracking.ToBackendModel(),
		Priority: t.Priority,

		Feed: feed,

		CloseDate: t.CloseDate,
		Created:   t.Created,
		Updated:   t.Updated,
		Version:   t.Version,
	}
}

func (t *ticketInMongo) ToBackendModelWithFetch(ctx context.Context, namespace string, useCache bool, clientRepository *client.ClientRepository, departmentRepository *department.DepartmentRepository, performerRepository *performer.PerformerRepository, projectRepository *project.ProjectRepository, kanbanRepository *KanbanRepository) (*models.Ticket, error) {
	tasksWaitGroup := sync.WaitGroup{}

	var clientModel *models.Client
	var clientModelError error
	tasksWaitGroup.Add(1)
	go func() {
		defer tasksWaitGroup.Done()
		clientModel, clientModelError = clientRepository.Get(ctx, t.ClientUUID.Hex(), useCache)
	}()

	var departmentModel *models.Department
	var departmentModelError error
	tasksWaitGroup.Add(1)
	go func() {
		defer tasksWaitGroup.Done()
		departmentModel, departmentModelError = departmentRepository.Get(ctx, t.DepartmentUUID.Hex(), useCache)
	}()

	var performerModel *models.Performer
	var performerModelError error
	tasksWaitGroup.Add(1)
	go func() {
		defer tasksWaitGroup.Done()
		performerModel, performerModelError = performerRepository.Get(ctx, t.PerformerUUID.Hex(), useCache)
	}()

	var projectModel *models.Project
	var projectModelError error
	tasksWaitGroup.Add(1)
	go func() {
		defer tasksWaitGroup.Done()
		projectModel, projectModelError = projectRepository.Get(ctx, t.ProjectUUID.Hex(), useCache)
	}()

	var stageModel *models.TicketStage
	var stageModelError error
	tasksWaitGroup.Add(1)
	go func() {
		defer tasksWaitGroup.Done()
		stageModel, stageModelError = kanbanRepository.GetStage(ctx, t.StageUUID.Hex(), useCache)
	}()

	feed := make([]models.TicketFeedEntry, len(t.Feed))
	for i, entry := range t.Feed {
		feed[i] = entry.ToBackendModel()
	}

	tasksWaitGroup.Wait()

	if clientModelError != nil && !errors.Is(clientModelError, models.ErrClientNotFound) {
		return nil, errors.Join(clientModelError, errors.New("failed to fetch client"))
	}
	if departmentModelError != nil && !errors.Is(departmentModelError, models.ErrDepartmentNotFound) {
		return nil, errors.Join(departmentModelError, errors.New("failed to fetch department"))
	}
	if performerModelError != nil && !errors.Is(performerModelError, models.ErrPerformerNotFound) {
		return nil, errors.Join(performerModelError, errors.New("failed to fetch performer"))
	}
	if projectModelError != nil && !errors.Is(projectModelError, models.ErrProjectNotFound) {
		return nil, errors.Join(projectModelError, errors.New("failed to fetch project"))
	}
	if stageModelError != nil && !errors.Is(stageModelError, models.ErrTicketStageNotFound) {
		return nil, errors.Join(stageModelError, errors.New("failed to fetch stage"))
	}

	var contactPersonModel *models.ContactPerson
	if clientModelError == nil && t.ContactPersonUUID != primitive.NilObjectID {
		for _, contactPerson := range clientModel.ContactPersons {
			if contactPerson.UUID == t.ContactPersonUUID.Hex() {
				contactPersonModel = &contactPerson
				break
			}
		}
	}

	return &models.Ticket{
		Namespace: namespace,
		UUID:      t.UUID.Hex(),

		Name:        t.Name,
		Description: t.Description,
		Files:       t.Files,

		ClientUUID:        t.ClientUUID.Hex(),
		Client:            clientModel,
		ContactPersonUUID: t.ContactPersonUUID.Hex(),
		ContactPerson:     contactPersonModel,
		DepartmentUUID:    t.DepartmentUUID.Hex(),
		Department:        departmentModel,
		PerformerUUID:     t.PerformerUUID.Hex(),
		Performer:         performerModel,
		ProjectUUID:       t.ProjectUUID.Hex(),
		Project:           projectModel,
		StageUUID:         t.StageUUID.Hex(),
		Stage:             stageModel,

		Planning: t.Planning.ToBackendModel(),
		Tracking: t.Tracking.ToBackendModel(),
		Priority: t.Priority,

		Feed: feed,

		CloseDate: t.CloseDate,
		Created:   t.Created,
		Updated:   t.Updated,
		Version:   t.Version,
	}, nil
}
