import { AuthAPI } from "./auth"
import { UserAPI } from './user'

export class MeAPI {
    public readonly auth: AuthAPI
    public readonly user: UserAPI

    constructor() {
        this.auth = new AuthAPI()
        this.user = new UserAPI()
    }
}