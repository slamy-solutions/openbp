import { LoginAPI } from './login'
import { BootstrapAPI } from './bootstrap'
import { NamespaceAPI } from './namespace'
import { AccessControlAPI } from './accessControl'

export interface API {
    login: LoginAPI,
    bootstrap: BootstrapAPI,
    namespace: NamespaceAPI,
    accessControl: AccessControlAPI
}

export const api = {
    login: new LoginAPI(),
    bootstrap: new BootstrapAPI(),
    namespace: new NamespaceAPI(),
    accessControl: new AccessControlAPI()
} as API

export default api 