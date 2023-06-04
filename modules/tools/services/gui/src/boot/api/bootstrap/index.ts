import { APIModuleBase } from "../model"

export interface GetStatusResponse {
    vaultSealed: boolean
    rootUserCreated: boolean
    rootUserCreationBlocked: boolean
}

export interface CreateRootUserRequest {
    login: string
    password: string
}

export class BootstrapAPI extends APIModuleBase {
    // Get bootstrap status
    async getStatus(): Promise<GetStatusResponse> {
        const response = await BootstrapAPI._axios.get<GetStatusResponse>('/bootstrap/status')
        return response.data
    }

    async createRootUser(credentials: CreateRootUserRequest): Promise<void> {
        await BootstrapAPI._axios.post('/bootstrap/rootUser', credentials)
    } 
}