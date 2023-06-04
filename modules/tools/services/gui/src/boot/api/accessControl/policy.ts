import { APIModuleBase } from '../model'
import { NotManaged, IdentityManaged, ServiceManaged, BuiltInManaged } from './model'

export interface Policy {
    namespace: string
    uuid: string
    name: string
    description: string

    managed: NotManaged | IdentityManaged | ServiceManaged | BuiltInManaged

    resources: string[]
    actions: string[]
    namespaceIndependent: boolean

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
    policies: Policy[]
    totalCount: number
}

export interface CreateRequest {
    namespace: string
    name: string
    description: string
    resources: string[]
    actions: string[]
    namespaceIndependent: boolean
}

export interface UpdateRequest {
    namespace: string
    uuid: string
    name: string
    description: string
    resources: string[]
    actions: string[]
    namespaceIndependent: boolean
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
    policy: Policy
}

export class PolicyAPI extends APIModuleBase {

    // Get list of policies
    async list(params: ListRequest): Promise<ListResponse> {
        const response = await PolicyAPI._axios.get<ListResponse>('/accessControl/iam/policy/list', { params })
        const policies = response.data.policies.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { policies, totalCount: response.data.totalCount } as ListResponse
    } 

    // Create new policy
    async create(params: CreateRequest): Promise<Policy> {
        const response = await PolicyAPI._axios.post<{policy: Policy}>('/accessControl/iam/policy', params)
        const policy = response.data.policy
        policy.created = new Date(policy.created) // we are receiving string in ISO format
        policy.updated = new Date(policy.updated)
        return policy
    }

    // Update policy and get updated version
    async update(params: UpdateRequest): Promise<Policy> {
        const response = await PolicyAPI._axios.patch<{policy: Policy}>('/accessControl/iam/policy', params)
        const policy = response.data.policy
        policy.created = new Date(policy.created) // we are receiving string in ISO format
        policy.updated = new Date(policy.updated)
        return policy
    }

    // Delete policy
    async delete(params: DeleteRequest): Promise<void> {
        await PolicyAPI._axios.delete('/accessControl/iam/policy', { params })
    }

    // Get policy
    async get(params: GetRequest): Promise<Policy> {
        const response = await PolicyAPI._axios.get<GetResponse>('/accessControl/iam/policy', { params })
        const policy = response.data.policy
        policy.created = new Date(policy.created) // we are receiving string in ISO format
        policy.updated = new Date(policy.updated)
        return policy
    }
}