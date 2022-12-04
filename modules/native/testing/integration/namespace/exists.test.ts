import { Status } from '@grpc/grpc-js/build/src/constants'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'
import { client as grpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../tools/namespace/grpc'

const GLOBAL_DB_NAME = 'openbp_global'

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_namespace').deleteMany({})
    } catch (e) {
        if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) throw e
    }
    await cacheClient.flushall()
})

afterEach(async ()=>{
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_namespace').deleteMany({})
    } catch (e) {
        if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) throw e
    }
    await cacheClient.flushall()
})

afterAll(async () => {
    await closeNativeNamespace()
    await closeCache()
    await closeMongo()
})

/**
 * @group native/namespace/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Adds value to cache, when cache enabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Exists({ name: name, useCache: true })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
    })
    test("Doesnt add value to cache, when cache disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Exists({ name: name, useCache: false })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })
    test("Gets value from cache, when cache enabled. Doesnt get value from cache, when cache disabled.", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })

        // Change value in database to contain wrong data without invalidating cache
        await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('native_namespace').deleteOne({ name })

        const cachedResponse = await grpc.Exists({ name, useCache: true })
        expect(cachedResponse.exist).toBeTruthy()
        const uncachedResponse = await grpc.Exists({ name, useCache: false })
        expect(uncachedResponse.exist).toBeFalsy()
    })
})

/**
 * @group native/namescape/get/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
     test("Returs actual namespace status (cache disabled)", async () => {
        const name = "testname" 
        let response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeFalsy()
        await grpc.Ensure({ name })
        response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeTruthy()
        await grpc.Delete({ name })
        response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeFalsy()
     })

     test("Returs actual namespace status (cache enabled)", async () => {
        const name = "testname" 
        let response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeFalsy()
        await grpc.Ensure({ name })
        response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeTruthy()
        await grpc.Delete({ name })
        response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeFalsy()
     })
 })