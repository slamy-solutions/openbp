import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { GridFSBucket, ObjectId } from 'mongodb'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../tools/namespace/grpc'
import { client as fileGrpc, connect as connectToNativeFile, close as closeNativeFile } from '../../tools/file/grpc'
import { TestFile } from '../../tools/file/testfile'
import { observable, Observable } from 'rxjs'
import { FileCreateRequest } from '../../tools/file/proto/file'

const TEST_NAMESPACE_NAME = "filetestnamespace"
const DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
    await connectToNativeFile()
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
    await closeNativeFile()
})

/**
 * @group native/file/delete/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Deletes information in database", async () => {
        const file = new TestFile(1, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.create()
        const entry = await mongoClient.db(DB_NAME).collection<{ dataId: ObjectId }>('native_file').findOne({ _id: file.mongoId })

        await fileGrpc.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: file.UUID
        })

        const result = await mongoClient.db(DB_NAME).collection('native_file').findOne({ _id: file.mongoId })
        expect(result).toBeNull()
    })

    test("Deletes bucket with binary data", async () => {
        const file = new TestFile(1, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.create()
        const entry = await mongoClient.db(DB_NAME).collection<{ dataId: ObjectId }>('native_file').findOne({ _id: file.mongoId })

        await fileGrpc.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: file.UUID
        })

        const bucket = new GridFSBucket(mongoClient.db(DB_NAME), { bucketName: "native_file_bucket" })
        await bucket.find({ _id: entry?.dataId }).forEach(() => {
            fail()
        })
    })
})

/**
 * @group native/file/delete/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with INVALID_ARGUMENT on bad UUID format", async () => {
        try {
            await fileGrpc.Delete({
                namespace: TEST_NAMESPACE_NAME,
                uuid: "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Doesnt fail when file does not exist", async () => {
        await fileGrpc.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: new ObjectId().toHexString() 
        })
    })

    test("File cant be accessed after deletion", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "custom", true)
        await file.create()
        await fileGrpc.Stat({
            namespace: TEST_NAMESPACE_NAME,
            useCache: true,
            uuid: file.UUID 
        })
        await fileGrpc.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: file.UUID
        })
        try {
            await fileGrpc.Stat({
                namespace: TEST_NAMESPACE_NAME,
                useCache: true,
                uuid: file.UUID
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})