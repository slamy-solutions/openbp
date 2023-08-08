import { APIModuleBase } from '../model';
import { FleetAPI } from './fleet'
import { DeviceAPI } from './device'
import { IntegrationAPI } from './integration'

export class IoTAPI extends APIModuleBase {
    public readonly fleet: FleetAPI;
    public readonly device: DeviceAPI;
    public readonly integration: IntegrationAPI;

    constructor() {
        super();
        this.fleet = new FleetAPI()
        this.device = new DeviceAPI()
        this.integration = new IntegrationAPI()
    }
}