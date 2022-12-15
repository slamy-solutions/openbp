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
const GLOBAL_DB = 'openbp_global'
const DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`

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
 * @group native/keyvaluestorage/set/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in namespace collection", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })
        const response = await mongoClient.db(DB_NAME).collection("native_keyvaluestorage").findOne<{ key: string, value: Buffer }>({ key })
        expect(response?.key).toBe(key)
        expect(response?.value.buffer).toStrictEqual(value)
    })

    test("Creates entry in global collection", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)
        await keyvaluestorageGrpc.Set({
            namespace: "",
            key,
            value
        })
        const response = await mongoClient.db(GLOBAL_DB).collection("native_keyvaluestorage").findOne<{ key: string, value: Buffer }>({ key })
        expect(response?.key).toBe(key)
        expect(response?.value.buffer).toStrictEqual(value)
    })

    test("Deletes entry in cache", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)

        const cacheKey = `native_keyvaluestorage_key_${TEST_NAMESPACE_NAME}_${key}`
        await cacheClient.set(cacheKey, value)
        let hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(1)

        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })

        hasKey = await cacheClient.exists(cacheKey)
        expect(hasKey).toBe(0)
    })
})

/**
 * @group native/keyvaluestorage/set/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Sets data", async () => {
        const key = randomBytes(32).toString("hex")
        const value = randomBytes(32)
        await keyvaluestorageGrpc.Set({
            namespace: TEST_NAMESPACE_NAME,
            key,
            value
        })
        const response = await keyvaluestorageGrpc.Get({
            key: key,
            namespace: TEST_NAMESPACE_NAME,
            useCache: false
        })

        expect(response.value).toStrictEqual(value)
    })

    test("Failes with FAILED_PRECONDITION if namespace doesnt exist", async () => {
        try {
            await keyvaluestorageGrpc.Set({
                namespace: TEST_NAMESPACE_NAME + "notexist",
                key: randomBytes(32).toString("hex"),
                value: randomBytes(32)
            })
            fail()
        } catch (e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })

    test("Failes with INVALID_ARGUMENT if key or value is too big", async () => {
        try {
            await keyvaluestorageGrpc.Set({
                namespace: TEST_NAMESPACE_NAME,
                key: randomBytes(32).toString("hex"),
                value: randomBytes(1024*1024*15+10)
            })
            fail()
        } catch (e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })
})