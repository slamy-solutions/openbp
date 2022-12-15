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
 * @group native/namescape/get/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Value is added to the cache on cache enabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name: name, useCache: true })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
    })

    test("Value is not added to the cache on cache disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name: name, useCache: false })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })

    test("Value is returned from cache when cache enabled and not from cache when disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })

        // Change value in database to contain wrong data without invalidating cache
        const newName = "newtestname"
        await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('native_namespace').updateOne({ name }, { "$set": { name: newName } })

        const cachedResponse = await grpc.Get({ name, useCache: true })
        expect(cachedResponse.namespace?.name).toBe(name)

        try {
            await grpc.Get({ name, useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }

        const uncachedResponse = await grpc.Get({ name: newName, useCache: false })
        expect(uncachedResponse.namespace?.name).toBe(newName)
    })
})

/**
 * @group native/namescape/get/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Get value same as created", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        const response = await grpc.Get({ name: name, useCache: false })
        expect(response.namespace?.name).toBe(name)
    })

    test("Get value for asked namespace", async () => {
        const namespaces = new Array<string>(16).fill("").map((_v, index) => `namespace${index}`)
        const r = await Promise.all(namespaces.map((namespace) => grpc.Ensure({ name: namespace })))
        const results = await Promise.all(namespaces.map((namespace) => grpc.Get({ name: namespace, useCache: true })))
        const resultsNames = results.map(result => result.namespace?.name)
        expect(namespaces).toEqual(resultsNames)
    })

    test("Returns error if not found", async () => {
        try {
            await grpc.Get({ name: "somevalue", useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})