export interface Config {
    vault: {
        url: string
    }
    redis: {
        url: string
    }
    nats: {
        url: string
    }
    cache: {
        url: string
    }
    db: {
        url: string
    }
    telemetry: {
        url: string
    }
}

export const makeDefaultConfig: ()=>Config = () => {
    return {
        vault: {
            url: process.env.SYSTEM_VAULT_URL || "system_vault:80"
        },
        redis: {
            url: process.env.SYSTEM_REDIS_URL || "redis://system_redis"
        },
        nats: {
            url: process.env.SYSTEM_NATS_URL || "nats://system_nats:4222"
        },
        cache: {
            url: process.env.SYSTEM_CACHE_URL || "redis://system_redis"
        },
        db: {
            url: process.env.SYSTEM_DB_URL || "mongodb://root:example@system_db/admin"
        },
        telemetry: {
            url: process.env.SYSTEM_TELEMETRY_EXPORTER_ENDPOINT || "system_telemetry:55680"
        }
    }
}