import { APIModuleBase } from '../model'

export type BackendType = 'NATIVE' | 'ONE_C'

export interface NativeBackendSettings {
    backendType: 'NATIVE'
}
export interface OneCBackendSettings {
    backendType: 'ONE_C'
    backendURL: string
    token: string
}

export type Settings = NativeBackendSettings | OneCBackendSettings

export interface GetRequest {
    namespace: string
}
export interface GetResponse {
    settings: Settings
}

export interface UpdateRequest {
    namespace: string
    settings: Settings
}

export interface CheckOneCConnectionRequest {
    backendURL: string
    token: string
}
export interface CheckOneCConnectionResponse {
    success: boolean
    statusCode: number
    message: string
}


export class SettingsAPI extends APIModuleBase {

    async get(params: GetRequest): Promise<GetResponse> {
        const response = await SettingsAPI._axios.get<GetResponse>('/crm/settings', { params })
        return response.data
    }

    async update(params: UpdateRequest): Promise<void> {
        await SettingsAPI._axios.patch('/crm/settings', params)
    }

    async checkOneCConnection(params: CheckOneCConnectionRequest): Promise<CheckOneCConnectionResponse> {
        const response = await SettingsAPI._axios.post<CheckOneCConnectionResponse>('/crm/settings/onec/connection', params)
        return response.data
    }
}