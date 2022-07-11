import { randomBytes } from 'crypto'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { managerClient as lambdaManagerGrpc, connect as connectToNativeLambda, close as closeNativeLambda } from '../../../tools/lambda/grpc'


const TEST_NAMESPACE_NAME = "lambdanamespace"
const BUNDLE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`

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
 * @group native/lambda/manager/exists/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns true when lambda exist", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle: randomBytes(32),
            data: randomBytes(32),
            ensureExactlyOneDelivery: true,
            runtime: randomBytes(32).toString('hex')
        })
        const response = await lambdaManagerGrpc.Exists({ namespace: TEST_NAMESPACE_NAME, uuid })
        expect(response.exists).toBeTruthy()
    })

    test("Returns false when lambda doesnt exist", async () => {
        const response = await lambdaManagerGrpc.Exists({ namespace: TEST_NAMESPACE_NAME, uuid: "nonexisting" })
        expect(response.exists).toBeFalsy()
    })

    test("Returns false after lambda deletion", async () => {
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
        const response = await lambdaManagerGrpc.Exists({ namespace: TEST_NAMESPACE_NAME, uuid })
        expect(response.exists).toBeFalsy()
    })
})