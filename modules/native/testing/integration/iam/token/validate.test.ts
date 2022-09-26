import { ObjectId } from 'mongodb'
import { sign as signJWT, decode as decodeJWT } from 'jsonwebtoken'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { tokenClient as nativeIAmTokenGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'
import { ValidateResponse_Status } from '../../../tools/iam/proto/token'

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
 * @group native/iam/token/validate/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Uses global cache if cache enabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(0)
        const r = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(1)
        const updateResponse = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_token").updateOne({"_id": ObjectId.createFromHexString(uuid)}, {"$set": {"disabled": true}})
        expect(updateResponse.modifiedCount).toBe(1)

        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.OK)
    })

    test("Doesnt use global cache if cache disabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(0)
        await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(await cacheClient.exists(`native_iam_token_data__${uuid}`)).toBe(1)
        const updateResponse = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_token").updateOne({"_id": ObjectId.createFromHexString(uuid)}, {"$set": {"disabled": true}})
        expect(updateResponse.modifiedCount).toBe(1)

        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: false
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.DISABLED)
    })

    test("Uses namespace cache if cache enabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: TEST_NAMESPACE_NAME,
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(0)
        await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(1)
        const updateResponse = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").updateOne({"_id": ObjectId.createFromHexString(uuid)}, {"$set": {"disabled": true}})
        expect(updateResponse.modifiedCount).toBe(1)

        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.OK)
    })

    test("Doesnt use namespace cache if cache disabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: TEST_NAMESPACE_NAME,
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(0)
        await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(await cacheClient.exists(`native_iam_token_data_${TEST_NAMESPACE_NAME}_${uuid}`)).toBe(1)
        const updateResponse = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_token").updateOne({"_id": ObjectId.createFromHexString(uuid)}, {"$set": {"disabled": true}})
        expect(updateResponse.modifiedCount).toBe(1)

        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: false
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.DISABLED)
    })
})

/**
 * @group native/iam/token/validate/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns OK status and actual token data if it is valid", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: TEST_NAMESPACE_NAME,
            metadata: "",
            scopes: []
        })
        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: false
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.OK)
    })
    
    test("Returns INVALID status if token is not JWT token", async () => {
        const repsponse = await nativeIAmTokenGRPC.Validate({
            token: "not a JWT",
            useCache: false
        })
        expect(repsponse.status).toBe(ValidateResponse_Status.INVALID)
    })

    test("Returns INVALID status if token is signed with wrong key", async () => {
        const token = signJWT({}, "some invalid key")
        const response = await nativeIAmTokenGRPC.Validate({
            token,
            useCache: false
        })
        expect(response.status).toBe(ValidateResponse_Status.INVALID)
    })

    test("Returns NOT_FOUND status if token was deleted", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })
        
        const repsponseBeforeDelete = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponseBeforeDelete.status).toBe(ValidateResponse_Status.OK)

        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid: createResponse.tokenData?.uuid as string
        })

        const repsponseAfterDelete = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponseAfterDelete.status).toBe(ValidateResponse_Status.NOT_FOUND)
    })

    test("Returns NOT_FOUND status if token namespace was deleted", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: TEST_NAMESPACE_NAME,
            metadata: "",
            scopes: []
        })
        
        const repsponseBeforeDelete = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: false
        })
        expect(repsponseBeforeDelete.status).toBe(ValidateResponse_Status.OK)

        await nativeNamespaceGRPC.Delete({
            name: TEST_NAMESPACE_NAME
        })

        const repsponseAfterDelete = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: false
        })
        expect(repsponseAfterDelete.status).toBe(ValidateResponse_Status.NOT_FOUND)
    })

    test("Returns DISABLED status if token was manually disabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })
        
        const repsponseBeforeDisable = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponseBeforeDisable.status).toBe(ValidateResponse_Status.OK)

        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: "",
            uuid: createResponse.tokenData?.uuid as string
        })

        const repsponseAfterDisable = await nativeIAmTokenGRPC.Validate({
            token: createResponse.token,
            useCache: true
        })
        expect(repsponseAfterDisable.status).toBe(ValidateResponse_Status.DISABLED)
    })

    test("Returns EXPIRED status if token is too old", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })

        // TODO: must be received dynamically
        const jwtSecretKey = "my_secret_key"
        const payload = decodeJWT(createResponse.token as string) as object as {exp?: string}
        delete payload["exp"]

        const token = signJWT(payload, jwtSecretKey, { expiresIn: "1ms" })

        await new Promise(resolve => setTimeout(resolve, 1000))

        const repsponse = await nativeIAmTokenGRPC.Validate({
            token,
            useCache: false
        })

        expect(repsponse.status).toBe(ValidateResponse_Status.EXPIRED)
    })
})