import { ServerAPI } from './server'
import { DeviceAPI } from './device'
import { SyncAPI } from './sync'
import { ToolsAPI } from './tools'


export class BalenaAPI {
    public readonly server: ServerAPI;
    public readonly device: DeviceAPI;
    public readonly sync: SyncAPI;
    public readonly tools: ToolsAPI;

    constructor() {
        this.server = new ServerAPI()
        this.device = new DeviceAPI()
        this.sync = new SyncAPI()
        this.tools = new ToolsAPI()
    }
}