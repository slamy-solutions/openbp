import { LoginAPI } from './login'
import { BootstrapAPI } from './bootstrap'
import { NamespaceAPI } from './namespace'
import { AccessControlAPI } from './accessControl'
import { IoTAPI } from './iot'
import { MeAPI } from './me'
import { CRMAPI } from './crm'
import { ModulesAPI } from './modules'

export interface API {
    login: LoginAPI,
    bootstrap: BootstrapAPI,
    namespace: NamespaceAPI,
    accessControl: AccessControlAPI
    iot: IoTAPI
    me: MeAPI
    crm: CRMAPI
    modules: ModulesAPI
}

export const api = {
    login: new LoginAPI(),
    bootstrap: new BootstrapAPI(),
    namespace: new NamespaceAPI(),
    accessControl: new AccessControlAPI(),
    iot: new IoTAPI(),
    me: new MeAPI(),
    crm: new CRMAPI(),
    modules: new ModulesAPI()
} as API

export default api 