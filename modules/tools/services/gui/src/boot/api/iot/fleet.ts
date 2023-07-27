import { APIModuleBase } from '../model'

export interface Fleet {
    namespace: string
    uuid: string
    name: string
    description: string

    created: Date
    updated: Date
    version: number
}

export interface ListFleetsRequest {
    namespace: string
    skip: number
    limit: number
}

export interface ListFleetsResponse {
    fleets: Array<Fleet>
    totalCount: number
}

export class FleetAPI extends APIModuleBase {

    async listFleets(params: ListFleetsRequest): Promise<ListFleetsResponse> {
        const response = await FleetAPI._axios.get<ListFleetsResponse>('/iot/fleets', { params })
        const fleets = response.data.fleets.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { fleets, totalCount: response.data.totalCount } as ListFleetsResponse
    }
}