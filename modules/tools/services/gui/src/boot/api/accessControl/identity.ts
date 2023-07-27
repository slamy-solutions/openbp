import { APIModuleBase } from "../model"
import { NotManaged, IdentityManaged, ServiceManaged } from './model'

export interface IdentityRole {
    namespace: string
    uuid: string
}

export interface IdentityPolicy {
    namespace: string
    uuid: string
}

export interface Identity {
    namespace: string
    uuid: string
    name: string
    active: boolean

    managed: NotManaged | IdentityManaged | ServiceManaged

    roles: IdentityRole[]
    policies: IdentityPolicy[]

    created: Date
    updated: Date
    version: number
}

export interface ListRequest {
    namespace: string,
    skip?: number,
    limit?: number
}
export interface ListResponse {
    identities: Array<Identity>
    totalCount: number
}

export interface CreateRequest {
    namespace: string
    name: string
    initiallyActive: boolean
}

export interface GetRequest {
    namespace: string
    uuid: string
}

export interface UpdateRequest {
    namespace: string
    uuid: string
    newName: string
}

export interface AddPolicyRequest {
    identityNamespace: string
    identityUUID: string
    policyNamespace: string
    policyUUID: string
}
export interface RemovePolicyRequest {
    identityNamespace: string
    identityUUID: string
    policyNamespace: string
    policyUUID: string
}

export interface AddRoleRequest {
    identityNamespace: string
    identityUUID: string
    roleNamespace: string
    roleUUID: string
}
export interface RemoveRoleRequest {
    identityNamespace: string
    identityUUID: string
    roleNamespace: string
    roleUUID: string
}

export interface DeleteRequest {
    namespace: string
    uuid: string
}

export interface SetActiveRequest {
    namespace: string
    uuid: string
    active: boolean
}


export class IdentityAPI extends APIModuleBase {

    // Get list of identities
    async list(params: ListRequest): Promise<ListResponse> {
        const response = await IdentityAPI._axios.get<ListResponse>('/accessControl/iam/identity/list', { params })
        const identities = response.data.identities.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { identities, totalCount: response.data.totalCount } as ListResponse
    } 

    // Create new identity
    async create(params: CreateRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.post<{identity: Identity}>('/accessControl/iam/identity', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async get(params: GetRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.get<{identity: Identity}>('/accessControl/iam/identity', { params } )
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async update(params: UpdateRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async setActive(params: SetActiveRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity/active', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async addPolicy(params: AddPolicyRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity/addPolicy', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async removePolicy(params: RemovePolicyRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity/removePolicy', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async addRole(params: AddRoleRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity/addRole', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    async removeRole(params: RemoveRoleRequest): Promise<Identity> {
        const response = await IdentityAPI._axios.patch<{identity: Identity}>('/accessControl/iam/identity/removeRole', params)
        const identity = response.data.identity
        identity.created = new Date(identity.created) // we are receiving string in ISO format
        identity.updated = new Date(identity.updated)
        return identity
    }

    // Delete identity
    async delete(params: DeleteRequest): Promise<void> {
        await IdentityAPI._axios.delete<DeleteRequest>('/accessControl/iam/identity', { params })
    }
}