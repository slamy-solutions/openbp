import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { tokenClient as nativeIAmTokenGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `openbp_global`
const TEST_NAMESPACE_NAME = "iamtokentestnamespace"
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
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_token').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_token').deleteMany({})
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
 * @group native/iam/token/disableByUUID/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Removes global cache", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.Get({
            namespace: "",
            useCache: true,
            uuid
        })

        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(1)

        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: "",
            uuid
        })

        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(0)
    })

    test("Removes namespace cache", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.Get({
            namespace: TEST_NAMESPACE_NAME,
            useCache: true,
            uuid
        })

        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(1)

        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })

        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(0)
    })

    test("Updates entry in DB", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string
        const id = ObjectId.createFromHexString(uuid)

        const responseBeforeDisable = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").findOne<{ disabled: boolean }>({"_id": id})
        expect(responseBeforeDisable?.disabled).toBeFalsy()
        
        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })

        const responseAfterDisable = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").findOne<{ disabled: boolean }>({"_id": id})
        expect(responseAfterDisable?.disabled).toBeTruthy()
    })
})

/**
 * @group native/iam/token/disableByUUID/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("After disabling token data changed", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string
        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: "",
            uuid
        })
        const getResponse = await nativeIAmTokenGRPC.Get({
            namespace: "",
            uuid,
            useCache: false
        })
        expect(getResponse.tokenData?.disabled).toBeTruthy()
    })
    test("Disabling several times is OK", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string
        for (let i = 0; i < 5; i++) {
            await nativeIAmTokenGRPC.DisableByUUID({
                namespace: "",
                uuid
            })
            const getResponse = await nativeIAmTokenGRPC.Get({
                namespace: "",
                uuid,
                useCache: false
            })
            expect(getResponse.tokenData?.disabled).toBeTruthy()
        }
    })
    test("Failes with INVALID_ARGUMENT if UUID has bad format", async () => {
        try {
            await nativeIAmTokenGRPC.DisableByUUID({
                namespace: "",
                uuid: "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })
    test("Failes with NOT_FOUND if token with specified UUID doesnt exist", async () => {
        try {
            await nativeIAmTokenGRPC.Create({
                namespace: "",
                identity: randomBytes(32).toString("hex"),
                metadata: "",
                scopes: []
            })

            await nativeIAmTokenGRPC.DisableByUUID({
                namespace: "",
                uuid: new ObjectId().toHexString()
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
    test("Failes with NOT_FOUND if token namespace doesnt exist", async () => {
        try {
            await nativeIAmTokenGRPC.DisableByUUID({
                namespace: randomBytes(32).toString("hex"),
                uuid: new ObjectId().toHexString()
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})