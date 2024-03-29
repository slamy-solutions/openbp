import { APIModuleBase } from '../model';
import { IdentityAPI } from './identity'
import { RoleAPI } from './role';
import { PolicyAPI } from './policy'

import { UserAPI } from './actor/user'

import { PasswordAPI } from './auth/password'
import { CertificateAPI } from './auth/certificate'

import { ConfigAPI } from './config'

export interface ActorAPIs {
    user: UserAPI
}

export interface AuthAPIs {
    password: PasswordAPI
    certificate: CertificateAPI
}

export class AccessControlAPI extends APIModuleBase {
    public readonly identity: IdentityAPI;
    public readonly role: RoleAPI
    public readonly policy: PolicyAPI
    public readonly actor: ActorAPIs
    public readonly auth: AuthAPIs
    public readonly config: ConfigAPI

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
        this.config = new ConfigAPI()
    }
}