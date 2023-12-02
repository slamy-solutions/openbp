import { APIModuleBase } from '../model'

export interface Department {
    namespace: string
    uuid: string
    name: string
}

export interface CreateDepartmentRequest {
    namespace: string
    name: string
}
export interface CreateDepartmentResponse {
    department: Department
}

export interface GetAllDepartmentsRequest {
    namespace: string
}
export interface GetAllDepartmentsResponse {
    departments: Array<Department>
}

export interface GetDepartmentRequest {
    namespace: string
    uuid: string
}
export interface GetDepartmentResponse {
    department: Department
}

export interface UpdateDepartmentRequest {
    namespace: string
    uuid: string
    name: string
}
export interface UpdateDepartmentResponse {
    department: Department
}

export interface DeleteDepartmentRequest {
    namespace: string
    uuid: string
}
export interface DeleteDepartmentResponse {}

export class DepartmentAPI extends APIModuleBase {
    async create(params: CreateDepartmentRequest): Promise<CreateDepartmentResponse> {
        const response = await DepartmentAPI._axios.post<CreateDepartmentResponse>('/crm/departments/department', params)
        return response.data
    }

    async getAll(params: GetAllDepartmentsRequest): Promise<GetAllDepartmentsResponse> {
        const response = await DepartmentAPI._axios.get<GetAllDepartmentsResponse>('/crm/departments', { params })
        return response.data
    }

    async get(params: GetDepartmentRequest): Promise<GetDepartmentResponse> {
        const response = await DepartmentAPI._axios.get<GetDepartmentResponse>('/crm/departments/department', { params })
        return response.data
    }

    async update(params: UpdateDepartmentRequest): Promise<UpdateDepartmentResponse> {
        const response = await DepartmentAPI._axios.patch<UpdateDepartmentResponse>('/crm/departments/department', params)
        return response.data
    }

    async delete(params: DeleteDepartmentRequest): Promise<DeleteDepartmentResponse> {
        const response = await DepartmentAPI._axios.delete<DeleteDepartmentResponse>('/crm/departments/department', { params })
        return response.data
    }
}