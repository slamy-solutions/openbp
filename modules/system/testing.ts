import { Testing } from '../scheme'

import { connect as connectToMongo, close as closeMongoConnection } from './testing/tools/mongo'
import { connect as connectToCache, close as closeCacheConnection } from './testing/tools/cache'

export default {
    async setup () {
        await connectToMongo()
        await connectToCache()
    },
    async teardown () {
        await closeMongoConnection()
        await closeCacheConnection()
    }
} as Testing