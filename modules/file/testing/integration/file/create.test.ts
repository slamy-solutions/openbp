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
const DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`

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
 * @group native/file/create/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Creates entry in database", async () => {
        const size = 64328
        const mimeType = "text"
        const file = new TestFile(size, TEST_NAMESPACE_NAME, false, false, mimeType, true)
        await file.with(async (f) => {
            const entry = await mongoClient.db(DB_NAME).collection<{ mimeType: string, size: number }>('native_file').findOne({ _id: f.mongoId })
            expect(entry).not.toBeNull()
            expect(entry?.mimeType).toBe(mimeType)
            expect(entry?.size).toBe(size)
        })
    })

    test("Creates entry with same data in gridfs", async () => {
        const size = 64328
        const mimeType = "text"
        const file = new TestFile(size, TEST_NAMESPACE_NAME, false, false, mimeType, true)
        await file.with(async (f) => {
            const entry = await mongoClient.db(DB_NAME).collection<{ dataId: ObjectId }>('native_file').findOne({ _id: f.mongoId })
            const bucket = new GridFSBucket(mongoClient.db(DB_NAME), { bucketName: "native_file_bucket" })
            const download = bucket.openDownloadStream(entry?.dataId as ObjectId)

            const reader = new Observable<Buffer>((observable) => {
                download.addListener('data', (data: Buffer) => observable.next(data))
                download.addListener('error', (e: Error) => observable.error(e))
                download.addListener('end', () => observable.complete())
                download.resume()
            })

            let currentIndex = 0
            await reader.forEach((data: Buffer) => {
                expect(f.data.slice(currentIndex, currentIndex + data.length).compare(data)).toBe(0)
                currentIndex += data.length
            })
            download.destroy()
            expect(currentIndex).toBe(f.data.length)
        })
    })
})

/**
 * @group native/file/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with ABORTED error on sensing zero packages", async () => {
        const sender = new Observable<FileCreateRequest>((obs) => {
            obs.complete()
        })

        try {
            await fileGrpc.Create(sender)
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.ABORTED)
        }
    })

    test("Fails with DATA_LOSS error when first package is not info package", async () => {
        const sender = new Observable<FileCreateRequest>((obs) => {
            obs.next({
                chunk: {
                    data: Buffer.from("")
                },
                info: undefined
            })

            obs.complete()
        })

        try {
            await fileGrpc.Create(sender)
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.DATA_LOSS)
        }
    })

    test("Fails with FAILED_PRECONDITION error when namespace doesnt exist", async () => {
        const file = new TestFile(1, "unexistingnamespace", false, false, "", true)
        try {
            await file.create()
            fail()
        } catch (e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })
    
    test("Creates entry and returns data on creation", async () => {
        const size = 64328
        const mimeType = "text"
        const file = new TestFile(size, TEST_NAMESPACE_NAME, false, false, mimeType, true)
        const creationResponse = await file.create()
        expect(creationResponse.file?.size).toBe(size)
        expect(creationResponse.file?.mimeType).toBe(mimeType)
        expect(creationResponse.file?.namespace).toBe(TEST_NAMESPACE_NAME)
    })

    test("Accessible after creation", async () => {
        const file = new TestFile(1, TEST_NAMESPACE_NAME, false, false, "", true)
        const creationResponse = await file.create()
        await fileGrpc.Stat({ namespace: TEST_NAMESPACE_NAME, useCache: false, uuid: creationResponse.file?.uuid as string })
    })

    test("Has SHA512 calculated on creation", async () => {
        const file = new TestFile(1, TEST_NAMESPACE_NAME, false, false, "", true)
        const creationResponse = await file.create()
        expect(creationResponse.file?.SHA512HashCalculated).toBeTruthy()
        expect(creationResponse.file?.SHA512Hash).not.toBe("")
    })

    test("Big file (100 MB) has valid data after creation", async () => {
        const size = 100 * 1024 * 1024
        const file = new TestFile(size, TEST_NAMESPACE_NAME, false, false, "", true)
        const creationResponse = await file.create()
        expect(creationResponse.file?.size).toBe(size)

        const stream = fileGrpc.Read({ namespace: TEST_NAMESPACE_NAME, start: 0, toRead: 0, uuid: file.UUID })
        let index = 0;
        await stream.forEach((data) => {
            expect(data.chunkStart).toBe(index)
            expect(data.transfered).toBe(index + data.chunk.length)
            expect(data.chunk.compare(file.data.slice(index, index + data.chunk.length))).toBe(0)
            index += data.chunk.length
        })
    })
})