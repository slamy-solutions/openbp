import { APIModuleBase } from "../model"

export interface Device {
    namespace: string
    uuid: string
    name: string
    description: string
    identity: string

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
    device: Device
}

export interface DeleteRequest {
    namespace: string
    uuid: string
}

export class FleetAPI extends APIModuleBase {
    async createDevice(params: CreateRequest): Promise<CreateResponse> {
        const response = await FleetAPI._axios.post<CreateResponse>('/iot/devices/device', params)
        const device = response.data.device
        device.created = new Date(device.created)
        device.updated = new Date(device.updated)
        return { device }
    }

    async deleteDevice(params: CreateRequest): Promise<void> {
        await FleetAPI._axios.delete<CreateResponse>('/iot/devices/device', { params })
    }
}