import { APIModuleBase } from '../../../model';
import { SyncLogEntry } from './models'

export interface SyncNowRequest {
    serverUUID: string
}
export interface SyncNowResponse {}

export interface ListLogRequest {
    serverUUID: string
    skip: number
    limit: number
}
export interface ListLogResponse {
    logs: Array<SyncLogEntry>
    totalCount: number
}

export class SyncAPI extends APIModuleBase {
    async syncNow(params: SyncNowRequest): Promise<SyncNowResponse> {
        const response = await SyncAPI._axios.post<SyncNowResponse>('/iot/integration/balena/sync/now', params)
        return response.data
    }

    async listLog(params: ListLogRequest): Promise<ListLogResponse> {
        const response = await SyncAPI._axios.get<ListLogResponse>('/iot/integration/balena/sync/log', { params })
        const logs = response.data.logs
        logs.forEach(log => {
            log.timestamp = new Date(log.timestamp)
        })
        return response.data
    }
}