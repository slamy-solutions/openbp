import { OAuthAPI } from './oauth'

export class ConfigAPI {
    public readonly oauth: OAuthAPI

    constructor() {
        this.oauth = new OAuthAPI()
    }
}