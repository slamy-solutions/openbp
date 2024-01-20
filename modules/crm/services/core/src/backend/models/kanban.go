package models

import (
	"context"
	"errors"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/client"
	department "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/department"
	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/kanban"
	performer "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/performer"
	project "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/project"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrTicketStageUUIDInvalid        = errors.New("stage uuid is invalid")
	ErrTicketStageNotFound           = errors.New("stage not found")
	ErrTicketStageArragementConflict = errors.New("stage arrangement index conflict")

	ErrTicketCreationInfoInvalid = errors.New("ticket creation info is invalid")
	ErrTicketNotFound            = errors.New("ticket not found")
	ErrTicketUUIDInvalid         = errors.New("ticket uuid is invalid")
)

type TicketStage struct {
	Namespace        string `json:"namespace"`
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	ArrangementIndex int64  `json:"arrangementIndex"`
	DepartmentUUID   string `json:"departmentUUID"`
}

func (s *TicketStage) ToGRPC() *kanban.TicketStage {
	return &kanban.TicketStage{
		Namespace:        s.Namespace,
		Uuid:             s.UUID,
		Name:             s.Name,
		ArrangementIndex: s.ArrangementIndex,
		DepartmentUUID:   s.DepartmentUUID,
	}
}

type TicketFeedEntryType string

const (
	TicketFeedEntryTypeCommentPerformer TicketFeedEntryType = "commentPerformer"
	TicketFeedEntryTypeCommentClient    TicketFeedEntryType = "commentClient"
	TicketFeedEntryTypeCallIn           TicketFeedEntryType = "callIn"
	TicketFeedEntryTypeCallOut          TicketFeedEntryType = "callOut"
)

type TicketFeedEntry struct {
	Type      TicketFeedEntryType `json:"type"`
	Files     []string            `json:"files"`
	Timestamp time.Time           `json:"timestamp"`
}

type TicketPlanningInfo struct {
	ExpectedStartDate *time.Time `json:"expectedStartDate"`
}

type TicketTrackingInfo struct {
	StoryPointsPlan uint32
	StoryPointsFact uint32

	TrackedTime time.Duration
}

type TicketCreationInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Priority    int32    `json:"priority"`

	ClientUUID        string `json:"clientUUID"`
	ContactPersonUUID string `json:"contactPersonUUID"`
	DepartmentUUID    string `json:"departmentUUID"`
	PerformerUUID     string `json:"performerUUID"`
	ProjectUUID       string `json:"projectUUID"`

	TrackingStoryPointsPlan uint32 `json:"trackingStoryPointsPlan"`
}

type Ticket struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Priority    int32    `json:"priority"`

	ClientUUID        string         `json:"clientUUID"`
	Client            *Client        `json:"client"`
	ContactPersonUUID string         `json:"contactPersonUUID"`
	ContactPerson     *ContactPerson `json:"contactPerson"`
	DepartmentUUID    string         `json:"departmentUUID"`
	Department        *Department    `json:"department"`
	PerformerUUID     string         `json:"performerUUID"`
	Performer         *Performer     `json:"performer"`
	ProjectUUID       string         `json:"projectUUID"`
	Project           *Project       `json:"project"`
	StageUUID         string         `json:"ticketStageUUID"`
	Stage             *TicketStage   `json:"ticketStage"`

	Planning TicketPlanningInfo `json:"planning"`
	Tracking TicketTrackingInfo `json:"tracking"`

	Feed []TicketFeedEntry `json:"feed"`

	CloseDate *time.Time `json:"closeDate"`
	Created   time.Time  `json:"createdDate"`
	Updated   time.Time  `json:"updatedDate"`
	Version   int        `json:"version"`
}

func (t *Ticket) ToGRPC() *kanban.Ticket {

	var clientData *client.Client
	if t.Client != nil {
		clientData = t.Client.ToGRPC()
	}

	var contactPersonData *client.ContactPerson
	if t.ContactPerson != nil {
		contactPersonData = t.ContactPerson.ToGRPC()
	}

	var departmentData *department.Department
	if t.Department != nil {
		departmentData = t.Department.ToGRPC()
	}

	var performerData *performer.Performer
	if t.Performer != nil {
		performerData = t.Performer.ToGRPC()
	}

	var projectData *project.Project
	if t.Project != nil {
		projectData = t.Project.ToGRPC()
	}

	var stageData *kanban.TicketStage
	if t.Stage != nil {
		stageData = t.Stage.ToGRPC()
	}

	var planningExpectedStartDate *timestamppb.Timestamp
	if t.Planning.ExpectedStartDate != nil {
		planningExpectedStartDate = timestamppb.New(*t.Planning.ExpectedStartDate)
	}

	var feed []*kanban.TicketFeedEntry
	for _, entry := range t.Feed {
		feed = append(feed, &kanban.TicketFeedEntry{
			Type:      kanban.TicketFeedEntryType(kanban.TicketFeedEntryType_value[string(entry.Type)]),
			Files:     entry.Files,
			Timestamp: timestamppb.New(entry.Timestamp),
		})
	}

	var closeDate *timestamppb.Timestamp
	if t.CloseDate != nil {
		closeDate = timestamppb.New(*t.CloseDate)
	}

	return &kanban.Ticket{
		Namespace: t.Namespace,
		UUID:      t.UUID,

		Name:        t.Name,
		Description: t.Description,
		Files:       t.Files,
		Priority:    t.Priority,

		ClientUUID:        t.ClientUUID,
		Client:            clientData,
		ContactPersonUUID: t.ContactPersonUUID,
		ContactPerson:     contactPersonData,
		DepartmentUUID:    t.DepartmentUUID,
		Department:        departmentData,
		PerformerUUID:     t.PerformerUUID,
		Performer:         performerData,
		ProjectUUID:       t.ProjectUUID,
		Project:           projectData,
		StageUUID:         t.StageUUID,
		Stage:             stageData,

		Planning: &kanban.Ticket_Planning{
			ExpectedStartDate: planningExpectedStartDate,
		},
		Tracking: &kanban.Ticket_Tracking{
			StoryPointsPlan: t.Tracking.StoryPointsPlan,
			StoryPointsFact: t.Tracking.StoryPointsFact,
			TrackedTime:     uint64(t.Tracking.TrackedTime),
		},

		Feed: feed,

		CloseDate: closeDate,
		Created:   timestamppb.New(t.Created),
		Updated:   timestamppb.New(t.Updated),
		Version:   int32(t.Version),
	}
}

type TicketsFilter struct {
	// If seted return tickts only for specified department
	DepartmentUUID *string `json:"departmentUUID"`
	PerformerUUID  *string `json:"performerUUID"`
}

type KanbanRepository interface {
	CreateStage(ctx context.Context, name string, departmentUUID string) (*TicketStage, error)
	GetStage(ctx context.Context, uuid string, useCache bool) (*TicketStage, error)
	GetStages(ctx context.Context, departmentUUID string, useCache bool) ([]TicketStage, error)
	UpdateStage(ctx context.Context, uuid string, name string) (*TicketStage, error)
	DeleteStage(ctx context.Context, uuid string) (*TicketStage, error)
	SwapStagesOrder(ctx context.Context, uuid1 string, uuid2 string) error

	CreateTicket(ctx context.Context, ticket *TicketCreationInfo) (*Ticket, error)
	GetTicket(ctx context.Context, uuid string, useCache bool) (*Ticket, error)
	GetTickets(ctx context.Context, useCache bool, filter TicketsFilter) ([]Ticket, error)
	UpdateTicketBasicInfo(ctx context.Context, uuid string, name string, description string, files []string) (*Ticket, error)
	DeleteTicket(ctx context.Context, uuid string) (*Ticket, error)

	UpdateTicketStage(ctx context.Context, ticketUUID string, ticketStageUUID string) (*Ticket, error)
	UpdateTicketPriority(ctx context.Context, ticketUUID string, priority uint32) (*Ticket, error)
	CloseTicket(ctx context.Context, ticketUUID string) (*Ticket, error)
}
