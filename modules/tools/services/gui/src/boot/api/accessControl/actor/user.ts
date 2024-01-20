import { APIModuleBase } from "../../model"

export interface User {
    namespace: string
    uuid: string
    login: string
    identity: string

    fullName: string
    avatar: string
    email: string

    created: Date
    updated: Date
    version: number
}

export interface CreateRequest {
    namespace: string
    login: string
    fullName: string
    email: string
}

export interface ListRequest {
    namespace: string
    skip: number
    limit: number
}
export interface ListResponse {
    users: Array<User>
    totalCount: number
}

export interface DeleteRequest {
    namespace: string
    uuid: string
}

export class UserAPI extends APIModuleBase {
    async create(params: CreateRequest): Promise<User> {
        const response = await UserAPI._axios.post<{user: User}>('/accessControl/iam/actor/user', params)
        const user = response.data.user
        user.created = new Date(user.created) // we are receiving string in ISO format
        user.updated = new Date(user.updated)
        return user
    }
    
    async list(params: ListRequest): Promise<ListResponse> {
        const response = await UserAPI._axios.get<ListResponse>('/accessControl/iam/actor/user', { params })
        const users = response.data.users.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        
        return {users, totalCount: response.data.totalCount}
    }

    async delete(params: DeleteRequest): Promise<void> {
        await UserAPI._axios.delete('/accessControl/iam/actor/user', { params })
    }
}