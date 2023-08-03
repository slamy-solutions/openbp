import { APIModuleBase } from '../../../model';
import { Server, SyncLogEntry } from './models'

export interface ListRequest {
    namespace: string
    skip: number
    limit: number
}
export interface ListResponse {
    servers: Array<{
        server: Server
        lastSyncLog?: SyncLogEntry 
    }>
    totalCount: number
}

export interface CreateRequest {
    namespace: string
    name: string
    description: string
    baseUrl: string
    authToken: string
}
export interface CreateResponse {
    server: Server
}

export interface UpdateRequest {
    uuid: string
    newDescription: string
    newBaseUrl: string
    newAuthToken: string
}
export interface UpdateResponse {
    server: Server
}

export interface DeleteRequest {
    uuid: string
}
export interface DeleteResponse {}

export interface SetEnabledRequest {
    uuid: string
    enabled: boolean
}
export interface SetEnabledResponse {
    server: Server
}

export class ServerAPI extends APIModuleBase {

    async list(params: ListRequest): Promise<ListResponse> {
        const response = await ServerAPI._axios.get<ListResponse>('/iot/integration/balena/servers', { params })
        const servers = response.data.servers
        servers.forEach(server => {
            server.server.created = new Date(server.server.created)
            server.server.updated = new Date(server.server.updated)
            if (server.lastSyncLog) {
                server.lastSyncLog.timestamp = new Date(server.lastSyncLog.timestamp)
            }
        })
        return { servers, totalCount: response.data.totalCount }
    }

    async create(params: CreateRequest): Promise<CreateResponse> {
        const response = await ServerAPI._axios.post<CreateResponse>('/iot/integration/balena/servers/server', params)
        const server = response.data.server
        server.created = new Date(server.created)
        server.updated = new Date(server.updated)
        return { server }
    }

    async update(params: UpdateRequest): Promise<UpdateResponse> {
        const response = await ServerAPI._axios.patch<UpdateResponse>('/iot/integration/balena/servers/server', params)
        const server = response.data.server
        server.created = new Date(server.created)
        server.updated = new Date(server.updated)
        return { server }
    }

    async delete(params: DeleteRequest): Promise<DeleteResponse> {
        await ServerAPI._axios.delete<DeleteResponse>('/iot/integration/balena/servers/server', { params })
        return {}
    }

    async setEnabled(params: SetEnabledRequest): Promise<SetEnabledResponse> {
        const response = await ServerAPI._axios.patch<SetEnabledResponse>('/iot/integration/balena/servers/server/enabled', params)
        const server = response.data.server
        server.created = new Date(server.created)
        server.updated = new Date(server.updated)
        return { server }
    }
}