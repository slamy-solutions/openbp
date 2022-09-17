import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { managerClient as lambdaManagerGrpc, connect as connectToNativeLambda, close as closeNativeLambda } from '../../../tools/lambda/grpc'


const TEST_NAMESPACE_NAME = "lambdanamespace"
const BUNDLE_DB_NAME = `openbp_global`

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
 * @group native/lambda/manager/getBundle/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns existing lambda bundle", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const runtime = randomBytes(32).toString("hex")
        const data = randomBytes(32)

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle,
            data: data,
            ensureExactlyOneDelivery: true,
            runtime
        })
        const response = await lambdaManagerGrpc.GetBundle({ bundle })
        expect(response.data).toStrictEqual(data)
    })

    test("Returns NOT_FOUND error when lambda doesnt exist", async () => {
        try {
            await lambdaManagerGrpc.Get({ namespace: TEST_NAMESPACE_NAME, uuid: "nonexisting" })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})