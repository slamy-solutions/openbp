import { ClientAPI } from './client'
import { SettingsAPI } from './settings'

export class CRMAPI {
    public readonly settings: SettingsAPI;
    public readonly client: ClientAPI;

    constructor() {
        this.settings = new SettingsAPI()
        this.client = new ClientAPI()
    }
}