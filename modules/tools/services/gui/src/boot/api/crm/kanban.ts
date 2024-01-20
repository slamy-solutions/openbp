import { APIModuleBase } from '../model'

export interface TicketStage {
    namespace: string;
    uuid: string;
    name: string;
    departmentUUID: string;
    arrangementIndex: number;
}

export interface CreateTicketStageRequest {
    namespace: string
    departmentUUID: string
    name: string
}
export interface CreateTicketStageResponse {
    ticketStage: TicketStage
}

export interface GetTicketStagesRequest {
    namespace: string
    departmentUUID: string
}
export interface GetTicketStagesResponse {
    stages: TicketStage[]
}

export interface UpdateTicketStageRequest {
    namespace: string
    uuid: string
    name: string
}
export interface UpdateTicketStageResponse {
    ticketStage: TicketStage
}

export interface DeleteTicketStageRequest {
    namespace: string
    uuid: string
}
export interface DeleteTicketStageResponse {}

export interface SwapTicketStagesPrioritiesRequest {
    namespace: string
    uuid1: string
    uuid2: string
}
export interface SwapTicketStagesPrioritiesResponse {}

export class KanbanAPI extends APIModuleBase {
    async createTicketStage(params: CreateTicketStageRequest): Promise<CreateTicketStageResponse> {
        const response = await KanbanAPI._axios.post<CreateTicketStageResponse>('/crm/kanban/stage', params)
        return response.data
    }

    async getTicketStages(params: GetTicketStagesRequest): Promise<GetTicketStagesResponse> {
        const response = await KanbanAPI._axios.get<GetTicketStagesResponse>('/crm/kanban/stages', { params })
        response.data.stages.sort((a, b) => a.arrangementIndex - b.arrangementIndex)
        return response.data
    }

    async updateTicketStage(params: UpdateTicketStageRequest): Promise<UpdateTicketStageResponse> {
        const response = await KanbanAPI._axios.patch<UpdateTicketStageResponse>('/crm/kanban/stage', params)
        return response.data
    }

    async deleteTicketStage(params: DeleteTicketStageRequest): Promise<DeleteTicketStageResponse> {
        const response = await KanbanAPI._axios.delete<DeleteTicketStageResponse>('/crm/kanban/stage', { params })
        return response.data
    }

    async swapTicketStagesPriorities(params: SwapTicketStagesPrioritiesRequest): Promise<SwapTicketStagesPrioritiesResponse> {
        const response = await KanbanAPI._axios.patch<SwapTicketStagesPrioritiesResponse>('/crm/kanban/stage/swapPriority', params)
        return response.data
    }
}