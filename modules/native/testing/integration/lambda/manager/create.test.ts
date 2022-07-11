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
const BUNDLE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const INFO_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`

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
 * @group native/lambda/manager/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in info database", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const data = randomBytes(32)
        const runtime = randomBytes(32).toString("hex")

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })

        const dbEntry = await mongoClient.db(INFO_DB_NAME).collection<Omit<Lambda, "bundle"> & { bundle: Binary }>('native_lambda_manager_info').findOne({ uuid })
        expect(dbEntry).not.toBeNull()
        expect(dbEntry?.uuid).toBe(uuid)
        expect(dbEntry?.bundle.buffer).toStrictEqual(bundle)
        expect(dbEntry?.runtime).toBe(runtime)
        expect(dbEntry?.ensureExactlyOneDelivery).toBeTruthy()
    })

    test("Creates bundle in bundle database", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const data = randomBytes(32)
        const runtime = randomBytes(32).toString("hex")

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })

        const dbEntry = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Binary, data: Binary, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntry).not.toBeNull()
        expect(dbEntry?.uuid.buffer).toStrictEqual(bundle)
        expect(dbEntry?.data.buffer).toStrictEqual(data)
        expect(dbEntry?.references).toBe(1)
    })

    test("Increases references count for same bundle", async () => {
        const uuid1 = "customlambda" + randomBytes(20).toString("hex")
        const uuid2 = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const data = randomBytes(32)
        const runtime = randomBytes(32).toString("hex")

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid1,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })
        const dbEntry1 = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntry1?.references).toBe(1)

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: uuid2,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })
        const dbEntry2 = await mongoClient.db(BUNDLE_DB_NAME).collection<{uuid: Buffer, data: Buffer, references: number}>("native_lambda_manager_bundle").findOne({ uuid: bundle })
        expect(dbEntry2?.references).toBe(2)
    })

    test("Info collection automatically creates UUID index on first insert", async () => {
        const collections = await mongoClient.db(INFO_DB_NAME).collections()
        expect(collections.findIndex((l) => l.collectionName === "native_lambda_manager_info")).toBe(-1)

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid: "customlambda" + randomBytes(20).toString("hex"),
            bundle: randomBytes(32),
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })

        const existAfter = await mongoClient.db(INFO_DB_NAME).collection("native_lambda_manager_info").indexExists("unique_uuid")
        expect(existAfter).toBeTruthy()
    })
})

/**
 * @group native/lambda/manager/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with FAILED_PRECONDITION error when namespace doesnt exist", async () => {
        try {
            await lambdaManagerGrpc.Create({
                namespace: TEST_NAMESPACE_NAME + "invalid",
                uuid: "customlambda" + randomBytes(20).toString("hex"),
                bundle: randomBytes(32),
                data: randomBytes(32),
                ensureExactlyOneDelivery: true,
                runtime: randomBytes(32).toString('hex')
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })

    test("Fails with ALREADY_EXISTS error when UUID is not unique", async () => {
        try {
            const uuid = "customlambda" + randomBytes(20).toString("hex")
            await lambdaManagerGrpc.Create({
                namespace: TEST_NAMESPACE_NAME,
                uuid,
                bundle: randomBytes(32),
                data: randomBytes(32),
                ensureExactlyOneDelivery: true,
                runtime: randomBytes(32).toString('hex')
            })
            await lambdaManagerGrpc.Create({
                namespace: TEST_NAMESPACE_NAME,
                uuid,
                bundle: randomBytes(32),
                data: randomBytes(32),
                ensureExactlyOneDelivery: true,
                runtime: randomBytes(32).toString('hex')
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.ALREADY_EXISTS)
        }
    })

    test("Creates and returns created value", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const data = randomBytes(32)
        const runtime = randomBytes(32).toString("hex")

        const response = await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })
        expect(response.lambda?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(response.lambda?.bundle).toStrictEqual(bundle)
        expect(response.lambda?.uuid).toBe(uuid)
        expect(response.lambda?.runtime).toBe(runtime)
        expect(response.lambda?.ensureExactlyOneDelivery).toBeTruthy()
    })
})