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
    await closeMongo()
    await closeCache()
    await closeNativeNamespace()
})

/**
 * @group native/namescape/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Value is deleted from the database", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        let entry = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('native_namespace').findOne({ name })
        expect(entry).not.toBeNull()
        await grpc.Delete({ name })
        entry = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('native_namespace').findOne({ name })
        expect(entry).toBeNull()
    })

    test("Cache for namespace data is deleted from redis", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })
        let response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
        await grpc.Delete({ name })
        response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })

    test("Cache for namespaces list is deleted from redis", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        let response = await cacheClient.get("native_namespace_list")
        expect(response).not.toBeNull()
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        await grpc.Delete({ name })
        response = await cacheClient.get("native_namespace_list")
        expect(response).toBeNull()
    })

    test("Database for namespace is deleted from MongoDB server", async () => {
        const name = "testname"
        await grpc.Ensure({ name })

        const collectionName = "testingcollection"
        await mongoClient.db(`openbp_namespace_${name}`).collection<{ name: string }>(collectionName).insertOne({ name: "testinsert" })

        await grpc.Delete({ name })

        await new Promise<void>((resolve) => setTimeout(resolve, 4000))

        await mongoClient.db(`openbp_namespace_${name}`).listCollections({ name: collectionName }).forEach((c) => {
            fail()
        })
    })
})

/**
 * @group native/namescape/delete/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Can\'t be accessed after deletion", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: false })
        await grpc.Delete({ name })
        
        try {
            await grpc.Get({ name, useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== 5) fail()
        }
    })
 })