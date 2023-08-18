import { createClient, RedisClientType } from 'redis'
import { MongoClient } from 'mongodb'
import { Config, makeDefaultConfig } from './config'
import { Client as GRPCClient } from './client'
import { VaultServiceClientImpl } from './proto/vault'

export interface SystemServices {
    redis: RedisClientType
    db: MongoClient
    vault: VaultServiceClientImpl
}


export enum Service {
    VAULT,
    REDIS,
    CACHE,
    DB,
    TELEMETRY
}

export class SystemStub {
    readonly services: SystemServices
    private readonly config: Config
    private readonly loadedServices: Set<Service>

    private readonly grpcClients: Array<GRPCClient>

    constructor(services: Array<Service>, config?: Config) {
        if (config === undefined) {
            config = makeDefaultConfig()
        }

        this.loadedServices = new Set(services)
        this.config = config
        
        this.grpcClients = []
        const grpcConfig = {
            'grpc.service_config': JSON.stringify({ loadBalancingConfig: [{ round_robin: {} }], })
        }
        const client = new GRPCClient(this.config.vault.url, undefined, grpcConfig)
        this.grpcClients.push(client)

        this.services = {
            redis: createClient({
                url: config.redis.url
            }),
            db: new MongoClient(config.db.url),
            vault: new VaultServiceClientImpl(client)
        }
    }

    async connect() {
        if (this.loadedServices.has(Service.REDIS)) {
            await this.services.redis.connect()
        }
        if (this.loadedServices.has(Service.DB)) {
            await this.services.db.connect()
        }
        for (const client of this.grpcClients) {
            await client.connect()
        }
    }

    async close() {
        if (this.loadedServices.has(Service.REDIS)) {
            await this.services.redis.disconnect()
        }
        if (this.loadedServices.has(Service.DB)) {
            await this.services.db.close()
        }
        for (const client of this.grpcClients) {
            await client.close()
        }
    }
}