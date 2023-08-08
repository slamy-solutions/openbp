import { APIModuleBase } from '../model'
import { Device } from './device'

export interface Fleet {
    namespace: string
    uuid: string
    name: string
    description: string

    created: Date
    updated: Date
    version: number
}

export interface CreateRequest {
    namespace: string
    name: string
    description: string
}
export interface CreateResponse {
    fleet: Fleet
}

export interface GetRequest {
    namespace: string
    uuid: string
}
export interface GetResponse {
    fleet: Fleet
}

export interface UpdateRequest {
    namespace: string
    uuid: string
    newDescription: string
}
export interface UpdateResponse {
    fleet: Fleet
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

export interface DeleteRequest {
    namespace: string
    uuid: string
}

export interface ListDevicesRequest {
    namespace: string
    uuid: string
    skip: number
    limit: number
}
export interface ListDevicesResponse {
    devices: Array<Device>
    totalCount: number
}

export interface AddDeviceRequest {
    namespace: string
    deviceUUID: string
    fleetUUID: string
}
export interface RemoveDeviceRequest {
    namespace: string
    deviceUUID: string
    fleetUUID: string
}

export class FleetAPI extends APIModuleBase {

    async createFleet(params: CreateRequest): Promise<CreateResponse> {
        const response = await FleetAPI._axios.post<CreateResponse>('/iot/fleets/fleet', params)
        const fleet = response.data.fleet
        fleet.created = new Date(fleet.created)
        fleet.updated = new Date(fleet.updated)
        return { fleet }
    }

    async listFleets(params: ListFleetsRequest): Promise<ListFleetsResponse> {
        const response = await FleetAPI._axios.get<ListFleetsResponse>('/iot/fleets', { params })
        const fleets = response.data.fleets.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { fleets, totalCount: response.data.totalCount } as ListFleetsResponse
    }

    async getFleet(params: GetRequest): Promise<GetResponse> {
        const response = await FleetAPI._axios.get<GetResponse>('/iot/fleets/fleet', { params })
        const fleet = response.data.fleet
        fleet.created = new Date(fleet.created)
        fleet.updated = new Date(fleet.updated)
        return { fleet }
    }

    async updateFleet(params: UpdateRequest): Promise<UpdateResponse> {
        const response = await FleetAPI._axios.patch<UpdateResponse>('/iot/fleets/fleet', params)
        const fleet = response.data.fleet
        fleet.created = new Date(fleet.created)
        fleet.updated = new Date(fleet.updated)
        return { fleet }
    }

    async deleteFleet(params: DeleteRequest): Promise<void> {
        await FleetAPI._axios.delete<CreateResponse>('/iot/fleets/fleet', { params })
    }

    async listDevices(params: ListDevicesRequest): Promise<ListDevicesResponse> {
        const response = await FleetAPI._axios.get<ListDevicesResponse>('/iot/fleets/fleet/devices', { params })
        const devices = response.data.devices.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        return { devices, totalCount: response.data.totalCount } as ListDevicesResponse
    }

    async addDevice(params: AddDeviceRequest): Promise<void> {
        await FleetAPI._axios.put('/iot/fleets/fleet/devices/device', params )
    }

    async removeDevice(params: RemoveDeviceRequest): Promise<void> {
        await FleetAPI._axios.delete('/iot/fleets/fleet/devices/device', { params } )
    }
}