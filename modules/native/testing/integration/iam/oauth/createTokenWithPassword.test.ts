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
import { CreateTokenWithPasswordResponse_Status } from '../../../tools/iam/proto/oauth'
import { Scope, ValidateResponse_Status } from '../../../tools/iam/proto/token'
import { Policy } from '../../../tools/iam/proto/policy'
import { ObjectID } from 'bson'

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

async function createRandomPoliciesForIdentity(identity: string, namespace: string = TEST_NAMESPACE_NAME) {
    interface PolicyCreationData {
        namespace: string
        resources: Array<string>
        actions: Array<string>
    }

    const policies = new Array<PolicyCreationData>(6).fill({
        namespace: "",
        resources: [],
        actions: []
    }).map(() => {
        return {
            namespace,
            resources: new Array(4).fill("").map(() => randomBytes(32).toString('hex')),
            actions: new Array(4).fill("").map(() => randomBytes(32).toString('hex'))
        }
    })

    const policyResponses = await Promise.all(policies.map((p) => {
        return nativeIAmPolicyGRPC.Create({
            name: randomBytes(16).toString('hex'),
            namespace: p.namespace,
            resources: p.resources,
            actions: p.actions
        })
    }))
    await Promise.all(policyResponses.map((p) => {
        nativeIAmIdentityGRPC.AddPolicy({
            identityNamespace: namespace,
            identityUUID: identity,
            policyNamespace: namespace,
            policyUUID: p.policy?.uuid as string
        })
    }))

    return policyResponses.map((p)=> { return p.policy as Policy })
}

/**
 * @group native/iam/oauth/createTokenWithPassword/blackbox
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
        
        expect(tokenResponse.status).toBe(CreateTokenWithPasswordResponse_Status.OK)
    })

    test("On success returns valid access and refresh tokens", async () => {
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

        const accessTokenResponse = await nativeIAmTokenGRPC.Validate({
            token: tokenResponse.accessToken,
            useCache: false
        })
        expect(accessTokenResponse.status).toBe(ValidateResponse_Status.OK)
        const refreshTokenResponse = await nativeIAmTokenGRPC.Validate({
            token: tokenResponse.accessToken,
            useCache: false
        })
        expect(refreshTokenResponse.status).toBe(ValidateResponse_Status.OK)
    })

    test("On success returns tokens with same scopes as was in the request", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME,
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })

        const policies = await createRandomPoliciesForIdentity(identity, TEST_NAMESPACE_NAME)

        let scopes = [policies[1], policies[4], policies[5]]
        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: scopes.map((s) => {
                return {
                    namespace: s.namespace,
                    resources: s.resources,
                    actions: s.actions
                }
            })
        })

        const getTokenResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })
        let getedScopes = getTokenResponse.tokenData?.scopes as Scope[]


        scopes = scopes.sort((s1, s2) => s1.resources[0].localeCompare(s2.resources[0]))
        getedScopes = getedScopes.sort((s1, s2) => s1.resources[0].localeCompare(s2.resources[0]))
        
        expect(getedScopes).toHaveLength(scopes.length)
        for (let i = 0; i < scopes.length; i++) {
            expect(scopes[i].namespace).toBe(getedScopes[i].namespace)
            expect(scopes[i].resources.sort()).toStrictEqual(getedScopes[i].resources.sort())
            expect(scopes[i].actions.sort()).toStrictEqual(getedScopes[i].actions.sort())
        }
    })

    test("On success returns tokens with all posible scopes if in request list of scopes was empty", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME,
        })
        const identity = identityResponse.identity?.uuid as string

        const password = randomBytes(16).toString('hex')
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password
        })

        let policies = await createRandomPoliciesForIdentity(identity, TEST_NAMESPACE_NAME)

        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })

        const getTokenResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })
        let getedScopes = getTokenResponse.tokenData?.scopes as Scope[]


        policies = policies.sort((s1, s2) => s1.resources[0].localeCompare(s2.resources[0]))
        getedScopes = getedScopes.sort((s1, s2) => s1.resources[0].localeCompare(s2.resources[0]))
        
        expect(getedScopes).toHaveLength(policies.length)
        for (let i = 0; i < policies.length; i++) {
            expect(policies[i].namespace).toBe(getedScopes[i].namespace)
            expect(policies[i].resources.sort()).toStrictEqual(getedScopes[i].resources.sort())
            expect(policies[i].actions.sort()).toStrictEqual(getedScopes[i].actions.sort())
        }
    })

    test("On success returns created token has metadata provided by this request", async () => {
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
        
        const metadata = randomBytes(64).toString('hex')
        const tokenCreateResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata,
            password,
            scopes: []
        })

        const getTokenResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenCreateResponse.accessToken,
            useCache: false
        })

        expect(getTokenResponse.tokenData?.creationMetadata).toBe(metadata)
    })

    test("Returns CREDENTIALS_INVALID status if password authentication method wasnt set for this identity", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string
        
        const tokenCreateResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password: "notcreated",
            scopes: []
        })

        expect(tokenCreateResponse.status).toBe(CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID)
    })

    test("Returns CREDENTIALS_INVALID status if password for this identity is not valid", async () => {
        const identityResponse = await nativeIAmIdentityGRPC.Create({
            name: randomBytes(16).toString('hex'),
            initiallyActive: true,
            namespace: TEST_NAMESPACE_NAME
        })
        const identity = identityResponse.identity?.uuid as string
        
        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "password"
        })

        const tokenCreateResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password: "invalid",
            scopes: []
        })

        expect(tokenCreateResponse.status).toBe(CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID)
    })

    test("Returns CREDENTIALS_INVALID status if there is no such identity", async () => {
        const tokenCreateResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity: new ObjectID().toHexString(),
            metadata: "",
            password: "invalid",
            scopes: []
        })

        expect(tokenCreateResponse.status).toBe(CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID)
    })

    test("Returns CREDENTIALS_INVALID status if there is no such namespace", async () => {
        await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })

        const tokenCreateResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity: new ObjectID().toHexString(),
            metadata: "",
            password: "invalid",
            scopes: []
        })

        expect(tokenCreateResponse.status).toBe(CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID)
    })

    test("Returns IDENTITY_NOT_ACTIVE status if identity is not active", async () => {
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
        
        await nativeIAmIdentityGRPC.SetActive({
            namespace: TEST_NAMESPACE_NAME,
            uuid: identity,
            active: false
        })

        const tokenResponse = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
            namespace: TEST_NAMESPACE_NAME,
            identity,
            metadata: "",
            password,
            scopes: []
        })
        
        expect(tokenResponse.status).toBe(CreateTokenWithPasswordResponse_Status.IDENTITY_NOT_ACTIVE)
    })

    describe("Authorization cases", () => {
        for(const el of makeCases(TEST_NAMESPACE_NAME)) {
            const testCase = el
            test(testCase.name, async () => {
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

                for (const policy of testCase.has) {
                    const policyCreateResponse = await nativeIAmPolicyGRPC.Create({
                        name: "",
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
                }

                const response = await nativeIAmOAuthGRPC.CreateTokenWithPassword({
                    identity,
                    metadata: "",
                    namespace: TEST_NAMESPACE_NAME,
                    password,
                    scopes: testCase.requests
                })

                expect(response.status == CreateTokenWithPasswordResponse_Status.OK).toBe(testCase.authorized)
            })
        }
    })
})