import { randomBytes, createHash } from 'crypto'
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
 * @group native/file/sha512/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Doesnt recalculate if already calculated", async () => {
        const file = new TestFile(1, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.create()
        await mongoClient.db(DB_NAME).collection<{ dataId: ObjectId }>('native_file').updateOne({ _id: file.mongoId }, { $set: { sha512hash: Buffer.from("123") }})

        const response = await fileGrpc.CalculateSHA512({
            namespace: TEST_NAMESPACE_NAME,
            uuid: file.UUID
        })
        expect(response.SHA512.toString()).toBe("123")

        const findResponse = await mongoClient.db(DB_NAME).collection<{ sha512hash: Buffer }>('native_file').findOne({ _id: file.mongoId })
        expect(findResponse?.sha512hash.toString()).toBe("123")
    })
})

/**
 * @group native/file/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with INVALID_ARGUMENT on bad UUID format", async () => {
        try {
            await fileGrpc.CalculateSHA512({
                namespace: TEST_NAMESPACE_NAME,
                uuid: "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Fails with NOT_FOUND error when file does not exist", async () => {
        try {
            await fileGrpc.CalculateSHA512({
                namespace: TEST_NAMESPACE_NAME,
                uuid: new ObjectId().toHexString() 
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
    test("Calculates hash on creation", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "custom", true)
        const response = await file.create()
        expect(response.file?.SHA512HashCalculated).toBeTruthy()
    })

    test("Returs actual hash after execution", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "custom", true)
        const response = await file.create()
        const hasher = createHash("sha512")
        hasher.update(file.data)
        const hash = hasher.digest()
        expect(hash).toStrictEqual(response.file?.SHA512Hash)
    })
})