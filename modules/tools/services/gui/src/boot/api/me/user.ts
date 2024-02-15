import { APIModuleBase } from "../model"
import { User } from '../accessControl/actor/user'

export interface GetUserInfoResponse {
    user: User
}

export class UserAPI extends APIModuleBase {
    async getInfo(): Promise<GetUserInfoResponse> {
        const response = await UserAPI._axios.get<GetUserInfoResponse>('/me/user')
        return response.data
    }
}