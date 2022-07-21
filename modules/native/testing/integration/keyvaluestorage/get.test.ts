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
    test("Creates cache entry if cache enabled", async () => {fail()})

    test("Doesnt create cache entry if cache disabled", async () => {fail()})

    test("Gets from cache if cache enabled", async () => {fail()})

    test("Gets from DB if cache disabled", async () => {fail()})
})

/**
 * @group native/keyvaluestorage/get/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Fails with NOT_FOUND error if key doesnt exist", async () => {fail()})

    test("Returns value for specified key", async () => {fail()})
})