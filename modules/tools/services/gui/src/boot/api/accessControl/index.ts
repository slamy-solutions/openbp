import { APIModuleBase } from '../model';
import { IdentityAPI } from './identity'
import { RoleAPI } from './role';
import { PolicyAPI } from './policy'


export class AccessControlAPI extends APIModuleBase {
    public identity: IdentityAPI;
    public role: RoleAPI
    public policy: PolicyAPI

    constructor() {
        super();
        this.identity = new IdentityAPI()
        this.role = new RoleAPI()
        this.policy = new PolicyAPI()
    }
}