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
 * @group native/iam/identity/active/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Updates identity in global database", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "test",
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.SetActive({
            namespace: "",
            uuid: id.toHexString(),
            active: true
        })

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ active: boolean }>({"_id": id})
        expect(entry?.active).toBe(true)
    })

    test("Updates identity in namespace database", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "test",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            active: true
        })

        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ active: boolean }>({"_id": id})
        expect(entry?.active).toBe(true)
    })

    test("Clears cache on update", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "test",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        let existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(1)
        
        await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            active: true
        })

        existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/identity/active/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns new identity information as a result of update", async () => {
        const name = randomBytes(32).toString('hex')
        
        const createResponse = await nativeIAmIdentityGRPC.Create({
            name,
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

        const activeResponse = await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            active: true
        })

        expect(activeResponse.identity?.uuid).toBe(id.toHexString())
        expect(activeResponse.identity?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(activeResponse.identity?.active).toBe(true)
        expect(activeResponse.identity?.name).toBe(name)
        expect(activeResponse.identity?.policies).toStrictEqual([])
    })

    test("Get request will return new data after update", async () => {
        const createResponse = await nativeIAmIdentityGRPC.Create({
            name: "test",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

        await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            active: true
        })

        const getResponse = await nativeIAmIdentityGRPC.Get({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            useCache: false
        })

        expect(getResponse.identity?.active).toBe(true)
    })

    test("Multiple updates resuts in same result", async () => {
        const createResponse = await nativeIAmIdentityGRPC.Create({
            name: "test",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

        for (let i = 0; i < 5; i += 1) {
            await nativeIAmIdentityGRPC.SetActive({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                active: true
            })
            const getResponse = await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                useCache: false
            })
            expect(getResponse.identity?.active).toBe(true)
        }

        for (let i = 0; i < 5; i += 1) {
            await nativeIAmIdentityGRPC.SetActive({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                active: false
            })
            const getResponse = await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: id.toHexString(),
                useCache: false
            })
            expect(getResponse.identity?.active).toBe(false)
        }
    })

    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmIdentityGRPC.SetActive({
                namespace: TEST_NAMESPACE_NAME,
                uuid: "invalid",
                active: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with NOT_FOUND error if identity doesnt exist", async () => {
        try {
            await nativeIAmIdentityGRPC.SetActive({
                namespace: TEST_NAMESPACE_NAME,
                uuid: new ObjectId().toHexString(),
                active: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})