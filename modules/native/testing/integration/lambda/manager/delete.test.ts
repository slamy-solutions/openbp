import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { Binary, GridFSBucket, ObjectId } from 'mongodb'
import { observable, Observable } from 'rxjs'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { managerClient as lambdaManagerGrpc, connect as connectToNativeLambda, close as closeNativeLambda } from '../../../tools/lambda/grpc'
import { Lambda } from '../../../tools/lambda/proto/lambda'


const TEST_NAMESPACE_NAME = "lambdanamespace"
const BUNDLE_DB_NAME = `openbp_global`
const INFO_DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
    await connectToNativeLambda()
})

beforeEach(async () => {
    await namespaceGrpc.Ensure({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
    try {
        await mongoClient.db(BUNDLE_DB_NAME).collection("native_lambda_manager_bundle").deleteMany({})
    } catch {}
})

afterEach(async ()=>{
    await namespaceGrpc.Delete({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
    try {
        await mongoClient.db(BUNDLE_DB_NAME).collection("native_lambda_manager_bundle").deleteMany({})
    } catch {}
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeNamespace()
    await closeNativeLambda()
})

/**
 * @group native/lambda/manager/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Removes entry from the info database", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle: randomBytes(32),
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid })
        const dbEntry = await mongoClient.db(INFO_DB_NAME).collection('native_lambda_manager_info').findOne({ uuid })
        expect(dbEntry).toBeNull()
    })

    test("Reduces bundle references count on delete", async () => {
        const uuid1 = "customlambda" + randomBytes(20).toString("hex")
        const uuid2 = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid1,
            bundle,
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid2,
            bundle,
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })

        const dbEntryBeforeDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryBeforeDelete?.references).toBe(2)

        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })

        const dbEntryAfterDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryAfterDelete?.references).toBe(1)
    })

    test("Deletes bundle from database on zero references", async () => {
        const uuids = new Array<string>(16).fill("").map(() => randomBytes(20).toString("hex"))
        const bundle = randomBytes(32)
        await Promise.all(uuids.map((uuid)=>{
            return lambdaManagerGrpc.Create({
                namespace: TEST_NAMESPACE_NAME,
                uuid: uuid,
                bundle,
                data: randomBytes(32),
                ensureExactlyOneDelivery: true,
                runtime: randomBytes(32).toString('hex')
            })
        }))

        const dbEntryBeforeDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryBeforeDelete?.references).toBe(16)


        await Promise.all(uuids.map((uuid) => {
            return lambdaManagerGrpc.Delete({ 
                namespace: TEST_NAMESPACE_NAME,
                uuid
            })
        }))

        const dbEntryAfterDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryAfterDelete).toBeNull()
    })

    test("Delete already deleted lambda doesnt reduce references count", async () => {
        const uuid1 = "customlambda" + randomBytes(20).toString("hex")
        const uuid2 = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid1,
            bundle,
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid2,
            bundle,
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })

        const dbEntryBeforeDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryBeforeDelete?.references).toBe(2)

        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: uuid1 })

        const dbEntryAfterDelete = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntryAfterDelete?.references).toBe(1)
    })
})

/**
 * @group native/lambda/manager/delete/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Lambda does not exist after deletion", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle: randomBytes(32),
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid })
        const existResponse = await lambdaManagerGrpc.Exists({ namespace: TEST_NAMESPACE_NAME, uuid })
        expect(existResponse.exists).toBeFalsy()
    })

    test("Can be deleted multiple times without error", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle: randomBytes(32),
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        for (let times = 0; times < 10; times +=1 ) {
            await lambdaManagerGrpc.Delete({ namespace: TEST_NAMESPACE_NAME, uuid })
        }
    })
})