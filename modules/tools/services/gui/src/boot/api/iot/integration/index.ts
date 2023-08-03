import { BalenaAPI } from './balena'

export class IntegrationAPI {
    public readonly balena: BalenaAPI;

    constructor() {
        this.balena = new BalenaAPI()
    }
}