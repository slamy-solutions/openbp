import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { GridFSBucket, ObjectId } from 'mongodb'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../system/testing/tools/cache'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../tools/namespace/grpc'
import { client as keyvaluestorageGrpc, connect as connectToNativeKeyValueStorage, close as closeNativeKeyValueStorage } from '../../tools/keyvaluestorage/grpc'
import { observable, Observable } from 'rxjs'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'

const TEST_NAMESPACE_NAME = "keyvaluestoragetestnamespace"
const GLOBAL_DB = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
    await connectToNativeKeyValueStorage()
})

beforeEach(async () => {
    await namespaceGrpc.Ensure({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await namespaceGrpc.Delete({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeNamespace()
    await closeNativeKeyValueStorage()
})

/**
 * @group native/keyvaluestorage/get/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates cache entry if cache enabled", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)

        const cacheKey = `native_keyvaluestorage_key_${TEST_NAMESPACE_NAME}_${key}`
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })

        let hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(0)

        await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: true
        })

        hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(1)
    })

    test("Doesnt create cache entry if cache disabled", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)

        const cacheKey = `native_keyvaluestorage_key_${TEST_NAMESPACE_NAME}_${key}`
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })

        let hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(0)

        await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: false
        })

        hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(0)
    })

    test("Gets from cache if cache enabled", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)

        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })

        await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: true
        })

        const updateResult = await mongoClient.db(DB_NAME).collection("native_keyvaluestorage").updateOne({ key }, { "$set": {value: Buffer.from("123")} })
        expect(updateResult.modifiedCount).toBe(1)

        const response = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: true
        })

        expect(response.value).toStrictEqual(value)
    })

    test("Gets from DB if cache disabled", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)

        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })

        await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: true
        })

        const updateResult = await mongoClient.db(DB_NAME).collection("native_keyvaluestorage").updateOne({ key }, { "$set": {value: Buffer.from("123")} })
        expect(updateResult.modifiedCount).toBe(1)

        const response = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key,
            useCache: false
        })

        expect(response.value).toStrictEqual(Buffer.from("123"))
    })
})

/**
 * @group native/keyvaluestorage/get/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Fails with NOT_FOUND error if key doesnt exist", async () => {
        try {
            await keyvaluestorageGrpc.Get({
                namespace: TEST_NAMESPACE_NAME,
                key: "123lala",
                useCache: false
            })
            fail()
        } catch (e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })

    test("Returns value for specified key. Cache disabled.", async () => {
        const key1 = randomBytes(32).toString("hex")
        const value1 = randomBytes(32)
        const key2 = randomBytes(32).toString("hex")
        const value2 = randomBytes(32)

        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key: key1,
            value: value1
        })
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key: key2,
            value: value2
        })

        const r1 = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key: key1,
            useCache: false
        })
        expect(r1.value).toStrictEqual(value1)

        const r2 = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key: key2,
            useCache: false
        })
        expect(r2.value).toStrictEqual(value2)
    })

    test("Returns value for specified key. Cache enabled.", async () => {
        const key1 = randomBytes(32).toString("hex")
        const value1 = randomBytes(32)
        const key2 = randomBytes(32).toString("hex")
        const value2 = randomBytes(32)

        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key: key1,
            value: value1
        })
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key: key2,
            value: value2
        })

        const r1 = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key: key1,
            useCache: true
        })
        expect(r1.value).toStrictEqual(value1)

        const r2 = await keyvaluestorageGrpc.Get({
            namespace: TEST_NAMESPACE_NAME,
            key: key2,
            useCache: true
        })
        expect(r2.value).toStrictEqual(value2)
    })
})