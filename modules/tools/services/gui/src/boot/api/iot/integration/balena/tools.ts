import { APIModuleBase } from '../../../model';
import { SyncLogEntry } from './models'

export type ConnectionStatus = 'OK' | 'BAD_URL' | 'SERVER_UNAVAILABLE' | 'SERVER_BAD_RESPONSE'

export interface VerifyConnectiopnDataRequest {
    url: string
    apiToken: string
}
export interface VerifyConnectiopnDataResponse {
    status: ConnectionStatus
    message: string
}

export class ToolsAPI extends APIModuleBase {
    async verifyConnectionData(params: VerifyConnectiopnDataRequest): Promise<VerifyConnectiopnDataResponse> {
        const response = await ToolsAPI._axios.post<VerifyConnectiopnDataResponse>('/iot/integration/balena/tools/verifyConnectionData', params)
        return response.data
    }
}