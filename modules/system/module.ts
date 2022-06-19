import { Module } from '../scheme'

declare global {
    namespace NodeJS {
        interface ProcessEnv {
            SYSTEM_DB_URL: string
            SYSTEM_DB_PREFIX: string
            SYSTEM_CACHE_URL: string
            SYSTEM_BIGCACHE_URL: string
            SYSTEM_TELEMETRY_EXPORTER_ENDPOINT: string
        }
    }
}

process.env.SYSTEM_DB_URL = process.env.SYSTEM_DB_URL || "mongodb://root:example@system_db/admin"
process.env.SYSTEM_DB_PREFIX = process.env.SYSTEM_DB_PREFIX || "openerp_"
process.env.SYSTEM_CACHE_URL = process.env.SYSTEM_CACHE_URL || "redis://system_cache"
process.env.SYSTEM_BIGCACHE_URL = process.env.SYSTEM_BIGCACHE_URL || "redis://system_bigcache"
process.env.SYSTEM_TELEMETRY_EXPORTER_ENDPOINT = process.env.SYSTEM_TELEMETRY_EXPORTER_ENDPOINT || "system_telemetry:55680"

export default {
    name: 'System',
    uuid: 'system'
} as Module