import { APIModuleBase } from '../model'

export interface Performer {
    namespace: string
    uuid: string
    departmentUUID: string
    userUUID: string

    name: string
    avatarURL: string
}

export interface CreatePerformerRequest {
    namespace: string
    departmentUUID: string
    userUUID: string
}
export interface CreatePerformerResponse {
    performer: Performer
}

export interface GetPerformersRequest {
    namespace: string
}
export interface GetPerformersResponse {
    performers: Performer[]
}

export interface DeletePerformerRequest {
    namespace: string
    uuid: string
}
export interface DeletePerformerResponse {}

export class PerformerAPI extends APIModuleBase {
    async create(params: CreatePerformerRequest): Promise<CreatePerformerResponse> {
        const response = await PerformerAPI._axios.post<CreatePerformerResponse>('/crm/performers/performer', params)
        return response.data
    }

    async getAll(params: GetPerformersRequest): Promise<GetPerformersResponse> {
        const response = await PerformerAPI._axios.get<GetPerformersResponse>('/crm/performers', { params })
        return response.data
    }

    async delete(params: DeletePerformerRequest): Promise<DeletePerformerResponse> {
        const response = await PerformerAPI._axios.delete<DeletePerformerResponse>('/crm/performers/performer', { params })
        return response.data
    }
}