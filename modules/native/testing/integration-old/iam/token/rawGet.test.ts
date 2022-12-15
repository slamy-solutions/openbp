import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { sign as signJWT } from 'jsonwebtoken'

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
 * @group native/iam/token/rawGet/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Uses global cache if cache enabled", async () => {
        const identity = randomBytes(32).toString("hex")
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity,
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true,
        })
        const id = ObjectId.createFromHexString(uuid)
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_token').updateOne({ "_id": id }, {"$set": {"identity": "invalid"}})
    
        const response = await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true,
        })
        expect(response.tokenData?.identity).toBe(identity)
    })
    test("Doesnt use global cache if cache disabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true,
        })
        const id = ObjectId.createFromHexString(uuid)
        const newIdentity = randomBytes(32).toString("hex")
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_token').updateOne({ "_id": id }, {"$set": {"identity": newIdentity}})
    
        const response = await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: false,
        })
        expect(response.tokenData?.identity).toBe(newIdentity)
    })
    test("Uses namespace cache if cache enabled", async () => {
        const identity = randomBytes(32).toString("hex")
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true
        })
        const id = ObjectId.createFromHexString(uuid)
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_token').updateOne({ "_id": id }, {"$set": {"identity": "invalid"}})
    
        const response = await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true
        })
        expect(response.tokenData?.identity).toBe(identity)
    })
    test("Doesnt use namespace cache if cache disabled", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        const uuid = createResponse.tokenData?.uuid as string

        await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: true
        })
        const id = ObjectId.createFromHexString(uuid)
        const newIdentity = randomBytes(32).toString("hex")
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_token').updateOne({ "_id": id }, {"$set": {"identity": newIdentity}})
    
        const response = await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: false
        })
        expect(response.tokenData?.identity).toBe(newIdentity)
    })
})

/**
 * @group native/iam/token/rawGet/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Gets actual data", async () => {
        const identity = randomBytes(32).toString("hex")

        const scopeNamespace = randomBytes(32).toString("hex")
        const scopeAction = randomBytes(32).toString("hex")
        const scopeResource = randomBytes(32).toString("hex")
        
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: "",
            identity,
            metadata: "",
            scopes: [{
                namespace: scopeNamespace,
                actions: [scopeAction],
                resources: [scopeResource]
            }]
        })

        const response = await nativeIAmTokenGRPC.RawGet({
            token: createResponse.token,
            useCache: false
        })

        expect(response.tokenData?.namespace).toBe("")
        expect(response.tokenData?.identity).toBe(identity)
        expect(response.tokenData?.scopes).toHaveLength(1)
        expect(response.tokenData?.scopes[0].namespace).toBe(scopeNamespace)
        expect(response.tokenData?.scopes[0].actions).toHaveLength(1)
        expect(response.tokenData?.scopes[0].resources).toHaveLength(1)
        expect(response.tokenData?.scopes[0].actions[0]).toBe(scopeAction)
        expect(response.tokenData?.scopes[0].resources[0]).toBe(scopeResource)
    })

    test("Failes with INVALID_ARGUMENT if token has invalid format", async () => {
        try {
            await nativeIAmTokenGRPC.RawGet({
                token: randomBytes(32).toString("hex"),
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with INVALID_ARGUMENT if token was signed by invalid key", async () => {
        try {
            const token = signJWT({}, "some invalid key")
            await nativeIAmTokenGRPC.RawGet({
                token,
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with NOT_FOUND error if token is valid but namespace of the token doesnt exist (was deleted)", async () => {
        const createResponse = await nativeIAmTokenGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            identity: randomBytes(32).toString("hex"),
            metadata: "",
            scopes: []
        })
        await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })
        
        try {
            await nativeIAmTokenGRPC.RawGet({
                token: createResponse.token,
                useCache: false,
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })

    test("Failes with NOT_FOUND error if token is valid but doesnt exist (was deleted)", async () => {
        try {
            const createResponse = await nativeIAmTokenGRPC.Create({
                namespace: TEST_NAMESPACE_NAME,
                identity: randomBytes(32).toString("hex"),
                metadata: "",
                scopes: []
            })
            await nativeIAmTokenGRPC.Delete({namespace: TEST_NAMESPACE_NAME, uuid: createResponse.tokenData?.uuid as string})

            await nativeIAmTokenGRPC.RawGet({
                token: createResponse.token,
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})