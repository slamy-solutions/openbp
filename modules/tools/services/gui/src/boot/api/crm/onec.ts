import { APIModuleBase } from '../model'

export interface SyncEvent {
    uuid: string
    success: boolean
    errorMessage: string
    timestamp: Date
    log: string
}

export interface SyncNowRequest {
    namespace: string
}
export interface SyncNowResponse {
    success: boolean
    errorMessage: string
}

export interface GetSyncLogRequest {
    namespace: string
    skip: number
    limit: number
}
export interface GetSyncLogResponse {
    events: Array<SyncEvent>
    totalCount: number
}

export class OneCAPI extends APIModuleBase {
    async syncNow(params: SyncNowRequest): Promise<SyncNowResponse> {
        const response = await OneCAPI._axios.post<SyncNowResponse>('/crm/onec/sync/now', params)
        return response.data
    }

    async getSyncLog(params: GetSyncLogRequest): Promise<GetSyncLogResponse> {
        const response = await OneCAPI._axios.get<GetSyncLogResponse>('/crm/onec/sync/log', { params })
        for (const event of response.data.events) {
            event.timestamp = new Date(event.timestamp)
        }
        return response.data
    }
}