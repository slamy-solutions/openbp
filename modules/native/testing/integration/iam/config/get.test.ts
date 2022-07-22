
import { Status } from '@grpc/grpc-js/build/src/constants'
import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { configClient as iamConfigGRPC, connect as connectToNativeIAMConfig, close as closeNativeIAMConfig } from '../../../tools/iam/grpc'
import { client as nativeKeyValueStorageClient, connect as connectToNativeKeyValueStorage, close as closeNativeKeyValueStorage } from '../../../tools/keyvaluestorage/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeIAMConfig()
    await connectToNativeKeyValueStorage()
})

afterEach(async ()=>{
    await mongoClient.db(GLOBAL_DB_NAME).collection('native_keyvaluestorage').deleteMany({})
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeIAMConfig()
    await closeNativeKeyValueStorage()
})

/**
 * @group native/iam/config/set/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Gets default value when config is not set", async () => {
        const config = await iamConfigGRPC.Get({ useCache: false })
        expect(config.configuration).not.toBeNull()
        expect(config.configuration).not.toBeUndefined()

        try {
            await nativeKeyValueStorageClient.Get({ namespace: "", key: "native_iam_config_global", useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})