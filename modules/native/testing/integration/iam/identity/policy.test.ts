import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { identityClient as nativeIAmIdentityGRPC, policyClient as nativeIAmPolicyGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const TEST_NAMESPACE_NAME = "iamidentitytestnamespace"
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
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_identity').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_identity').deleteMany({})
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
 * @group native/iam/identity/policy/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    describe("Add", () => {
        test("Clears global cache on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const id = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: id.toHexString(),
                useCache: true
            })

            let existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
            expect(existResponse).toBe(1)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
            expect(existResponse).toBe(0)
        })

        test("Clears namespace cache on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: TEST_NAMESPACE_NAME
            })
            const id = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                useCache: true
            })

            let existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
            expect(existResponse).toBe(1)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
            expect(existResponse).toBe(0)
        })

        test("Updates data in global DB on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            let entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(1)
            expect(entry?.policies[0]).toBe(`:${policyId.toHexString()}`)
        })

        test("Updates data in namespace DB on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: TEST_NAMESPACE_NAME
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            let entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(1)
            expect(entry?.policies[0]).toBe(`:${policyId.toHexString()}`)
        })
    })

    describe("Remove", () => {
        test("Clears global cache on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const id = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: id.toHexString(),
                useCache: true
            })

            let existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
            expect(existResponse).toBe(1)

            await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: "",
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
            expect(existResponse).toBe(0)
        })

        test("Updates data in global DB on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            let entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(1)

            await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: "",
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(0)
        })

        test("Clears namespace cache on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: TEST_NAMESPACE_NAME
            })
            const id = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                useCache: true
            })

            let existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
            expect(existResponse).toBe(1)

            await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
            expect(existResponse).toBe(0)
        })

        test("Updates data in namespace DB on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: TEST_NAMESPACE_NAME
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            let entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(1)

            await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: TEST_NAMESPACE_NAME,
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ policies: Array<string> }>({"_id": identityId})
            expect(entry).not.toBeNull()
            expect(entry?.policies).toHaveLength(0)
        })
    })
})

/**
 * @group native/iam/identity/policy/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    describe("Add", () => {
        test("Fails with INVALID_ARGUMENT error if identity UUID has bad format", async () => {
            try {
                const policy = await nativeIAmPolicyGRPC.Create({
                    name: "",
                    namespace: "",
                    resources: [],
                    actions: []
                })
                const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

                await nativeIAmIdentityGRPC.AddPolicy({
                    identityNamespace: "",
                    identityUUID: "invalid",
                    policyNamespace: "",
                    policyUUID: policyId.toHexString()
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
            }
        })
        test("Fails with INVALID_ARGUMENT error if policy UUID has bad format", async () => {
            try {
                const response = await nativeIAmIdentityGRPC.Create({
                    name: "123123",
                    initiallyActive: false,
                    namespace: TEST_NAMESPACE_NAME
                })
                const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

                await nativeIAmIdentityGRPC.AddPolicy({
                    identityNamespace: "",
                    identityUUID: identityId.toHexString(),
                    policyNamespace: "",
                    policyUUID: "invalid"
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
            }
        })
        test("Fails with FAILED_PRECONDITION if policy doesnt exist", async () => {
            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            try {
                await nativeIAmIdentityGRPC.AddPolicy({
                    identityNamespace: "",
                    identityUUID: identityId.toHexString(),
                    policyNamespace: "",
                    policyUUID: new ObjectId().toHexString()
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
            }
        })
        test("Fails with NOT_FOUND error if identity doesnt exist", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            try {
                await nativeIAmIdentityGRPC.AddPolicy({
                    identityNamespace: "",
                    identityUUID: new ObjectId().toHexString(),
                    policyNamespace: "",
                    policyUUID: policyId.toHexString()
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
            }
        })
        test("After adding, 'get' command returns actual data", async () => {
            const policy1 = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policy1Id = ObjectId.createFromHexString(policy1.policy?.uuid as string)

            const policy2 = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policy2Id = ObjectId.createFromHexString(policy2.policy?.uuid as string)

            const createResponse = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policy1Id.toHexString()
            })

            const get1Response = await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: id.toHexString(),
                useCache: true
            })

            expect(get1Response.identity?.policies).toHaveLength(1)
            expect(get1Response.identity?.policies[0].uuid).toBe(policy1Id.toHexString())

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: id.toHexString(),
                policyNamespace: "",
                policyUUID: policy2Id.toHexString()
            })

            const get2Response = await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: id.toHexString(),
                useCache: false
            })
            expect(get2Response.identity?.policies).toHaveLength(2)
            expect(get2Response.identity?.policies.findIndex((i) => i.uuid === policy1Id.toHexString())).toBeGreaterThanOrEqual(0)
            expect(get2Response.identity?.policies.findIndex((i) => i.uuid === policy2Id.toHexString())).toBeGreaterThanOrEqual(0)
        })
        test("New state of the identity is returned on update", async () => {
            const name = randomBytes(32).toString("hex")
            
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name,
                initiallyActive: false,
                namespace: ""
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)

            await nativeIAmIdentityGRPC.SetActive({
                namespace: "",
                uuid: identityId.toHexString(),
                active: true
            })

            const updateResponse = await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID: identityId.toHexString(),
                policyNamespace: "",
                policyUUID: policyId.toHexString()
            })

            expect(updateResponse.identity?.active).toBeTruthy()
            expect(updateResponse.identity?.name).toBe(name)
            expect(updateResponse.identity?.policies).toHaveLength(1)
            expect(updateResponse.identity?.policies[0].uuid).toBe(policyId.toHexString())
        })
        test("Adding same policy several times result in adding only one item", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyId = ObjectId.createFromHexString(policy.policy?.uuid as string)

            const response = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityId = ObjectId.createFromHexString(response.identity?.uuid as string)


            for(let i = 1; i < 5; i+=1) {
                await nativeIAmIdentityGRPC.AddPolicy({
                    identityNamespace: "",
                    identityUUID: identityId.toHexString(),
                    policyNamespace: "",
                    policyUUID: policyId.toHexString()
                })
                const getResponse = await nativeIAmIdentityGRPC.Get({
                    namespace: "",
                    uuid: identityId.toHexString(),
                    useCache: false
                })
                expect(getResponse.identity?.policies).toHaveLength(1)
            }
        })
    })

    describe("Remove", () => {
        test("Fails with INVALID_ARGUMENT error if identity UUID has bad format", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyUUID = policy.policy?.uuid as string

            const identity = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityUUID = identity.identity?.uuid as string

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID
            })

            try {
                await nativeIAmIdentityGRPC.RemovePolicy({
                    identityNamespace: "",
                    identityUUID: "invalid",
                    policyNamespace: "",
                    policyUUID
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
            }
        })
        test("Fails with NOT_FOUND error if identity doesnt exist", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyUUID = policy.policy?.uuid as string
            
            try {
                await nativeIAmIdentityGRPC.RemovePolicy({
                    identityNamespace: "",
                    identityUUID: new ObjectId().toHexString(),
                    policyNamespace: "",
                    policyUUID
                })
                fail()
            } catch(e) {
                expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
            }
        })
        test("After removing, 'get' command returns actual data", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyUUID = policy.policy?.uuid as string

            const identity = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityUUID = identity.identity?.uuid as string

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID
            })

            await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID
            })
            
            const response = await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: identityUUID,
                useCache: false
            })
            expect(response.identity?.policies).toHaveLength(0)
        })
        test("New state of the identity is returned on update", async () => {
            const policy = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policyUUID = policy.policy?.uuid as string

            const identity = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityUUID = identity.identity?.uuid as string

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID
            })

            await nativeIAmIdentityGRPC.SetActive({
                namespace: "",
                uuid: identityUUID,
                active: true
            })

            const response = await nativeIAmIdentityGRPC.RemovePolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID
            })
            expect(response.identity?.active).toBeTruthy()
            expect(response.identity?.policies).toHaveLength(0)
        })
        test("Removing same policy several times result in removing only one item", async () => {
            const policy1 = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policy1UUID = policy1.policy?.uuid as string

            const policy2 = await nativeIAmPolicyGRPC.Create({
                name: "",
                namespace: "",
                resources: [],
                actions: []
            })
            const policy2UUID = policy2.policy?.uuid as string

            const identity = await nativeIAmIdentityGRPC.Create({
                name: "123123",
                initiallyActive: false,
                namespace: ""
            })
            const identityUUID = identity.identity?.uuid as string

            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID: policy1UUID
            })
            await nativeIAmIdentityGRPC.AddPolicy({
                identityNamespace: "",
                identityUUID,
                policyNamespace: "",
                policyUUID: policy2UUID
            })

            for(let i = 1; i < 5; i++) {
                await nativeIAmIdentityGRPC.RemovePolicy({
                    identityNamespace: "",
                    identityUUID,
                    policyNamespace: "",
                    policyUUID: policy1UUID
                })
                const response = await nativeIAmIdentityGRPC.Get({
                    namespace: "",
                    uuid: identityUUID,
                    useCache: false
                })
                expect(response.identity?.policies).toHaveLength(1)
                expect(response.identity?.policies[0].uuid).toBe(policy2UUID)
            }
        })
    })
})