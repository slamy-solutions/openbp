package models

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTicketStageUUIDInvalid = errors.New("stage uuid is invalid")
	ErrTicketStageNotFound    = errors.New("stage not found")

	ErrTicketCreationInfoInvalid = errors.New("ticket creation info is invalid")
	ErrTicketNotFound            = errors.New("ticket not found")
	ErrTicketUUIDInvalid         = errors.New("ticket uuid is invalid")
)

type TicketStage struct {
	Namespace        string `json:"namespace"`
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	ArrangementIndex uint32 `json:"arrangementIndex"`
	DepartmentUUID   string `json:"departmentUUID"`
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

	CloseDate   *time.Time `json:"closeDate"`
	CreatedDate time.Time  `json:"createdDate"`
	UpdatedDate time.Time  `json:"updatedDate"`
	Version     int        `json:"version"`
}

type TicketsFilter struct {
	// If seted return tickts only for specified department
	DepartmentUUID *string `json:"departmentUUID"`
	PerformerUUID  *string `json:"performerUUID"`
}

type KanbanRepository interface {
	CreateStage(ctx context.Context, name string, departmentUUID string, arrangementIndex uint32) (*TicketStage, error)
	GetStage(ctx context.Context, uuid string, useCache bool) (*TicketStage, error)
	GetStages(ctx context.Context, departmentUUID string, useCache bool) ([]TicketStage, error)
	UpdateStage(ctx context.Context, uuid string, name string, arrangementIndex uint32) (*TicketStage, error)
	DeleteStage(ctx context.Context, uuid string) (*TicketStage, error)

	CreateTicket(ctx context.Context, ticket *TicketCreationInfo) (*Ticket, error)
	GetTicket(ctx context.Context, uuid string, useCache bool) (*Ticket, error)
	GetTickets(ctx context.Context, useCache bool, filter TicketsFilter) ([]Ticket, error)
	DeleteTicket(ctx context.Context, uuid string) (*Ticket, error)

	UpdateTicketStage(ctx context.Context, ticketUUID string, ticketStageUUID string) (*Ticket, error)
	UpdateTicketPriority(ctx context.Context, ticketUUID string, priority uint32) (*Ticket, error)
	CloseTicket(ctx context.Context, ticketUUID string) (*Ticket, error)
}
