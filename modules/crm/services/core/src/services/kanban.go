package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/kanban"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KanbanService struct {
	kanban.UnimplementedKanbanServiceServer

	backend backend.BackendFactory
	logger  *slog.Logger
}

func NewKanbanServer(backend backend.BackendFactory, logger *slog.Logger) *KanbanService {
	return &KanbanService{
		backend: backend,
		logger:  logger,
	}
}

func (s *KanbanService) CreateStage(ctx context.Context, in *kanban.CreateStageRequest) (*kanban.CreateStageResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	stage, err := bkd.KanbanRepository().CreateStage(ctx, in.Name, in.DepartmentUUID)
	if err != nil {
		if errors.Is(err, models.ErrTicketStageArragementConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "stage with same arragement index already exists")
		}

		err := errors.Join(errors.New("failed to create stage"), err)
		s.logger.With(slog.String("route", "CreateStage")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create stage: %s", err.Error())
	}

	return &kanban.CreateStageResponse{
		Stage: stage.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) GetStage(ctx context.Context, in *kanban.GetStageRequest) (*kanban.GetStageResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	stage, err := bkd.KanbanRepository().GetStage(ctx, in.UUID, in.UseCache)
	if err != nil {
		if errors.Is(err, models.ErrTicketStageNotFound) {
			return nil, status.Errorf(codes.NotFound, "stage not found")
		}

		err := errors.Join(errors.New("failed to get stage"), err)
		s.logger.With(slog.String("route", "GetStage")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get stage: %s", err.Error())
	}

	return &kanban.GetStageResponse{
		Stage: stage.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) GetStages(ctx context.Context, in *kanban.GetStagesRequest) (*kanban.GetStagesResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	stages, err := bkd.KanbanRepository().GetStages(ctx, in.DepartmentUUID, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get stages"), err)
		s.logger.With(slog.String("route", "GetStages")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get stages: %s", err.Error())
	}

	var responseStages []*kanban.TicketStage = make([]*kanban.TicketStage, len(stages))
	for i, stage := range stages {
		responseStages[i] = stage.ToGRPC()
	}

	return &kanban.GetStagesResponse{
		Stages: responseStages,
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) UpdateStage(ctx context.Context, in *kanban.UpdateStageRequest) (*kanban.UpdateStageResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "UpdateStage")))
	if err != nil {
		return nil, err
	}

	stage, err := bkd.KanbanRepository().UpdateStage(ctx, in.UUID, in.Name)
	if err != nil {
		if errors.Is(err, models.ErrTicketStageArragementConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "stage with same arragement index already exists")
		}

		err := errors.Join(errors.New("failed to update stage"), err)
		s.logger.With(slog.String("route", "UpdateStage")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update stage: %s", err.Error())
	}

	return &kanban.UpdateStageResponse{
		Stage: stage.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) DeleteStage(ctx context.Context, in *kanban.DeleteStageRequest) (*kanban.DeleteStageResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "DeleteStage")))
	if err != nil {
		return nil, err
	}

	_, err = bkd.KanbanRepository().DeleteStage(ctx, in.UUID)
	if err != nil {
		if errors.Is(err, models.ErrTicketStageNotFound) {
			return nil, status.Errorf(codes.NotFound, "stage not found")
		}

		err := errors.Join(errors.New("failed to delete stage"), err)
		s.logger.With(slog.String("route", "DeleteStage")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete stage: %s", err.Error())
	}

	return &kanban.DeleteStageResponse{}, status.Error(codes.OK, "")
}
func (s *KanbanService) SwapStagesOrder(ctx context.Context, in *kanban.SwapStagesOrderRequest) (*kanban.SwapStagesOrderResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "SwapStagesOrder")))
	if err != nil {
		return nil, err
	}

	err = bkd.KanbanRepository().SwapStagesOrder(ctx, in.StageUUID1, in.StageUUID2)
	if err != nil {
		if errors.Is(err, models.ErrTicketStageNotFound) {
			return nil, status.Errorf(codes.NotFound, "stage not found")
		}

		err := errors.Join(errors.New("failed to swap stages order"), err)
		s.logger.With(slog.String("route", "SwapStagesOrder")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to swap stages order: %s", err.Error())
	}

	return &kanban.SwapStagesOrderResponse{}, status.Error(codes.OK, "")
}
func (s *KanbanService) CreateTicket(ctx context.Context, in *kanban.CreateTicketRequest) (*kanban.CreateTicketResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "CreateTicket")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().CreateTicket(ctx, &models.TicketCreationInfo{
		Name:        in.Name,
		Description: in.Description,
		Files:       in.Files,
		Priority:    in.Priority,

		DepartmentUUID:    in.DepartmentUUID,
		PerformerUUID:     in.PerformerUUID,
		ClientUUID:        in.ClientUUID,
		ProjectUUID:       in.ProjectUUID,
		ContactPersonUUID: in.ContactPersonUUID,

		TrackingStoryPointsPlan: in.TrackingStoryPointsPlan,
	})
	if err != nil {
		if errors.Is(err, models.ErrTicketCreationInfoInvalid) {
			return nil, status.Errorf(codes.InvalidArgument, "ticket creation info invalid: %s", err.Error())
		}

		err := errors.Join(errors.New("failed to create ticket"), err)
		s.logger.With(slog.String("route", "CreateTicket")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create ticket: %s", err.Error())
	}

	return &kanban.CreateTicketResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) GetTicket(ctx context.Context, in *kanban.GetTicketRequest) (*kanban.GetTicketResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetTicket")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().GetTicket(ctx, in.UUID, in.UseCache)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}

		err := errors.Join(errors.New("failed to get ticket"), err)
		s.logger.With(slog.String("route", "GetTicket")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get ticket: %s", err.Error())
	}

	return &kanban.GetTicketResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) GetTickets(ctx context.Context, in *kanban.GetTicketsRequest) (*kanban.GetTicketsResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetTickets")))
	if err != nil {
		return nil, err
	}

	filter := models.TicketsFilter{
		DepartmentUUID: in.DepartmentUUID,
		PerformerUUID:  in.PerformerUUID,
	}

	tickets, err := bkd.KanbanRepository().GetTickets(ctx, in.UseCache, filter)
	if err != nil {
		err := errors.Join(errors.New("failed to get tickets"), err)
		s.logger.With(slog.String("route", "GetTickets")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get tickets: %s", err.Error())
	}

	var responseTickets []*kanban.Ticket = make([]*kanban.Ticket, len(tickets))
	for i, ticket := range tickets {
		responseTickets[i] = ticket.ToGRPC()
	}

	return &kanban.GetTicketsResponse{
		Tickets: responseTickets,
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) UpdateTicketBasicInfo(ctx context.Context, in *kanban.UpdateTicketBasicInfoRequest) (*kanban.UpdateTicketBasicInfoResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetTickets")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().UpdateTicketBasicInfo(ctx, in.UUID, in.Name, in.Description, in.Files)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}

		err := errors.Join(errors.New("failed to update ticket basic info"), err)
		s.logger.With(slog.String("route", "UpdateTicketBasicInfo")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update ticket basic info: %s", err.Error())
	}

	return &kanban.UpdateTicketBasicInfoResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) DeleteTicket(ctx context.Context, in *kanban.DeleteTicketRequest) (*kanban.DeleteTicketResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "DeleteTicket")))
	if err != nil {
		return nil, err
	}

	_, err = bkd.KanbanRepository().DeleteTicket(ctx, in.UUID)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}

		err := errors.Join(errors.New("failed to delete ticket"), err)
		s.logger.With(slog.String("route", "DeleteTicket")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete ticket: %s", err.Error())
	}

	return &kanban.DeleteTicketResponse{}, status.Error(codes.OK, "")
}
func (s *KanbanService) UpdateTicketStage(ctx context.Context, in *kanban.UpdateTicketStageRequest) (*kanban.UpdateTicketStageResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "UpdateTicketStage")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().UpdateTicketStage(ctx, in.UUID, in.TicketStageUUID)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}
		if errors.Is(err, models.ErrTicketStageNotFound) {
			return nil, status.Errorf(codes.NotFound, "stage not found")
		}

		err := errors.Join(errors.New("failed to update ticket stage"), err)
		s.logger.With(slog.String("route", "UpdateTicketStage")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update ticket stage: %s", err.Error())
	}

	return &kanban.UpdateTicketStageResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) UpdateTicketPriority(ctx context.Context, in *kanban.UpdateTicketPriorityRequest) (*kanban.UpdateTicketPriorityResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "UpdateTicketPriority")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().UpdateTicketPriority(ctx, in.UUID, uint32(in.Priority))
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}

		err := errors.Join(errors.New("failed to update ticket priority"), err)
		s.logger.With(slog.String("route", "UpdateTicketPriority")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update ticket priority: %s", err.Error())
	}

	return &kanban.UpdateTicketPriorityResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *KanbanService) CloseTicket(ctx context.Context, in *kanban.CloseTicketRequest) (*kanban.CloseTicketResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "CloseTicket")))
	if err != nil {
		return nil, err
	}

	ticket, err := bkd.KanbanRepository().CloseTicket(ctx, in.UUID)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticket not found")
		}

		err := errors.Join(errors.New("failed to close ticket"), err)
		s.logger.With(slog.String("route", "CloseTicket")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to close ticket: %s", err.Error())
	}

	return &kanban.CloseTicketResponse{
		Ticket: ticket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
