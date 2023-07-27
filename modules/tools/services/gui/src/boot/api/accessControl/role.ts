import { APIModuleBase } from '../model'
import { NotManaged, IdentityManaged, ServiceManaged, BuiltInManaged } from './model'

export interface RolePolicy {
    namespace: string
    uuid: string
}

export interface Role {
    namespace: string
    uuid: string
    name: string
    description: string

    managed: NotManaged | IdentityManaged | ServiceManaged | BuiltInManaged

    policies: RolePolicy[]

    tags: string[]
    created: Date
    updated: Date
    version: number
}

export interface ListRequest {
    namespace: string
    skip?: number
    limit?: number
}
export interface ListResponse {
    roles: Role[]
    totalCount: number
}

export interface CreateRequest {
    namespace: string
    name: string
    description: string
}

export interface UpdateRequest {
    namespace: string
    uuid: string
    newName: string
    newDescription: string
}
export interface UpdateResponse {
    role: Role
}

export interface DeleteRequest {
    namespace: string
    uuid: string
}

export interface GetRequest {
    namespace: string
    uuid: string
}
export interface GetResponse {
    role: Role
}

export interface AddPolicyRequest {
    roleNamespace: string
    roleUUID: string
    policyNamespace: string
    policyUUID: string
}

export interface RemovePolicyRequest {
    roleNamespace: string
    roleUUID: string
    policyNamespace: string
    policyUUID: string
}

export class RoleAPI extends APIModuleBase {

    // Get list of roles
    async list(params: ListRequest): Promise<ListResponse> {
        const response = await RoleAPI._axios.get<ListResponse>('/accessControl/iam/role/list', { params })
        const roles = response.data.roles.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { roles, totalCount: response.data.totalCount } as ListResponse
    } 

    // Create new role
    async create(params: CreateRequest): Promise<Role> {
        const response = await RoleAPI._axios.post<{role: Role}>('/accessControl/iam/role', params)
        const role = response.data.role
        role.created = new Date(role.created) // we are receiving string in ISO format
        role.updated = new Date(role.updated)
        return role
    }

    async update(params: UpdateRequest): Promise<UpdateResponse> {
        const response = await RoleAPI._axios.patch<UpdateResponse>('/accessControl/iam/role', params)
        response.data.role.created = new Date(response.data.role.created) // we are receiving string in ISO format
        response.data.role.updated = new Date(response.data.role.updated)
        return response.data
    }

    // Delete role
    async delete(params: DeleteRequest): Promise<void> {
        await RoleAPI._axios.delete('/accessControl/iam/role', { params })
    }

    // Get role
    async get(params: GetRequest): Promise<Role> {
        const response = await RoleAPI._axios.get<GetResponse>('/accessControl/iam/role', { params })
        const role = response.data.role
        role.created = new Date(role.created) // we are receiving string in ISO format
        role.updated = new Date(role.updated)
        return role
    }

    // Add policy to the role. Do nothing if policy already added.
    async addPolicy(params: AddPolicyRequest): Promise<Role> {
        const response = await RoleAPI._axios.patch<{role: Role}>('/accessControl/iam/role/addPolicy', params)
        const role = response.data.role
        role.created = new Date(role.created) // we are receiving string in ISO format
        role.updated = new Date(role.updated)
        return role
    }

    // Remove policy from the role. Do nothing if policy already removed.
    async removePolicy(params: RemovePolicyRequest): Promise<Role> {
        const response = await RoleAPI._axios.patch<{role: Role}>('/accessControl/iam/role/removePolicy', params)
        const role = response.data.role
        role.created = new Date(role.created) // we are receiving string in ISO format
        role.updated = new Date(role.updated)
        return role
    }
}