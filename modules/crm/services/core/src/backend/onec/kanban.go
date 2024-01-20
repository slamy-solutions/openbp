package onec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
)

type ticketStage struct {
	ID             string `json:"id,omitempty"`
	IDX            string `json:"idx"`
	Name           string `json:"name"`
	DepartmentUUID string `json:"departmentId"`
}

type ticketData struct {
	Id          string `json:"id,omitempty"`
	ColumnId    int    `json:"columnId"`
	ColumnName  string `json:"columnName,omitempty"`
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Client      struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"client"`
	Contact struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
		Tel1 string `json:"tel1"`
		Tel2 string `json:"tel2"`
	} `json:"contact"`
	Department struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"department"`
	Performer struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"performer"`
	Project struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"project"`
	Files []struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"files"`
	CreatedDate    string `json:"createdDate"`
	StartDate      string `json:"startDate"`
	IsStartTiming  bool   `json:"isStartTiming,omitempty"`
	Storypoints    int    `json:"storypoints"`
	Fact           int    `json:"fact"`
	LeadTime       int    `json:"leadTime"`
	IsDayTask      bool   `json:"isDayTask"`
	TaskTemplateId string `json:"taskTemplateId"`
	NumberDayTasks bool   `json:"numberDayTasks"`
	IsTracking     bool   `json:"isTracking"`
}

func ticketDataFromTicket(ticket *models.Ticket) (*ticketData, error) {
	files := make([]struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	}, len(ticket.Files))
	for i, file := range ticket.Files {
		files[i].Id = file
	}

	columnId, err := strconv.Atoi(ticket.StageUUID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse column id"), err)
	}

	return &ticketData{
		Id:          ticket.UUID,
		ColumnId:    columnId,
		Priority:    int(ticket.Priority),
		Name:        ticket.Name,
		Description: ticket.Description,
		Client: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.ClientUUID},
		Contact: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
			Tel1 string `json:"tel1"`
			Tel2 string `json:"tel2"`
		}{Id: ticket.ContactPersonUUID},
		Department: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.DepartmentUUID},
		Performer: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.PerformerUUID},
		Project: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.ProjectUUID},
		Files:          files,
		CreatedDate:    ticket.Created.Format(time.RFC3339),
		StartDate:      ticket.Created.Format(time.RFC3339),
		IsStartTiming:  false,
		Storypoints:    int(ticket.Tracking.StoryPointsPlan),
		Fact:           int(ticket.Tracking.StoryPointsFact),
		LeadTime:       0,
		IsDayTask:      false,
		ColumnName:     "",
		TaskTemplateId: "",
		NumberDayTasks: false,
		IsTracking:     false,
	}, nil
}

func (t *ticketData) toTicket(ctx context.Context, namespace string, useCache bool, kanbanRep *kanbanRepository) (*models.Ticket, error) {
	files := make([]string, len(t.Files))
	for i, file := range t.Files {
		files[i] = file.Id
	}

	creationDate, err := time.Parse(time.RFC3339, t.CreatedDate)
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse creation date"), models.ErrBackendMissBehaviour)
	}

	dataCollectorsWaitGroup := sync.WaitGroup{}

	var client *models.Client
	var contactPerson *models.ContactPerson
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, clientRequestError := kanbanRep.clientRep.Get(ctx, t.Client.Id, useCache)
		if clientRequestError != nil {
			client = response
			for _, contact := range response.ContactPersons {
				if contact.UUID == t.Contact.Id {
					contactPerson = &contact
					return
				}
			}
		}
	}()

	var department *models.Department
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, departmentRequestError := kanbanRep.departmentRep.Get(ctx, t.Department.Id, useCache)
		if departmentRequestError != nil {
			department = response
		}
	}()

	var performer *models.Performer
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, performerRequestError := kanbanRep.performerRep.Get(ctx, t.Performer.Id, useCache)
		if performerRequestError != nil {
			performer = response
		}
	}()

	var project *models.Project
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, projectRequestError := kanbanRep.projectRep.Get(ctx, t.Project.Id, useCache)
		if projectRequestError != nil {
			project = response
		}
	}()

	var stage *models.TicketStage
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, stageRequestError := kanbanRep.GetStage(ctx, string(t.ColumnId), useCache)
		if stageRequestError != nil {
			stage = response
		}
	}()

	var feed *[]models.TicketFeedEntry
	dataCollectorsWaitGroup.Add(1)
	go func() {
		defer dataCollectorsWaitGroup.Done()
		response, feedRequestError := kanbanRep.getTicketFeed(ctx, t.Id, t.Client.Id, useCache)
		if feedRequestError != nil {
			feed = &response
		}
	}()

	dataCollectorsWaitGroup.Wait()
	return &models.Ticket{
		Namespace:   namespace,
		UUID:        t.Id,
		Name:        t.Name,
		Description: t.Description,

		Files:    files,
		Priority: int32(t.Priority),

		ClientUUID:        t.Client.Id,
		Client:            client,
		ContactPersonUUID: t.Contact.Id,
		ContactPerson:     contactPerson,
		DepartmentUUID:    t.Department.Id,
		Department:        department,
		PerformerUUID:     t.Performer.Id,
		Performer:         performer,
		ProjectUUID:       t.Project.Id,
		Project:           project,
		StageUUID:         string(t.ColumnId),
		Stage:             stage,
		Feed:              *feed,

		Planning: models.TicketPlanningInfo{
			ExpectedStartDate: nil,
		},
		Tracking: models.TicketTrackingInfo{
			StoryPointsPlan: uint32(t.Storypoints),
			StoryPointsFact: uint32(t.Fact),
			TrackedTime:     0,
		},

		Updated: creationDate,
		Created: creationDate,
		Version: 0,
	}, nil
}

type ticketComment struct {
	Id         string `json:"id"`
	Type       string `json:"type"` //'commentPerformer' | 'commentClient' | 'callIn' | 'callOut'
	TaskId     string `json:"taskId"`
	TaskName   string `json:"taskName"`
	ClientId   string `json:"clientId"`
	AuthorId   string `json:"authorId"`
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	Files      []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"files"`
	Date       string `json:"date"`
	Seen       bool   `json:"seen"`
	SeenClient bool   `json:"seen_client"`
}

type kanbanRepository struct {
	logger    *slog.Logger
	connector *connector.OneCConnector
	namespace string

	clientRep     models.ClientRepository
	departmentRep models.DepartmentRepository
	performerRep  models.PerformerRepository
	projectRep    models.ProjectRepository
}

func (r *kanbanRepository) CreateStage(ctx context.Context, name string, departmentUUID string) (*models.TicketStage, error) {

	arrangementIndex := time.Now().UnixMilli()

	columnCreateRequest := ticketStage{
		IDX:            fmt.Sprintf("%d", arrangementIndex),
		Name:           name,
		DepartmentUUID: departmentUUID,
	}

	response, statusCode, err := connector.MakeRequest[ticketStage](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		columnCreateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to create stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Created stage", slog.Group("stage", "name", name, "uuid", response.ID))
	return &models.TicketStage{
		Namespace:        r.namespace,
		ArrangementIndex: arrangementIndex,
		UUID:             response.ID,
		Name:             name,
		DepartmentUUID:   departmentUUID,
	}, nil
}
func (r *kanbanRepository) GetStage(ctx context.Context, uuid string, useCache bool) (*models.TicketStage, error) {
	response, statusCode, err := connector.MakeRequest[[]ticketStage](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to create stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	for _, stage := range *response {
		arrangementIndex, err := strconv.ParseInt(stage.IDX, 10, 64)
		if err != nil {
			return nil, errors.Join(errors.New("failed to parse arrangement index"), err)
		}

		if stage.ID == uuid {
			return &models.TicketStage{
				Namespace:        r.namespace,
				ArrangementIndex: int64(arrangementIndex),
				UUID:             stage.ID,
				Name:             stage.Name,
				DepartmentUUID:   stage.DepartmentUUID,
			}, nil
		}
	}

	return nil, models.ErrTicketStageNotFound
}
func (r *kanbanRepository) GetStages(ctx context.Context, departmentUUID string, useCache bool) ([]models.TicketStage, error) {
	response, statusCode, err := connector.MakeRequest[[]ticketStage](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to create stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	ticketStage := make([]models.TicketStage, len(*response))
	for i, stage := range *response {
		arrangementIndex, err := strconv.ParseInt(stage.IDX, 10, 64)
		if err != nil {
			return nil, errors.Join(errors.New("failed to parse arrangement index"), err)
		}

		ticketStage[i] = models.TicketStage{
			Namespace:        r.namespace,
			ArrangementIndex: arrangementIndex,
			UUID:             stage.ID,
			Name:             stage.Name,
			DepartmentUUID:   stage.DepartmentUUID,
		}
	}

	return ticketStage, nil
}
func (r *kanbanRepository) UpdateStage(ctx context.Context, uuid string, name string) (*models.TicketStage, error) {
	currentTicketStage, err := r.GetStage(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	currentTicketStage.Name = name

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		currentTicketStage,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Updated stage", slog.Group("stage", "name", name, "uuid", uuid))
	return currentTicketStage, nil
}
func (r *kanbanRepository) DeleteStage(ctx context.Context, uuid string) (*models.TicketStage, error) {
	currentTicketStage, err := r.GetStage(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"DELETE",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		currentTicketStage,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to delete stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Deleted stage", slog.Group("stage", "name", currentTicketStage.Name, "uuid", uuid))
	return currentTicketStage, nil
}
func (r *kanbanRepository) SwapStagesOrder(ctx context.Context, uuid1 string, uuid2 string) error {
	stage1, err := r.GetStage(ctx, uuid1, false)
	if err != nil {
		return err
	}

	stage2, err := r.GetStage(ctx, uuid2, false)
	if err != nil {
		return err
	}

	stage1.ArrangementIndex, stage2.ArrangementIndex = stage2.ArrangementIndex, stage1.ArrangementIndex

	//TODO: maybe something atomic here?

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		stage1,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update stage. Invalid status code from the backend: %d", statusCode), err)
		return err
	}

	_, statusCode, err = connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/column/%s", r.connector.ServerURL, r.connector.ServerToken),
		stage2,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update stage. Invalid status code from the backend: %d", statusCode), err)
		return err
	}

	r.logger.Info("Swapped stages", slog.Group("stage1", "name", stage1.Name, "uuid", uuid1), slog.Group("stage2", "name", stage2.Name, "uuid", uuid2))
	return nil
}

func (r *kanbanRepository) getTicketFeed(ctx context.Context, ticketUUID string, clientUUID string, useCache bool) ([]models.TicketFeedEntry, error) {
	comments, statusCode, err := connector.MakeRequest[[]ticketComment](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/comments/%s?taskId=%s&clientId=%s", r.connector.ServerURL, r.connector.ServerToken, ticketUUID, clientUUID),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get ticket comments"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to get ticket comments. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	ticketFeed := make([]models.TicketFeedEntry, len(*comments))
	for i, comment := range *comments {
		timestamp, err := time.Parse(time.RFC3339, comment.Date)
		if err != nil {
			return nil, errors.Join(errors.New("failed to parse comment timestamp"), models.ErrBackendMissBehaviour)
		}

		ticketFeed[i] = models.TicketFeedEntry{
			Type:      models.TicketFeedEntryType(comment.Type),
			Files:     make([]string, len(comment.Files)),
			Timestamp: timestamp,
		}
		for j, file := range comment.Files {
			ticketFeed[i].Files[j] = file.Id
		}
	}

	return ticketFeed, nil
}

func (r *kanbanRepository) CreateTicket(ctx context.Context, ticket *models.TicketCreationInfo) (*models.Ticket, error) {
	files := make([]struct {
		Id   string `json:"id"`
		Name string `json:"name,omitempty"`
	}, len(ticket.Files))
	for i, file := range ticket.Files {
		files[i].Id = file
	}

	ticketRequest := ticketData{
		Name:        ticket.Name,
		Description: ticket.Description,
		Files:       files,
		Priority:    int(ticket.Priority),

		Client: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.ClientUUID},

		Contact: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
			Tel1 string `json:"tel1"`
			Tel2 string `json:"tel2"`
		}{Id: ticket.ContactPersonUUID},

		Department: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.DepartmentUUID},

		Performer: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.PerformerUUID},

		Project: struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{Id: ticket.ProjectUUID},

		Storypoints: int(ticket.TrackingStoryPointsPlan),
	}

	response, statusCode, err := connector.MakeRequest[ticketData](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		ticketRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create ticket"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to create ticket. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Created ticket", slog.Group("ticket", "name", ticket.Name, "uuid", response.Id))
	return response.toTicket(ctx, r.namespace, false, r)
}
func (r *kanbanRepository) GetTicket(ctx context.Context, uuid string, useCache bool) (*models.Ticket, error) {
	//TODO: Better ticket get

	tickets, err := r.GetTickets(ctx, useCache, models.TicketsFilter{})
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.UUID == uuid {
			return &ticket, nil
		}
	}

	return nil, models.ErrTicketNotFound
}
func (r *kanbanRepository) GetTickets(ctx context.Context, useCache bool, filter models.TicketsFilter) ([]models.Ticket, error) {
	target, err := url.Parse(fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken))
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse url"), err)
	}
	q := target.Query()
	if filter.DepartmentUUID != nil {
		q.Set("departmentId", *filter.DepartmentUUID)
	}
	target.RawQuery = q.Encode()

	response, statusCode, err := connector.MakeRequest[[]ticketData](
		ctx,
		r.connector,
		"POST",
		target.String(),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create ticket"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to create ticket. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	tickets := make([]models.Ticket, 0, len(*response))
	for _, t := range *response {
		if filter.PerformerUUID != nil && t.Performer.Id != *filter.PerformerUUID {
			continue
		}

		//TODO: this part should be optimize / parallelized

		ticket, err := t.toTicket(ctx, r.namespace, false, r)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket)
	}

	return tickets, nil
}
func (r *kanbanRepository) UpdateTicketBasicInfo(ctx context.Context, uuid string, name string, description string, files []string) (*models.Ticket, error) {
	currentTicket, err := r.GetTicket(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	currentTicket.Name = name
	currentTicket.Description = description
	currentTicket.Files = files

	ticketRequest, err := ticketDataFromTicket(currentTicket)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		ticketRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update ticket basic info"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update ticket basic info. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Updated ticket basic info", slog.Group("ticket", "name", name, "uuid", uuid))
	return currentTicket, nil
}

func (r *kanbanRepository) DeleteTicket(ctx context.Context, uuid string) (*models.Ticket, error) {
	currentTicket, err := r.GetTicket(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"DELETE",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		uuid,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete ticket"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to delete ticket. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Deleted ticket", slog.Group("ticket", "name", currentTicket.Name, "uuid", uuid))
	return currentTicket, nil
}

func (r *kanbanRepository) UpdateTicketStage(ctx context.Context, ticketUUID string, ticketStageUUID string) (*models.Ticket, error) {
	currentTicket, err := r.GetTicket(ctx, ticketUUID, false)
	if err != nil {
		return nil, err
	}

	currentTicket.StageUUID = ticketStageUUID

	ticketRequest, err := ticketDataFromTicket(currentTicket)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		ticketRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update ticket stage"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update ticket stage. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Updated ticket stage", slog.Group("ticket", "name", currentTicket.Name, "uuid", ticketUUID))
	return currentTicket, nil
}
func (r *kanbanRepository) UpdateTicketPriority(ctx context.Context, ticketUUID string, priority uint32) (*models.Ticket, error) {
	currentTicket, err := r.GetTicket(ctx, ticketUUID, false)
	if err != nil {
		return nil, err
	}

	currentTicket.Priority = int32(priority)

	ticketRequest, err := ticketDataFromTicket(currentTicket)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		ticketRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update ticket priority"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to update ticket priority. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Updated ticket priority", slog.Group("ticket", "name", currentTicket.Name, "uuid", ticketUUID))
	return currentTicket, nil
}
func (r *kanbanRepository) CloseTicket(ctx context.Context, ticketUUID string) (*models.Ticket, error) {
	currentTicket, err := r.GetTicket(ctx, ticketUUID, false)
	if err != nil {
		return nil, err
	}

	ticketRequest, err := ticketDataFromTicket(currentTicket)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"PATCH",
		fmt.Sprintf("%s/task/%s", r.connector.ServerURL, r.connector.ServerToken),
		ticketRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to close ticket"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, fmt.Errorf("failed to close ticket. Invalid status code from the backend: %d", statusCode), err)
		return nil, err
	}

	r.logger.Info("Closed ticket", slog.Group("ticket", "name", currentTicket.Name, "uuid", ticketUUID))
	return currentTicket, nil
}
