import { APIModuleBase } from '../../../model';
import { Device } from './models'

export interface ListRequest{
    namespace: string
    skip: number
    limit: number
    bindingFilter: 'all' | 'binded' | 'unbinded'
}
export interface ListResponse {
    devices: Array<Device>
    totalCount: number
}

export interface BindRequest {
    deviceNamespace: string
    deviceUuid: string
    balenaDeviceUUID: number
}
export interface BindResponse {}

export interface UnbindRequest {
    balenaDeviceUUID: string
}
export interface UnbindResponse {}

export class DeviceAPI extends APIModuleBase {

    async list(params: ListRequest): Promise<ListResponse> {
        const response = await DeviceAPI._axios.get<ListResponse>('/iot/integration/balena/devices', { params })
        const devices = response.data.devices
        devices.forEach(device => {
            device.created = new Date(device.created)
            device.updated = new Date(device.updated)

            device.balenaData.lastConnectivityEvent = new Date(device.balenaData.lastConnectivityEvent)
        })
        return response.data
    }

    async bind(params: BindRequest): Promise<BindResponse> {
        await DeviceAPI._axios.patch<BindResponse>('/iot/integration/balena/devices/bind', params)
        return {}
    }

    async unbind(params: UnbindRequest): Promise<UnbindResponse> {
        await DeviceAPI._axios.patch<BindResponse>('/iot/integration/balena/devices/unbind', params)
        return {}
    }
}
