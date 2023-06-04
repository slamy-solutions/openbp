import { APIModuleBase } from "../model"
import { Namespace } from './models'


export interface CreateRequest {
    name: string
    fullName: string
    description: string
}

export interface DeleteRequest {
    name: string
}


export class ListAPI extends APIModuleBase {

    // Get list of namespaces
    async list(): Promise<Namespace[]> {
        const response = await ListAPI._axios.get<{namespaces: Array<Namespace>}>('/namespace/list')
        return response.data.namespaces
    } 

    // Create new namespace
    async create(params: CreateRequest): Promise<Namespace> {
        const response = await ListAPI._axios.post<{namespace: Namespace}>('/namespace/list/namespace', params)
        return response.data.namespace
    }

    // Delete namespace
    async delete(params: DeleteRequest): Promise<void> {
        await ListAPI._axios.delete<DeleteRequest>('/namespace/list/namespace', { params })
    }
}