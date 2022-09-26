import { sign as signJWT, decode as decodeJWT, Jwt } from 'jsonwebtoken'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { tokenClient as nativeIAmTokenGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'
import { RefreshResponse_Status, ValidateResponse_Status } from '../../../tools/iam/proto/token'

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
 * @group native/iam/token/refresh/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns OK status and new token if it is valid", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: "123",
            metadata: "",
            scopes: []
        })
        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.refreshToken
        })
        expect(refreshResponse.status).toBe(RefreshResponse_Status.OK)
    })

    test("Token received after refreshing is valid", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: "123",
            metadata: "",
            scopes: []
        })
        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.refreshToken
        })
        
        const validateResponse = await nativeIAmTokenGRPC.Validate({
            token: refreshResponse.token as string,
            useCache: false
        })
        expect(validateResponse.status).toBe(ValidateResponse_Status.OK)
    })

    test("Returns INVALID status if token is not a JWT", async () => { 
        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: "invalid JWT"
        })
        expect(refreshResponse.status).toBe(RefreshResponse_Status.INVALID)
    })

    test("Returns INVALID status if token was signed by wrong secret", async () => { 
        const refreshToken = signJWT({}, "some invalid key")
        const response = await nativeIAmTokenGRPC.Refresh({
            refreshToken
        })
        expect(response.status).toBe(RefreshResponse_Status.INVALID)
    })
    
    test("Returns NOT_FOUND status if token was deleted", async () => { 
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: "123",
            metadata: "",
            scopes: []
        })

        await nativeIAmTokenGRPC.Delete({
            namespace: "",
            uuid: createResponse.tokenData?.uuid as string
        })

        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.refreshToken
        })
        expect(refreshResponse.status).toBe(RefreshResponse_Status.NOT_FOUND)
    })

    test("Returns NOT_FOUND status if namespaces of the token was deleted", async () => { 
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: "123",
            metadata: "",
            scopes: []
        })

        await nativeNamespaceGRPC.Delete({
            name: TEST_NAMESPACE_NAME
        })

        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.refreshToken
        })
        expect(refreshResponse.status).toBe(RefreshResponse_Status.NOT_FOUND)
    })

    test("Returns DISABLED status if token was manually disabled", async () => { 
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: "123",
            metadata: "",
            scopes: []
        })

        await nativeIAmTokenGRPC.DisableByUUID({
            namespace: "",
            uuid: createResponse.tokenData?.uuid as string
        })

        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.refreshToken
        })
        expect(refreshResponse.status).toBe(RefreshResponse_Status.DISABLED)
    })

    test("Returns EXPIRED status if token was expired", async () => { 
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

        const refreshToken = signJWT(payload, jwtSecretKey, { expiresIn: "1ms" })

        await new Promise(resolve => setTimeout(resolve, 1000))

        const response = await nativeIAmTokenGRPC.Refresh({
            refreshToken
        })

        expect(response.status).toBe(RefreshResponse_Status.EXPIRED)
    })

    test("Returns NOT_REFRESH_TOKEN status if this token is not a refresh token", async () => { 
        const createResponse = await nativeIAmTokenGRPC.Create({
            identity: "123",
            namespace: "",
            metadata: "",
            scopes: []
        })

        const refreshResponse = await nativeIAmTokenGRPC.Refresh({
            refreshToken: createResponse.token
        })

        expect(refreshResponse.status).toBe(RefreshResponse_Status.NOT_REFRESH_TOKEN)
    })
})