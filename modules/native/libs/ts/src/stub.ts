import { Client, dummyClient, Rpc } from './client'
import { Config, makeDefaultConfig } from './config'

import { NamespaceServiceClientImpl } from './proto/namespace'
import { KeyValueStorageServiceClientImpl } from './proto/keyvaluestorage'

import { ActorUserServiceClientImpl } from './proto/actor/user'

import { IAMAuthenticationPasswordServiceClientImpl } from './proto/iam/authentication/password'
import { IAMIdentityServiceClientImpl } from './proto/iam/identity'
import { IAMPolicyServiceClientImpl } from './proto/iam/policy'
import { IAMOAuthServiceClientImpl } from './proto/iam/oauth'
import { IAMTokenServiceClientImpl } from './proto/iam/token'

interface NativeServices {
    namespace: NamespaceServiceClientImpl
    keyvaluestorage: KeyValueStorageServiceClientImpl
    actor: {
        user: ActorUserServiceClientImpl
    },
    iam: {
        authentication: {
            password: IAMAuthenticationPasswordServiceClientImpl
        },
        identity: IAMIdentityServiceClientImpl,
        policy: IAMPolicyServiceClientImpl,
        oauth: IAMOAuthServiceClientImpl,
        token: IAMTokenServiceClientImpl
    }
}

export enum Service {
    NAMESPACE,
    KEYVALUESTORAGE,
    ACTOR_USER,
    IAM_AUTHENTICAION_PASSWORD,
    IAM_IDENTITY,
    IAM_POLICY,
    IAM_OAUTH,
    IAM_TOKEN
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
        
        this.services = {
            namespace: new NamespaceServiceClientImpl(makeRpc(Service.NAMESPACE, config.urls.namespace)),
            keyvaluestorage: new KeyValueStorageServiceClientImpl(makeRpc(Service.KEYVALUESTORAGE, config.urls.keyvaluestorage)),
            actor: {
                user: new ActorUserServiceClientImpl(makeRpc(Service.ACTOR_USER, config.urls.actor.user))
            },
            iam: {
                authentication: {
                    password: new IAMAuthenticationPasswordServiceClientImpl(makeRpc(Service.IAM_AUTHENTICAION_PASSWORD, config.urls.iam.authentication.password))
                },
                identity: new IAMIdentityServiceClientImpl(makeRpc(Service.IAM_IDENTITY, config.urls.iam.identity)),
                oauth: new IAMOAuthServiceClientImpl(makeRpc(Service.IAM_OAUTH, config.urls.iam.oauth)),
                policy: new IAMPolicyServiceClientImpl(makeRpc(Service.IAM_POLICY, config.urls.iam.policy)),
                token: new IAMTokenServiceClientImpl(makeRpc(Service.IAM_TOKEN, config.urls.iam.token))
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