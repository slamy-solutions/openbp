import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { identityClient as nativeIAmIdentityGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

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
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_identity').deleteMany({})
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
 * @group native/iam/identity/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates identity in global database if no namespace", async () => {
        const name = randomBytes(16).toString("hex")

        const response = await nativeIAmIdentityGRPC.Create({
            name,
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })

    test("Creates identity in namespace database", async () => {
        const name = randomBytes(16).toString("hex")

        const response = await nativeIAmIdentityGRPC.Create({
            name,
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })
})

/**
 * @group native/iam/identity/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Creates identity and returns its data", async () => {
        const name = randomBytes(32).toString("hex")
        const initiallyActive = true

        const response = await nativeIAmIdentityGRPC.Create({
            name,
            initiallyActive: initiallyActive,
            namespace: TEST_NAMESPACE_NAME
        })

        expect(response.identity?.name).toBe(name)
        expect(response.identity?.active).toBe(initiallyActive)
        expect(response.identity?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(response.identity?.policies).toStrictEqual([])
    })

    test("Fails with FAILED_PRECONDITION if namespace doesnt exist", async () => {
        try {
            await nativeIAmIdentityGRPC.Create({
                name: "test",
                initiallyActive: false,
                namespace: NAMESPACE_DB_NAME + "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })
})