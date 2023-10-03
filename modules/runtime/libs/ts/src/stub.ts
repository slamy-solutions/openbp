import { Client, dummyClient, Rpc } from './client'
import { Config, makeDefaultConfig } from './config'

import { RuntimeServiceClientImpl } from './proto/manager/runtime'
import { RPCServiceClientImpl } from './proto/manager/rpc'

interface RuntimeServices {
    manager: {
        runtime: RuntimeServiceClientImpl,
        rpc: RPCServiceClientImpl
    }
}

export enum Service {MANAGER}

export class RuntimeStub {
    private clients: Array<Client>
    public services: RuntimeServices

    constructor(connectTo: Array<Service> = [Service.MANAGER], config?: Config) {
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
        
        const managerConnection = makeRpc(Service.MANAGER, config.urls.manager)

        this.services = {
            manager: {
                runtime: new RuntimeServiceClientImpl(managerConnection),
                rpc: new RPCServiceClientImpl(managerConnection)
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