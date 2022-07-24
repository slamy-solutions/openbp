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
 * @group native/iam/policy/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Removes value from global DB", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        let entry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).not.toBeNull()

        await nativeIAmPolicyGRPC.Delete({
            namespace: "",
            uuid: id.toHexString()
        })

        entry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).toBeNull()
    })

    test("Removes value from namespace DB", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        let entry = await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).not.toBeNull()

        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString()
        })

        entry = await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').findOne<{ name: string }>({ "_id": id })
        expect(entry).toBeNull()
    })

    test("Removes value from global cache", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        
        await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        let cacheExistResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(cacheExistResponse).toBe(1)

        await nativeIAmPolicyGRPC.Delete({
            namespace: "",
            uuid: id.toHexString()
        })

        cacheExistResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(cacheExistResponse).toBe(0)
    })

    test("Removes value from namespace cache", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        
        await nativeIAmPolicyGRPC.Get({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })
        let cacheExistResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(cacheExistResponse).toBe(1)

        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString()
        })

        cacheExistResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(cacheExistResponse).toBe(0)
    })
})

/**
 * @group native/iam/policy/delete/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Policy is not accesible after deletion", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const uuid = createResponse.policy?.uuid as string

        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })

        try {
            await nativeIAmPolicyGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid,
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
    test("Multiple deletions are OK", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: "",
            resources: [],
            actions: []
        })
        const uuid = createResponse.policy?.uuid as string
        
        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })
        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })
        await nativeIAmPolicyGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })
    })
    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmPolicyGRPC.Delete({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })
})