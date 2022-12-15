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
 * @group native/iam/policy/update/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Updates value in global DB", async () => {
        const name1 = randomBytes(32).toString("hex")
        const resources1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))

        const name2 = randomBytes(32).toString("hex")
        const resources2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: name1,
            resources: resources1,
            actions: actions1
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Update({
            namespace: "",
            uuid: id.toHexString(),
            name: name2,
            resources: resources2,
            actions: actions2
        })

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').findOne<{ name: string, resources: Array<string>, actions: Array<string> }>({ "_id": id })
        expect(entry?.name).toBe(name2)
        expect(entry?.resources).toStrictEqual(resources2)
        expect(entry?.actions).toStrictEqual(actions2)
    })

    test("Updates value in namespace DB", async () => {
        const name1 = randomBytes(32).toString("hex")
        const resources1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))

        const name2 = randomBytes(32).toString("hex")
        const resources2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: name1,
            resources: resources1,
            actions: actions1
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Update({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            name: name2,
            resources: resources2,
            actions: actions2
        })

        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').findOne<{ name: string, resources: Array<string>, actions: Array<string> }>({ "_id": id })
        expect(entry?.name).toBe(name2)
        expect(entry?.resources).toStrictEqual(resources2)
        expect(entry?.actions).toStrictEqual(actions2)
    })

    test("Clears cache on update", async () => {
        const name = randomBytes(32).toString("hex")
        const resources = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))

        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: name,
            resources: resources,
            actions: actions
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
    
        await nativeIAmPolicyGRPC.Get({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: true })
        let existResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(1)

        await nativeIAmPolicyGRPC.Update({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            name: "",
            resources: [],
            actions: []
        })

        existResponse = await cacheClient.exists(`native_iam_policy_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/policy/update/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Updates data so next requests will get updated data", async () => {
        const name1 = randomBytes(32).toString("hex")
        const resources1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions1 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))

        const name2 = randomBytes(32).toString("hex")
        const resources2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions2 = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name: name1,
            resources: resources1,
            actions: actions1
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)

        await nativeIAmPolicyGRPC.Update({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            name: name2,
            resources: resources2,
            actions: actions2
        })

        const response = await nativeIAmPolicyGRPC.Get({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            useCache: false
        })
        expect(response.policy?.name).toBe(name2)
        expect(response.policy?.resources).toStrictEqual(resources2)
        expect(response.policy?.actions).toStrictEqual(actions2)
    })

    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmPolicyGRPC.Update({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
                name: "",
                actions: [],
                resources: []
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with NOT_FOUND error if policy doesnt exist", async () => {
        try {
            await nativeIAmPolicyGRPC.Update({
                namespace: TEST_NAMESPACE_NAME,
                uuid: new ObjectId().toHexString(),
                name: "",
                resources: [],
                actions: []
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})