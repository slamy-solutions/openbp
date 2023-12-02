import { APIModuleBase } from '../model'

export interface Project {
    namespace: string
    uuid: string
    name: string

    clientUUID: string
    contactUUID: string
    departmentUUID: string

    notRelevant: boolean
}

export interface CreateProjectRequest {
    namespace: string
    name: string
    clientUUID: string
    contactUUID: string
    departmentUUID: string
}
export interface CreateProjectResponse {
    project: Project
}

export interface GetProjectsRequest {
    namespace: string
    clientUUID: string
    departmentUUID: string
}
export interface GetProjectsResponse {
    projects: Array<Project>
}

export interface GetProjectRequest {
    namespace: string
    uuid: string
}
export interface GetProjectResponse {
    project: Project
}

export interface UpdateProjectRequest {
    namespace: string
    uuid: string
    name: string
    clientUUID: string
    contactUUID: string
    departmentUUID: string
}
export interface UpdateProjectResponse {
    project: Project
}

export interface DeleteProjectRequest {
    namespace: string
    uuid: string
}
export interface DeleteProjectResponse {}

export class ProjectAPI extends APIModuleBase {
    async create(params: CreateProjectRequest): Promise<CreateProjectResponse> {
        const response = await ProjectAPI._axios.post<CreateProjectResponse>('/crm/projects/project', params)
        return response.data
    }

    async get(params: GetProjectRequest): Promise<GetProjectResponse> {
        const response = await ProjectAPI._axios.get<GetProjectResponse>('/crm/projects/project', { params })
        return response.data
    }

    async getAll(params: GetProjectsRequest): Promise<GetProjectsResponse> {
        const response = await ProjectAPI._axios.get<GetProjectsResponse>('/crm/projects', { params })
        return response.data
    }

    async update(params: UpdateProjectRequest): Promise<UpdateProjectResponse> {
        const response = await ProjectAPI._axios.patch<UpdateProjectResponse>('/crm/projects/project', params)
        return response.data
    }

    async delete(params: DeleteProjectRequest): Promise<DeleteProjectResponse> {
        await ProjectAPI._axios.delete('/crm/projects/project', { params })
        return {}
    }
}