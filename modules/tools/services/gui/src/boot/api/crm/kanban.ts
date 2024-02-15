import { APIModuleBase } from '../model'

import { Client, ContactPerson } from './client'
import { Department } from './department'
import { Performer } from './performer'
import { Project } from './project'

export interface Ticket {
    namespace: string
    UUID: string
    name: string
    description: string

    files: string[]
    priority: number

    clientUUID: string
    client?: Client
    contactPersonUUID: string
    contactPerson?: ContactPerson
    departmentUUID: string
    department?: Department
    performerUUID: string
    performer?: Performer
    projectUUID: string
    project?: Project
    ticketStageUUID: string
    ticketStage?: TicketStage
}

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

export interface GetTicketsRequest {
    namespace: string
    departmentUUID?: string
    performerUUID?: string
}
export interface GetTicketsResponse {
    tickets: Ticket[]
}

export interface CreateTicketRequest {
    namespace: string
    name: string
    description: string

    clientUUID: string
    contactPersonUUID: string
    departmentUUID: string
    performerUUID: string
    projectUUID: string

    trackingStoryPointsPlan: number
}
export interface CreateTicketResponse {
    ticket: Ticket
}

export interface DeleteTicketRequest {
    namespace: string
    uuid: string
}
export interface DeleteTicketResponse {}

export interface CloseTicketRequest {
    namespace: string
    uuid: string
}
export interface CloseTicketResponse {}

export interface UpdateTicketBasicInfoRequest {
    namespace: string
    uuid: string
    name: string
    description: string
    files: []
}
export interface UpdateTicketBasicInfoResponse {
    ticket: Ticket
}


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

    async getTickets(params: GetTicketsRequest): Promise<GetTicketsResponse> {
        const response = await KanbanAPI._axios.get<GetTicketsResponse>('/crm/kanban/tickets', { params })
        return response.data
    }

    async createTicket(params: CreateTicketRequest): Promise<CreateTicketResponse> {
        const response = await KanbanAPI._axios.post<CreateTicketResponse>('/crm/kanban/ticket', params)
        return response.data
    }

    async deleteTicket(params: DeleteTicketRequest): Promise<DeleteTicketResponse> {
        const response = await KanbanAPI._axios.delete<DeleteTicketResponse>('/crm/kanban/ticket', { params })
        return response.data
    }

    async closeTicket(params: CloseTicketRequest): Promise<CloseTicketResponse> {
        const response = await KanbanAPI._axios.patch<CloseTicketResponse>('/crm/kanban/ticket/close', params)
        return response.data
    }

    async updateTicketBasicInfo(params: UpdateTicketBasicInfoRequest): Promise<UpdateTicketBasicInfoResponse> {
        const response = await KanbanAPI._axios.patch<UpdateTicketBasicInfoResponse>('/crm/kanban/ticket', params)
        return response.data
    }
}