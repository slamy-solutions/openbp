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
 * @group native/iam/policy/get/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Gets value from cache if cache enabled", async () => {
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

        await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').updateOne({"_id": id}, {"$set": {
            "name": name2,
            "resources": resources2,
            "actions": actions2
        }})

        const response = await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        expect(response.policy?.name).toBe(name1)
        expect(response.policy?.resources).toStrictEqual(resources1)
        expect(response.policy?.actions).toStrictEqual(actions1)
    })

    test("Gets value from DB if cache disabled", async () => {
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

        await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').updateOne({"_id": id}, {"$set": {
            "name": name2,
            "resources": resources2,
            "actions": actions2
        }})

        const response = await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: false })
        expect(response.policy?.name).toBe(name2)
        expect(response.policy?.resources).toStrictEqual(resources2)
        expect(response.policy?.actions).toStrictEqual(actions2)
    })

    test("Puts value in cache if cache enabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        const existResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(existResponse).toBe(1)
    })

    test("Doesnt put value in cache if cache disabled", async () => {
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: "",
            name: "",
            resources: [],
            actions: []
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        await nativeIAmPolicyGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: false })
        const existResponse = await cacheClient.exists(`native_iam_policy_data__${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/policy/get/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Gets actual data", async () => {
        const name = randomBytes(32).toString("hex")
        const resources = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const actions = new Array<string>(10).fill("").map(() => randomBytes(16).toString("hex"))
        const createResponse = await nativeIAmPolicyGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name,
            resources,
            actions
        })
        const id = ObjectId.createFromHexString(createResponse.policy?.uuid as string)
        const response = await nativeIAmPolicyGRPC.Get({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        expect(response.policy?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(response.policy?.uuid).toBe(id.toHexString())
        expect(response.policy?.name).toBe(name)
        expect(response.policy?.resources).toStrictEqual(resources)
        expect(response.policy?.actions).toStrictEqual(actions)
    })


    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmPolicyGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with NOT_FOUND error if policy doesnt exist", async () => {
        try {
            await nativeIAmPolicyGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: new ObjectId().toHexString(),
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})