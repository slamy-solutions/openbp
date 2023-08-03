export interface Server {
    namespace: string
    uuid: string
    name: string
    description: string
    enabled: boolean

    created: Date
    updated: Date
    version: number
}

export interface DeviceData {
    uuid: string
    id: number
    isOnline: boolean
    status: string
    deviceName: string
    longitude: string
    latitude: string
    location: string
    lastConnectivityEvent: Date
    memoryUsage: number
    menoryTotal: number
    storageUsage: number
    cpuUsage: number
    cpuTemp: number
    isUndervolted: boolean
}

export interface Device {
    uuid: string
    bindedDeviceNamespace : string
    bindedDeviceUUID: string

    balenaServerNamespace: string
    balenaServerUUID: string
    balenaData: DeviceData

    created: Date
    updated: Date
    version: number
}

export interface SyncLogEntryStats {
    foundedDevicesOnServer: number
    foundedActiveDevices: number
    metricsUpdates: number
    executionTime: number
}

export interface SyncLogEntry {
    uuid: string
    serverUUID: string
    timestamp: Date
    status: string
    error: string
    stats: SyncLogEntryStats
}