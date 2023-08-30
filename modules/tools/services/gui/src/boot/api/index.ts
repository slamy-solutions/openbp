import { LoginAPI } from './login'
import { BootstrapAPI } from './bootstrap'
import { NamespaceAPI } from './namespace'
import { AccessControlAPI } from './accessControl'
import { IoTAPI } from './iot'
import { MeAPI } from './me'

export interface API {
    login: LoginAPI,
    bootstrap: BootstrapAPI,
    namespace: NamespaceAPI,
    accessControl: AccessControlAPI
    iot: IoTAPI
    me: MeAPI
}

export const api = {
    login: new LoginAPI(),
    bootstrap: new BootstrapAPI(),
    namespace: new NamespaceAPI(),
    accessControl: new AccessControlAPI(),
    iot: new IoTAPI(),
    me: new MeAPI()
} as API

export default api 