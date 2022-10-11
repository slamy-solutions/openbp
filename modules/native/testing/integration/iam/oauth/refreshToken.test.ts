import { randomBytes } from 'crypto'

import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import {
    tokenClient as nativeIAmTokenGRPC,
    oauthClient as nativeIAmOAuthGRPC,
    identityClient as nativeIAmIdentityGRPC,
    policyClient as nativeIAmPolicyGRPC,
    authenticationPasswordClient as nativeIAmAuthenticationPasswordGRPC,
    connect as connectToNativeIAM,
    close as closeNativeIAM
} from '../../../tools/iam/grpc'
import { RefreshTokenResponse_Status } from '../../../tools/iam/proto/oauth'
import { ValidateResponse_Status } from '../../../tools/iam/proto/token'
import { Policy } from '../../../tools/iam/proto/policy'

import { makeCases } from './authCases'

const TEST_NAMESPACE_NAME = "iamoauthtestnamespace"


beforeAll(async ()=>{
    await connectToNativeIAM()
    await connectToNativeNamespace()
    
})

beforeEach(async () => {
    await nativeNamespaceGRPC.Ensure({ name: TEST_NAMESPACE_NAME })
})

afterEach(async ()=>{
    await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })
})

afterAll(async ()=>{
    await closeNativeIAM()
    await closeNativeNamespace()
})

/**
 * @group native/iam/oauth/refreshToken/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns OK status if everything is ok", async () => { 
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.OK)
    })

    test("On success returns valid refresh tokens", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        const validationResponse = await nativeIAmTokenGRPC.Validate({
            token: refreshResponse.accessToken,
            useCache: false
        })

        expect(validationResponse.status).toBe(ValidateResponse_Status.OK)
    })

    test("Returns TOKEN_INVALID status if token has bad format", async () => {
        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: "asdasd"
        })
        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.TOKEN_INVALID)
    })

    test("Returns TOKEN_NOT_FOUND status if token was deleted", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const tokenGetResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })

        await nativeIAmTokenGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: tokenGetResponse.tokenData?.uuid as string
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.TOKEN_NOT_FOUND)
    })

    test("Returns TOKEN_NOT_FOUND status if token namespace was deleted", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        await nativeNamespaceGRPC.Delete({
            name: TEST_NAMESPACE_NAME
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.TOKEN_NOT_FOUND)
    })

    test("Returns TOKEN_DISABLED status if token namespace was deleted", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const tokenGetResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })

        await nativeIAmTokenGRPC.Disable({
            namespace: TEST_NAMESPACE_NAME,
            uuid: tokenGetResponse.tokenData?.uuid as string
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.TOKEN_DISABLED)
    })

    test("Returns TOKEN_IS_NOT_REFRESH_TOKEN status if token is not a refresh token", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.accessToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.TOKEN_IS_NOT_REFRESH_TOKEN)
    })

    test("Returns IDENTITY_NOT_FOUND status if identity was deleted", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        await nativeIAmIdentityGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: identity
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.IDENTITY_NOT_FOUND)
    })

    test("Returns IDENTITY_NOT_ACTIVE status if identity was mannualy disabled", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })
        
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: identity,
            active: false
        })

        const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
            refreshToken: tokenResponse.refreshToken
        })

        expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.IDENTITY_NOT_ACTIVE)
    })

    describe("Revalidates policies on refresh", () => {
        for(const el of makeCases(TEST_NAMESPACE_NAME)) {
            const testCase = el

            test(testCase.name, async () => {
                const identityResponse = await nativeIAmIdentityGRPC.Create({
                    name: randomBytes(16).toString('hex'),
                    initiallyActive: true,
                    namespace: TEST_NAMESPACE_NAME
                })
                const identity = identityResponse.identity?.uuid as string

                const hasPolicies = [] as Array<{namespace: string, uuid: string}>
                for(const policy of testCase.requests) {
                    const policyCreateResponse = await nativeIAmPolicyGRPC.Create({
                        name: randomBytes(16).toString('hex'),
                        namespace: policy.namespace,
                        resources: policy.resources,
                        actions: policy.actions
                    })
                    await nativeIAmIdentityGRPC.AddPolicy({
                        identityNamespace: TEST_NAMESPACE_NAME,
                        identityUUID: identity,
                        policyNamespace: policy.namespace,
                        policyUUID: policyCreateResponse.policy?.uuid as string
                    })
                    hasPolicies.push(policyCreateResponse.policy as Policy)
                }

                const password = randomBytes(16).toString('hex')
                await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
                    identity,
                    namespace: TEST_NAMESPACE_NAME,
                    password
                })
                
                const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
                    namespace: TEST_NAMESPACE_NAME,
                    identity,
                    metadata: "",
                    password,
                    scopes: []
                })

                for(const policy of hasPolicies) {
                    await nativeIAmIdentityGRPC.RemovePolicy({
                        identityNamespace: TEST_NAMESPACE_NAME,
                        identityUUID: identity,
                        policyNamespace: policy.namespace,
                        policyUUID: policy.uuid
                    })
                    await nativeIAmPolicyGRPC.Delete(policy)
                }

                for(const scope of testCase.has) {
                    const policyCreateResponse = await nativeIAmPolicyGRPC.Create({
                        name: randomBytes(16).toString('hex'),
                        namespace: scope.namespace,
                        resources: scope.resources,
                        actions: scope.actions
                    })
                    await nativeIAmIdentityGRPC.AddPolicy({
                        identityNamespace: TEST_NAMESPACE_NAME,
                        identityUUID: identity,
                        policyNamespace: scope.namespace,
                        policyUUID: policyCreateResponse.policy?.uuid as string
                    })
                }

                const refreshResponse = await nativeIAmOAuthGRPC.RefreshToken({
                    refreshToken: tokenResponse.refreshToken
                })
        
                if (testCase.authorized) {
                    expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.OK)
                } else {
                    expect(refreshResponse.status).toBe(RefreshTokenResponse_Status.IDENTITY_UNAUTHENTICATED)
                }
            })
        }
    })
})