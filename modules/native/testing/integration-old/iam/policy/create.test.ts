import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { policyClient as nativeIAmPolicyGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `openbp_global`
const TEST_NAMESPACE_NAME = "liampolicytestnamespace"
const NAMESPACE_DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`


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
 * @group native/iam/policy/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in global database if no namespace", async () => {
        const name = randomBytes(32).toString("hex")
        
        const response = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name,
            actions: [],
            resources: []
        })
        const id = ObjectId.createFromHexString(response.policy?.uuid as string)
        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })
    test("Creates entry in database for specific namespace", async () => {
        const name = randomBytes(32).toString("hex")

        const response = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name,
            actions: [],
            resources: []
        })
        const id = ObjectId.createFromHexString(response.policy?.uuid as string)
        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })
})

/**
 * @group native/iam/policy/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Creates policy and returns its data", async () => {
        const name = randomBytes(32).toString("hex")
        const resources = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const response = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name,
            actions,
            resources
        })
        expect(response.policy?.name).toBe(name)
        expect(response.policy?.resources).toStrictEqual(resources)
        expect(response.policy?.actions).toStrictEqual(actions)
        expect(response.policy?.namespace).toBe(TEST_NAMESPACE_NAME)
    })
    test("Fails with FAILED_PRECONDITION if namespace doesnt exist", async () => {
        try {
            await nativeIAmPolicyGRPC.Create({
                namespace: TEST_NAMESPACE_NAME + "invalid",
                name: "",
                resources: [],
                actions: []
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })
})