import { APIModuleBase } from '../model';
import { FleetAPI } from './fleet'

export class IoTAPI extends APIModuleBase {
    public fleet: FleetAPI;

    constructor() {
        super();
        this.fleet = new FleetAPI()
    }
}