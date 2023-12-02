package cacheable

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const (
	kanbanStageDataCacheKeyPrefix = "crm_kanban_stage_data_"
	kanbanStageDataCacheTTL       = time.Minute * 5

	kanbanStageListCacheKeyPrefix = "crm_kandan_stage_list"
	kanbanStageListCacheTTL       = time.Minute * 5

	kanbanTicketDataCacheKeyPrefix = "crm_kanban_ticket_data_"
	kanbanTicketDataCacheTTL       = time.Second * 30

	kanbanTicketListCacheKeyPrefix = "crm_kanban_ticket_list_"
	kanbanTicketListCacheTTL       = time.Second * 15
)

func MakeKanbanStageDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", kanbanStageDataCacheKeyPrefix, namespace, uuid)
}

func MakeKanbanStageListCacheKey(namespace string, departmentUUID string) string {
	return fmt.Sprintf("%s_%s_%s", kanbanStageListCacheKeyPrefix, namespace, departmentUUID)
}

func MakeKanbanTicketDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", kanbanTicketDataCacheKeyPrefix, namespace, uuid)
}

func MakeKanbanTicketListCacheKey(namespace string, filter models.TicketsFilter) string {
	departmentUUID := ""
	if filter.DepartmentUUID != nil {
		departmentUUID = *filter.DepartmentUUID
	}

	performerUUID := ""
	if filter.PerformerUUID != nil {
		performerUUID = *filter.PerformerUUID
	}

	return fmt.Sprintf("%s_%s_%s_%s", kanbanTicketListCacheKeyPrefix, namespace, departmentUUID, performerUUID)
}

func MakePossibleCacheKeysForKanbanTicket(tiket *models.Ticket) []string {
	return []string{
		MakeKanbanTicketDataCacheKey(tiket.Namespace, tiket.UUID),
		MakeKanbanTicketListCacheKey(tiket.Namespace, models.TicketsFilter{
			DepartmentUUID: &tiket.DepartmentUUID,
			PerformerUUID:  &tiket.PerformerUUID,
		}),
		MakeKanbanTicketListCacheKey(tiket.Namespace, models.TicketsFilter{
			DepartmentUUID: nil,
			PerformerUUID:  &tiket.PerformerUUID,
		}),
		MakeKanbanTicketListCacheKey(tiket.Namespace, models.TicketsFilter{
			DepartmentUUID: &tiket.DepartmentUUID,
			PerformerUUID:  nil,
		}),
		MakeKanbanTicketListCacheKey(tiket.Namespace, models.TicketsFilter{
			DepartmentUUID: nil,
			PerformerUUID:  nil,
		}),
	}
}

type kanbanRepository struct {
	wrapedRespository models.KanbanRepository
	logger            *slog.Logger
	namespace         string
	systemStub        *system.SystemStub
}

func (r *kanbanRepository) CreateStage(ctx context.Context, name string, departmentUUID string, arrangementIndex uint32) (*models.TicketStage, error) {
	stage, err := r.wrapedRespository.CreateStage(ctx, name, departmentUUID, arrangementIndex)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeKanbanStageListCacheKey(r.namespace, departmentUUID))
	}

	return stage, err
}
func (r *kanbanRepository) GetStage(ctx context.Context, uuid string, useCache bool) (*models.TicketStage, error) {
	var cacheKey string
	if useCache {
		cacheKey = MakeKanbanStageDataCacheKey(r.namespace, uuid)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var stage models.TicketStage
			err := json.Unmarshal(cacheBytes, &stage)
			if err != nil {
				r.logger.Warn("failed to unmarshal stage from cache", "error", err.Error())
			} else {
				return &stage, nil
			}
		}
	}

	stage, err := r.wrapedRespository.GetStage(ctx, uuid, useCache)
	if err == nil {
		cacheBytes, _ := json.Marshal(stage)
		cacheSetErr := r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, kanbanStageDataCacheTTL)
		if cacheSetErr != nil {
			r.logger.Warn("failed to set stage to cache", "error", cacheSetErr.Error())
		}
	}

	return stage, err
}

func (r *kanbanRepository) GetStages(ctx context.Context, departmentUUID string, useCache bool) ([]models.TicketStage, error) {
	var cacheKey string
	if useCache {
		cacheKey = MakeKanbanStageListCacheKey(r.namespace, departmentUUID)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var stages []models.TicketStage
			err := json.Unmarshal(cacheBytes, &stages)
			if err != nil {
				r.logger.Warn("failed to unmarshal stages from cache", "error", err.Error())
			} else {
				return stages, nil
			}
		}
	}

	stages, err := r.wrapedRespository.GetStages(ctx, departmentUUID, useCache)
	if err == nil {
		cacheBytes, _ := json.Marshal(stages)
		cacheSetErr := r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, kanbanStageListCacheTTL)
		if cacheSetErr != nil {
			r.logger.Warn("failed to set stages to cache", "error", cacheSetErr.Error())
		}
	}

	return stages, err
}
func (r *kanbanRepository) UpdateStage(ctx context.Context, uuid string, name string, arrangementIndex uint32) (*models.TicketStage, error) {
	stage, err := r.wrapedRespository.UpdateStage(ctx, uuid, name, arrangementIndex)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakeKanbanStageListCacheKey(r.namespace, stage.DepartmentUUID), MakeKanbanTicketListCacheKey(r.namespace, models.TicketsFilter{DepartmentUUID: &stage.DepartmentUUID}), MakeKanbanStageDataCacheKey(r.namespace, uuid))
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove stages from cache", "error", cacheRemoveErr.Error())
		}
	}
	return stage, err
}
func (r *kanbanRepository) DeleteStage(ctx context.Context, uuid string) (*models.TicketStage, error) {
	stage, err := r.wrapedRespository.DeleteStage(ctx, uuid)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakeKanbanStageListCacheKey(r.namespace, stage.DepartmentUUID), MakeKanbanStageDataCacheKey(r.namespace, uuid), MakeKanbanTicketListCacheKey(r.namespace, models.TicketsFilter{DepartmentUUID: &stage.DepartmentUUID}))
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove stages from cache", "error", cacheRemoveErr.Error())
		}
	}
	return stage, err
}

func (r *kanbanRepository) CreateTicket(ctx context.Context, ticket *models.TicketCreationInfo) (*models.Ticket, error) {
	tiket, err := r.wrapedRespository.CreateTicket(ctx, ticket)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakePossibleCacheKeysForKanbanTicket(tiket)...)
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove ticket from cache", "error", cacheRemoveErr.Error())
		}
	}
	return tiket, err

}
func (r *kanbanRepository) GetTicket(ctx context.Context, uuid string, useCache bool) (*models.Ticket, error) {
	var cacheKey string
	if useCache {
		cacheKey = MakeKanbanTicketDataCacheKey(r.namespace, uuid)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var ticket models.Ticket
			err := json.Unmarshal(cacheBytes, &ticket)
			if err != nil {
				r.logger.Warn("failed to unmarshal ticket from cache", "error", err.Error())
			} else {
				return &ticket, nil
			}
		}
	}

	tiket, err := r.wrapedRespository.GetTicket(ctx, uuid, useCache)
	if err == nil {
		cacheBytes, _ := json.Marshal(tiket)
		cacheSetErr := r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, kanbanTicketDataCacheTTL)
		if cacheSetErr != nil {
			r.logger.Warn("failed to set ticket to cache", "error", cacheSetErr.Error())
		}
	}

	return tiket, err
}
func (r *kanbanRepository) GetTickets(ctx context.Context, useCache bool, filter models.TicketsFilter) ([]models.Ticket, error) {
	var cacheKey string
	if useCache {
		cacheKey = MakeKanbanTicketListCacheKey(r.namespace, filter)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var tickets []models.Ticket
			err := json.Unmarshal(cacheBytes, &tickets)
			if err != nil {
				r.logger.Warn("failed to unmarshal tickets from cache", "error", err.Error())
			} else {
				return tickets, nil
			}
		}
	}

	tickets, err := r.wrapedRespository.GetTickets(ctx, useCache, filter)
	if err == nil {
		cacheBytes, _ := json.Marshal(tickets)
		cacheSetErr := r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, kanbanTicketListCacheTTL)
		if cacheSetErr != nil {
			r.logger.Warn("failed to set tickets to cache", "error", cacheSetErr.Error())
		}
	}

	return tickets, err
}
func (r *kanbanRepository) DeleteTicket(ctx context.Context, uuid string) (*models.Ticket, error) {
	tiket, err := r.wrapedRespository.DeleteTicket(ctx, uuid)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakePossibleCacheKeysForKanbanTicket(tiket)...)
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove ticket from cache", "error", cacheRemoveErr.Error())
		}
	}
	return tiket, err
}

func (r *kanbanRepository) UpdateTicketStage(ctx context.Context, ticketUUID string, ticketStageUUID string) (*models.Ticket, error) {
	tiket, err := r.wrapedRespository.UpdateTicketStage(ctx, ticketUUID, ticketStageUUID)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakePossibleCacheKeysForKanbanTicket(tiket)...)
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove ticket from cache", "error", cacheRemoveErr.Error())
		}
	}
	return tiket, err
}
func (r *kanbanRepository) UpdateTicketPriority(ctx context.Context, ticketUUID string, priority uint32) (*models.Ticket, error) {
	tiket, err := r.wrapedRespository.UpdateTicketPriority(ctx, ticketUUID, priority)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakePossibleCacheKeysForKanbanTicket(tiket)...)
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove ticket from cache", "error", cacheRemoveErr.Error())
		}
	}
	return tiket, err
}
func (r *kanbanRepository) CloseTicket(ctx context.Context, ticketUUID string) (*models.Ticket, error) {
	tiket, err := r.wrapedRespository.CloseTicket(ctx, ticketUUID)
	if err == nil {
		cacheRemoveErr := r.systemStub.Cache.Remove(ctx, MakePossibleCacheKeysForKanbanTicket(tiket)...)
		if cacheRemoveErr != nil {
			r.logger.Warn("failed to remove ticket from cache", "error", cacheRemoveErr.Error())
		}
	}
	return tiket, err
}
