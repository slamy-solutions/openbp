import { createClient, RedisClientType } from 'redis'
import { connect as connectToNats, NatsConnection } from 'nats'
import { MongoClient } from 'mongodb'
import { Config, makeDefaultConfig } from './config'
import { Client as GRPCClient } from './client'
import { VaultServiceClientImpl } from './proto/vault'
import { Cache } from './cache'

export interface SystemServices {
    redis: RedisClientType
    cache: Cache
    db: MongoClient
    vault: VaultServiceClientImpl
    nats: NatsConnection
}


export enum Service {
    VAULT,
    NATS,
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
            vault: new VaultServiceClientImpl(client),
            nats: null as unknown as NatsConnection,
            cache: new Cache(config.cache.url)
        }
    }

    async connect() {
        if (this.loadedServices.has(Service.NATS)) {
            const url = new URL(this.config.nats.url)
            const server = url.hostname + (url.port !== "") ? ":" + url.port : ""
            this.services.nats = await connectToNats({
                servers: server
            })
        }
        if (this.loadedServices.has(Service.REDIS)) {
            await this.services.redis.connect()
        }
        if (this.loadedServices.has(Service.CACHE)) {
            await this.services.cache.connect()
        }
        if (this.loadedServices.has(Service.DB)) {
            await this.services.db.connect()
        }
        for (const client of this.grpcClients) {
            await client.connect()
        }
    }

    async close() {
        if (this.loadedServices.has(Service.NATS)) {
            await this.services.nats.drain()
        }
        if (this.loadedServices.has(Service.REDIS)) {
            await this.services.redis.disconnect()
        }
        if (this.loadedServices.has(Service.CACHE)) {
            await this.services.cache.close()
        }
        if (this.loadedServices.has(Service.DB)) {
            await this.services.db.close()
        }
        for (const client of this.grpcClients) {
            await client.close()
        }
    }
}