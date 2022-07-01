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
 * @group native/file/exists/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Fails with INVALID_ARGUMENT error if UUID has invalid format", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "", true)
        try {
            await file.with(async (f) => {
                await fileGrpc.Exists({
                    namespace: TEST_NAMESPACE_NAME,
                    useCache: true,
                    uuid: "invalid"
                })
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Returns true after creation", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.with(async (f) => {
            const response = await fileGrpc.Exists({
                namespace: TEST_NAMESPACE_NAME,
                useCache: true,
                uuid: f.UUID
            })
            expect(response.exist).toBeTruthy()
        })
    })

    test("Returns false if wasnt created", async () => {
        const response = await fileGrpc.Exists({
            namespace: TEST_NAMESPACE_NAME,
            useCache: true,
            uuid: new ObjectId().toHexString()
        })
        expect(response.exist).toBeFalsy()
    })

    test("Returns true after deletion", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.with(async () => undefined )
        const response = await fileGrpc.Exists({
            namespace: TEST_NAMESPACE_NAME,
            useCache: true,
            uuid: file.UUID
        })
        expect(response.exist).toBeFalsy()
    })
})