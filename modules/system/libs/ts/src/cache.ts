import { RedisClientType, createClient } from 'redis';

export class Cache {
    public readonly redisClient: RedisClientType

    constructor(url: string) {
        this.redisClient = createClient({ url })
    }

    async connect() {}

    async close() {}
}