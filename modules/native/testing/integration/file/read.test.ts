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
 * @group native/file/read/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Fails with INVALID_ARGUMENT error if start index is out of range", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "", true)
        try {
            await file.with(async (f) => {
                await fileGrpc.Read({
                    namespace: TEST_NAMESPACE_NAME,
                    start: 15,
                    toRead: 1,
                    uuid: f.UUID
                }).forEach(() => undefined)
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Fails with INVALID_ARGUMENT error if bytes to read is out of range", async () => {
        const file = new TestFile(10, TEST_NAMESPACE_NAME, false, false, "", true)
        try {
            await file.with(async (f) => {
                await fileGrpc.Read({
                    namespace: TEST_NAMESPACE_NAME,
                    start: 5,
                    toRead: 10,
                    uuid: f.UUID
                }).forEach(() => undefined)
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Fails with INVALID_ARGUMENT on bad UUID format", async () => {
        try {
            await fileGrpc.Read({
                namespace: TEST_NAMESPACE_NAME,
                start: 5,
                toRead: 10,
                uuid: "invalid"
            }).forEach(() => undefined)
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Fails with NOT_FOUND error when file does not exist", async () => {
        try {
            await fileGrpc.Read({
                namespace: TEST_NAMESPACE_NAME,
                start: 5,
                toRead: 10,
                uuid: new ObjectId().toHexString() 
            }).forEach(() => undefined)
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })

    test("Reads from the middle of the file", async () => {
        const file = new TestFile(1024*1024*50, TEST_NAMESPACE_NAME, false, false, "", true)
        await file.with(async (f) => {
            const start = 1024*1024*15 + 3587
            const toRead = 1024*1024*15 + 7203
            const readStream = fileGrpc.Read({
                namespace: TEST_NAMESPACE_NAME,
                start,
                toRead,
                uuid: f.UUID
            })
            let readed = 0
            const readBuf = Buffer.allocUnsafe(toRead)
            await readStream.forEach((p) => {
                readed += p.chunk.length
                p.chunk.copy(readBuf, p.chunkStart - start, 0, p.chunk.length)
            })
            expect(readed).toBe(toRead)
            expect(readBuf.compare(file.data, start, start+toRead, 0, readBuf.length)).toBe(0)
        })
    })
})