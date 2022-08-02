import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { policyClient as nativeIAmPolicyGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const TEST_NAMESPACE_NAME = "liampolicytestnamespace"
const NAMESPACE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`


beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeIAM()
    await connectToNativeNamespace()
    
})

beforeEach(async () => {
    await nativeNamespaceGRPC.Ensure({ name: TEST_NAMESPACE_NAME })
})

afterEach(async ()=>{
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').deleteMany({})
    } catch {}
    await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeIAM()
    await closeNativeNamespace()
})


/**
 * @group native/iam/policy/exist/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Gets value from global cache if cache enabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').deleteOne({"_id": id})

        const response = await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: true })
        expect(response.exist).toBeTruthy()
    })
    
    test("Gets value from namespace cache if cache enabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })

        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').deleteOne({"_id": id})

        const response = await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })
        expect(response.exist).toBeTruthy()
    })

    test("Gets value from global DB if cache disabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').deleteOne({"_id": id})

        const response = await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: false })
        expect(response.exist).toBeFalsy()
    })

    test("Gets value from namespace DB if cache disabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })

        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').deleteOne({"_id": id})

        const response = await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        expect(response.exist).toBeFalsy()
    })

    test("Puts value in global cache if cache enabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: true })
        const existResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(existResponse).toBe(1)
    })

    test("Puts value in namespace cache if cache enabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })
        const existResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(1)
    })

    test("Doesnt put value in global cache if cache disabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Exist({ namespace: "", uuid: id.toHexString(), useCache: false })
        const existResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })

    test("Doesnt put value in namespace cache if cache disabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        const existResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/policy/exist/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Gets actual data", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        let response = await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        expect(response.exist).toBeTruthy()

        await nativeIAmPolicyGRPC.Delete({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString() })

        response = await nativeIAmPolicyGRPC.Exist({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        expect(response.exist).toBeFalsy()
    })


    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmPolicyGRPC.Exist({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })
})