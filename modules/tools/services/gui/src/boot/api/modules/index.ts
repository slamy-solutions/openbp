import { APIModuleBase } from "../model"

export interface GetModulesStatusResponse {
    system: boolean
    native: boolean
    tools: boolean
    iot: boolean
    crm: boolean
    erp: boolean
}

export class ModulesAPI extends APIModuleBase {
    async getStatus(): Promise<GetModulesStatusResponse> {
        const response = await ModulesAPI._axios.get<GetModulesStatusResponse>('/modules/status')
        return response.data
    }
}