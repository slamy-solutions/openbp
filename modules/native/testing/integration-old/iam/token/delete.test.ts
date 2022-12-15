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
 * @group native/iam/token/delete/whitebox
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

        await nativeIAmTokenGRPC.Delete({
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

        await nativeIAmTokenGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })

        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(0)
    })

    test("Removes entry from database", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string
        const id = ObjectId.createFromHexString(uuid)

        const beforeDeleteCount = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").countDocuments({"_id": id})
        expect(beforeDeleteCount).toBe(1)
        
        await nativeIAmTokenGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid
        })

        const afterDeleteCount = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").countDocuments({"_id": id})
        expect(afterDeleteCount).toBe(0)
    })
})

/**
 * @group native/iam/token/delete/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Token can not be get after deletion", async () => {
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

        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid
        })

        try {
            await nativeIAmTokenGRPC.Get({
                namespace: "",
                useCache: true,
                uuid
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })

    test("Multiple deletions are OK", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid
        })
        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid
        })
        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid
        })
    })

    test("Delete with invalid namespace is OK", async () => {
        await nativeIAmTokenGRPC.Delete({
            namespace: randomBytes(32).toString("hex"),
            uuid: new ObjectId().toHexString()
        })
    })

    test("Failes with INVALID_ARGUMENT if UUID has bad format", async () => {
        try {
            await nativeIAmTokenGRPC.Delete({
                namespace: "",
                uuid: "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })
})