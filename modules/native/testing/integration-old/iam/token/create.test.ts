import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { decode as decodeJWT, JwtPayload } from 'jsonwebtoken'

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
 * @group native/iam/token/create/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Creates entry in global database if no namespace", async () => {
        const identity = randomBytes(32).toString("hex")

        const response = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity,
            metadata: "",
            scopes: []
        })
        
        const id = ObjectId.createFromHexString(response.tokenData?.uuid as string)
        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_token').findOne<{ identity: string }>({ "_id": id })
        expect(entry).not.toBeNull()
        expect(entry?.identity).toBe(identity)
    })
    test("Creates entry in namespace database", async () => {
        const identity = randomBytes(32).toString("hex")

        const response = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            scopes: []
        })
        
        const id = ObjectId.createFromHexString(response.tokenData?.uuid as string)
        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_token').findOne<{ identity: string }>({ "_id": id })
        expect(entry).not.toBeNull()
        expect(entry?.identity).toBe(identity)
    })

    interface JWTTokenData {
        uuid?: string
        namespace?: string
        identity?: string
        scopes?: Array<{
            namespace: string,
            resources: Array<string>,
            actions: Array<string>
        }>
        refresh?: boolean
    }

    test("JWT token has provided scopes", async () => {
        const identity = randomBytes(32).toString("hex")
        const metadata = randomBytes(32).toString("hex")

        const scopeNamespace = randomBytes(32).toString("hex")
        const scopeAction = randomBytes(32).toString("hex")
        const scopeResource = randomBytes(32).toString("hex")
        
        const response = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata,
            scopes: [{
                namespace: scopeNamespace,
                actions: [scopeAction],
                resources: [scopeResource]
            }]
        })

        const decoded = decodeJWT(response.token) as JwtPayload & JWTTokenData
        expect(decoded.uuid).toBe(response.tokenData?.uuid)
        expect(decoded.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(decoded.identity).toBe(identity)
        expect(decoded?.scopes).toHaveLength(1)
        expect(decoded?.scopes?.[0].namespace).toBe(scopeNamespace)
        expect(decoded?.scopes?.[0].actions).toHaveLength(1)
        expect(decoded?.scopes?.[0].resources).toHaveLength(1)
        expect(decoded?.scopes?.[0].actions[0]).toBe(scopeAction)
        expect(decoded?.scopes?.[0].resources[0]).toBe(scopeResource)
        expect(decoded.refresh).toBe(false)
    })

    test("JWT refresh token has provided scopes", async () => {
        const identity = randomBytes(32).toString("hex")

        const scopeNamespace = randomBytes(32).toString("hex")
        const scopeAction = randomBytes(32).toString("hex")
        const scopeResource = randomBytes(32).toString("hex")
        
        const response = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            scopes: [{
                namespace: scopeNamespace,
                actions: [scopeAction],
                resources: [scopeResource]
            }]
        })

        const decoded = decodeJWT(response.refreshToken) as JwtPayload & JWTTokenData
        expect(decoded.uuid).toBe(response.tokenData?.uuid)
        expect(decoded.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(decoded.identity).toBe(identity)
        expect(decoded?.scopes).toHaveLength(0) // Refresh tokens have empty scopes because theirs data is inside DB
        expect(decoded.refresh).toBe(true)
    })
})

/**
 * @group native/iam/token/create/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Creates token and returns its data", async () => {
        const identity = randomBytes(32).toString("hex")

        const scopeNamespace = randomBytes(32).toString("hex")
        const scopeAction = randomBytes(32).toString("hex")
        const scopeResource = randomBytes(32).toString("hex")
        
        const response = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            scopes: [{
                namespace: scopeNamespace,
                actions: [scopeAction],
                resources: [scopeResource]
            }]
        })

        expect(response.tokenData?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(response.tokenData?.identity).toBe(identity)
        expect(response.tokenData?.scopes).toHaveLength(1)
        expect(response.tokenData?.scopes[0].namespace).toBe(scopeNamespace)
        expect(response.tokenData?.scopes[0].actions).toHaveLength(1)
        expect(response.tokenData?.scopes[0].resources).toHaveLength(1)
        expect(response.tokenData?.scopes[0].actions[0]).toBe(scopeAction)
        expect(response.tokenData?.scopes[0].resources[0]).toBe(scopeResource)
    })

    test("Fails with FAILED_PRECONDITION if namespace doesnt exist", async () => {
        try {
            await nativeIAmTokenGRPC.Create({
                namespace: randomBytes(32).toString("hex"),
                identity: randomBytes(32).toString("hex"),
                metadata: "",
                scopes: []
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })
})