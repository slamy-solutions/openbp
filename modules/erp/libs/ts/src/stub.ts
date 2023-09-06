import { Client, dummyClient, Rpc } from './client'
import { Config, makeDefaultConfig } from './config'

import { CatalogServiceClientImpl, CatalogEntryServiceClientImpl, CatalogIndexServiceClientImpl } from './proto/core/catalog'

interface ERPServices {
    core: {
        catalog: {
            catalog: CatalogServiceClientImpl,
            entry: CatalogEntryServiceClientImpl,
            index: CatalogIndexServiceClientImpl
        }
    }
}

export enum Service {CORE}

export class ERPStub {
    private clients: Array<Client>
    public services: ERPServices

    constructor(connectTo: Array<Service> = [Service.CORE], config?: Config) {
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
        
        const coreConnection = makeRpc(Service.CORE, config.urls.core)

        this.services = {
            core: {
                catalog: {
                    catalog: new CatalogServiceClientImpl(coreConnection),
                    entry: new CatalogEntryServiceClientImpl(coreConnection),
                    index: new CatalogIndexServiceClientImpl(coreConnection)
                }
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