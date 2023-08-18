import { Client, dummyClient, Rpc } from './client'
import { Config, makeDefaultConfig } from './config'

import { NamespaceServiceClientImpl } from './proto/namespace'
import { KeyValueStorageServiceClientImpl } from './proto/keyvaluestorage'

import { IAMAuthenticationPasswordServiceClientImpl } from './proto/iam/authentication/password'
import { IAMIdentityServiceClientImpl } from './proto/iam/identity'
import { IAMPolicyServiceClientImpl } from './proto/iam/policy'
import { IAMTokenServiceClientImpl } from './proto/iam/token'
import { ActorUserServiceClientImpl } from './proto/iam/actor/user'
import { IAMAuthenticationX509ServiceClientImpl } from './proto/iam/authentication/x509'
import { IAMRoleServiceClientImpl } from './proto/iam/role'
import { IAMAuthServiceClientImpl } from './proto/iam/auth'

interface NativeServices {
    namespace: NamespaceServiceClientImpl
    keyvaluestorage: KeyValueStorageServiceClientImpl
    iam: {
        actor: {
            user: ActorUserServiceClientImpl
        },
        authentication: {
            password: IAMAuthenticationPasswordServiceClientImpl,
            x509: IAMAuthenticationX509ServiceClientImpl
        },
        identity: IAMIdentityServiceClientImpl,
        policy: IAMPolicyServiceClientImpl,
        role: IAMRoleServiceClientImpl,
        token: IAMTokenServiceClientImpl,
        auth: IAMAuthServiceClientImpl
    }
}

export enum Service {
    NAMESPACE,
    KEYVALUESTORAGE,
    IAM
}

export class NativeStub {
    private clients: Array<Client>
    public services: NativeServices

    constructor(connectTo: Array<Service>, config?: Config) {
        if (config === undefined) {
            config = makeDefaultConfig()
        }
        this.clients = []
        const clientByUrl = new Map<string, Client>()
        
        const connectToSet = new Set(connectTo)
        
        const grpcConfig = {
            'grpc.service_config': JSON.stringify({ loadBalancingConfig: [{ round_robin: {} }], })
        }
        const makeRpc = (service: Service, url: string) => {
            if (connectToSet.has(service)) {
                const savedCliet = clientByUrl.get(url)
                if (savedCliet !== undefined) return savedCliet
                
                const client = new Client(url, undefined, grpcConfig)
                this.clients.push(client)
                clientByUrl.set(url, client)
                return client
            }
            return dummyClient
        }
        
        const iamConnection = makeRpc(Service.IAM, config.urls.iam)

        this.services = {
            namespace: new NamespaceServiceClientImpl(makeRpc(Service.NAMESPACE, config.urls.namespace)),
            keyvaluestorage: new KeyValueStorageServiceClientImpl(makeRpc(Service.KEYVALUESTORAGE, config.urls.keyvaluestorage)),
            iam: {
                actor: {
                    user: new ActorUserServiceClientImpl(iamConnection)
                },
                authentication: {
                    password: new IAMAuthenticationPasswordServiceClientImpl(iamConnection),
                    x509: new IAMAuthenticationX509ServiceClientImpl(iamConnection)
                },
                identity: new IAMIdentityServiceClientImpl(iamConnection),
                policy: new IAMPolicyServiceClientImpl(iamConnection),
                role: new IAMRoleServiceClientImpl(iamConnection),
                token: new IAMTokenServiceClientImpl(iamConnection),
                auth: new IAMAuthServiceClientImpl(iamConnection)
            }
        }
    }

    async connect() {
        for(const client of this.clients) {
            await client.connect()
        }
    }

    close() {
        for(const client of this.clients) {
            client.close()
        }
    }
}