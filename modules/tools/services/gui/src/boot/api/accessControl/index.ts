import { APIModuleBase } from '../model';
import { IdentityAPI } from './identity'
import { RoleAPI } from './role';
import { PolicyAPI } from './policy'

import { UserAPI } from './actor/user'

import { PasswordAPI } from './auth/password'
import { CertificateAPI } from './auth/certificate'

export interface ActorAPIs {
    user: UserAPI
}

export interface AuthAPIs {
    password: PasswordAPI
    certificate: CertificateAPI
}

export class AccessControlAPI extends APIModuleBase {
    public identity: IdentityAPI;
    public role: RoleAPI
    public policy: PolicyAPI
    public actor: ActorAPIs
    public auth: AuthAPIs

    constructor() {
        super();
        this.identity = new IdentityAPI()
        this.role = new RoleAPI()
        this.policy = new PolicyAPI()
        this.actor = {
            user: new UserAPI()
        }
        this.auth = {
            password: new PasswordAPI(),
            certificate: new CertificateAPI()
        }
    }
}