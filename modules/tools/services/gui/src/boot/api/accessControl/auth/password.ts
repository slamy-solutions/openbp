import { APIModuleBase } from "../../model"

export interface StatusRequest {
    namespace: string
    identityUUID: string
}
export interface StatusResponse {
    seted: boolean
}

export interface DisableRequest {
    namespace: string
    identityUUID: string
}

export interface SetOrUpdateRequest {
    namespace: string
    identityUUID: string
    newPassword: string
}

export class PasswordAPI extends APIModuleBase {

    async status(params: StatusRequest): Promise<StatusResponse> {
        const response = await PasswordAPI._axios.get<StatusResponse>('/accessControl/iam/auth/password', { params })
        return response.data
    }

    async disable(params: DisableRequest): Promise<void> {
        await PasswordAPI._axios.delete('/accessControl/iam/auth/password', { params })
    }

    async setOrUpdate(params: SetOrUpdateRequest): Promise<void> {
        await PasswordAPI._axios.put('/accessControl/iam/auth/password', params)
    }
}