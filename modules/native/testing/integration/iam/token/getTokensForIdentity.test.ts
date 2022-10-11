import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { sign as signJWT, decode as decodeJWT, Jwt } from 'jsonwebtoken'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { tokenClient as nativeIAmTokenGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'
import { GetTokensForIdentityRequest_ActiveFilter, RefreshResponse_Status, TokenData, ValidateResponse, ValidateResponse_Status } from '../../../tools/iam/proto/token'

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
 * @group native/iam/token/getTokensForIdentity/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Get tokens only for specified identity", async () => { 
        const identitiesCount = 5
        const tokensCount = 30
        
        const identities = new Array(identitiesCount).fill("").map(() => randomBytes(32).toString("hex"))
        const tokens = [] as string[]
        for(let i = 0; i < tokensCount; i++) {
            const response = await nativeIAmTokenGRPC.Create({
                identity: identities[i % identitiesCount],
                metadata: "",
                namespace: TEST_NAMESPACE_NAME,
                scopes: []
            })
            tokens.push(response.tokenData?.uuid as string)
        }
        const identityTokens = tokens.filter((_val, i) => {
            return i % identitiesCount == 0
        })

        const receivedTokens = [] as string[]
        await nativeIAmTokenGRPC.GetTokensForIdentity({
            activeFilter: GetTokensForIdentityRequest_ActiveFilter.ALL,
            identity: identities[0],
            limit: 0,
            skip: 0,
            namespace: TEST_NAMESPACE_NAME
        }).forEach((r) => {
            receivedTokens.push(r.tokenData?.uuid as string)
        })

        expect(identityTokens.sort()).toStrictEqual(receivedTokens.sort())
    })

    test("Returns newly created tokens first", async () => {
        const identity = "123"
        const tokensCount = 5
        const tokens = [] as string[]
        for(let i = 0; i < tokensCount; i++) {
            const response = await nativeIAmTokenGRPC.Create({
                identity,
                metadata: "",
                namespace: TEST_NAMESPACE_NAME,
                scopes: []
            })
            tokens.push(response.tokenData?.uuid as string)
        }

        const receivedTokens = [] as string[]
        await nativeIAmTokenGRPC.GetTokensForIdentity({
            activeFilter: GetTokensForIdentityRequest_ActiveFilter.ALL,
            identity,
            limit: 0,
            skip: 0,
            namespace: TEST_NAMESPACE_NAME
        }).forEach((r) => {
            receivedTokens.push(r.tokenData?.uuid as string)
        })

        expect(tokens).toStrictEqual(receivedTokens.reverse())
    })

    describe("Skips and limits stream of returned tokens", () => {
        test("Total: 3. Skip: 1. Limit: 1.", async () => {
            const identity = "123"
            const tokensCount = 3
            const tokens = [] as string[]
            for(let i = 0; i < tokensCount; i++) {
                const response = await nativeIAmTokenGRPC.Create({
                    identity,
                    metadata: "",
                    namespace: TEST_NAMESPACE_NAME,
                    scopes: []
                })
                tokens.push(response.tokenData?.uuid as string)
            }

            const receivedTokens = [] as string[]
            await nativeIAmTokenGRPC.GetTokensForIdentity({
                activeFilter: GetTokensForIdentityRequest_ActiveFilter.ALL,
                identity,
                limit: 1,
                skip: 1,
                namespace: TEST_NAMESPACE_NAME
            }).forEach((r) => {
                receivedTokens.push(r.tokenData?.uuid as string)
            })

            expect(receivedTokens).toHaveLength(1)
            expect(receivedTokens[0]).toBe(tokens[1])
        })
        test("Total: 3. Skip: 0. Limit: 1.", async () => {
            const identity = "123"
            const tokensCount = 3
            const tokens = [] as string[]
            for(let i = 0; i < tokensCount; i++) {
                const response = await nativeIAmTokenGRPC.Create({
                    identity,
                    metadata: "",
                    namespace: TEST_NAMESPACE_NAME,
                    scopes: []
                })
                tokens.push(response.tokenData?.uuid as string)
            }

            const receivedTokens = [] as string[]
            await nativeIAmTokenGRPC.GetTokensForIdentity({
                activeFilter: GetTokensForIdentityRequest_ActiveFilter.ALL,
                identity,
                limit: 1,
                skip: 0,
                namespace: TEST_NAMESPACE_NAME
            }).forEach((r) => {
                receivedTokens.push(r.tokenData?.uuid as string)
            })

            expect(receivedTokens).toHaveLength(1)
            expect(receivedTokens[0]).toBe(tokens[2])
        })
        test("Total: 3. Skip: 1. Limit: 0.", async () => {
            const identity = "123"
            const tokensCount = 3
            const tokens = [] as string[]
            for(let i = 0; i < tokensCount; i++) {
                const response = await nativeIAmTokenGRPC.Create({
                    identity,
                    metadata: "",
                    namespace: TEST_NAMESPACE_NAME,
                    scopes: []
                })
                tokens.push(response.tokenData?.uuid as string)
            }

            const receivedTokens = [] as string[]
            await nativeIAmTokenGRPC.GetTokensForIdentity({
                activeFilter: GetTokensForIdentityRequest_ActiveFilter.ALL,
                identity,
                limit: 0,
                skip: 1,
                namespace: TEST_NAMESPACE_NAME
            }).forEach((r) => {
                receivedTokens.push(r.tokenData?.uuid as string)
            })

            expect(receivedTokens).toHaveLength(2)
            expect(receivedTokens[0]).toBe(tokens[1])
            expect(receivedTokens[1]).toBe(tokens[0])  
        })
    })

    test("Returns only active tokens if ActiveFilter seted to ONLY_ACTIVE", async () => {
        const identity = "123"
        const tokensCount = 5
        const tokens = [] as string[]
        for(let i = 0; i < tokensCount; i++) {
            const response = await nativeIAmTokenGRPC.Create({
                identity,
                metadata: "",
                namespace: TEST_NAMESPACE_NAME,
                scopes: []
            })
            const uuid = response.tokenData?.uuid as string
            if (i%2 == 0) {
                await nativeIAmTokenGRPC.Disable({
                    namespace: TEST_NAMESPACE_NAME,
                    uuid
                })
            } else {
                tokens.push(uuid)
            }
        }

        const receivedTokens = [] as string[]
        await nativeIAmTokenGRPC.GetTokensForIdentity({
            activeFilter: GetTokensForIdentityRequest_ActiveFilter.ONLY_ACTIVE,
            identity,
            limit: 0,
            skip: 0,
            namespace: TEST_NAMESPACE_NAME
        }).forEach((r) => {
            receivedTokens.push(r.tokenData?.uuid as string)
        })

        expect(tokens.sort()).toStrictEqual(receivedTokens.sort())
    })

    test("Returns only not active tokens if ActiveFilter seted to ONLY_NOT_ACTIVE", async () => {
        const identity = "123"
        const tokensCount = 5
        const tokens = [] as string[]
        for(let i = 0; i < tokensCount; i++) {
            const response = await nativeIAmTokenGRPC.Create({
                identity,
                metadata: "",
                namespace: TEST_NAMESPACE_NAME,
                scopes: []
            })
            const uuid = response.tokenData?.uuid as string
            if (i%2 == 0) {
                await nativeIAmTokenGRPC.Disable({ 
                    namespace: TEST_NAMESPACE_NAME,
                    uuid
                })
                tokens.push(uuid)
            }
        }

        const receivedTokens = [] as string[]
        await nativeIAmTokenGRPC.GetTokensForIdentity({
            activeFilter: GetTokensForIdentityRequest_ActiveFilter.ONLY_NOT_ACTIVE,
            identity,
            limit: 0,
            skip: 0,
            namespace: TEST_NAMESPACE_NAME
        }).forEach((r) => {
            receivedTokens.push(r.tokenData?.uuid as string)
        })

        expect(tokens.sort()).toStrictEqual(receivedTokens.sort())
    })
})