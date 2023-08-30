import { AuthAPI } from "./auth"

export class MeAPI {
    public readonly auth: AuthAPI

    constructor() {
        this.auth = new AuthAPI()
    }
}