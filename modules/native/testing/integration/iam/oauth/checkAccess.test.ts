import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { sign as signJWT, decode as decodeJWT } from 'jsonwebtoken'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
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
import { CheckAccessResponse_Status, CreateTokenWithPasswordResponse_Status, RefreshTokenResponse_Status } from '../../../tools/iam/proto/oauth'
import { RefreshResponse_Status, Scope, ValidateResponse_Status } from '../../../tools/iam/proto/token'
import { Policy } from '../../../tools/iam/proto/policy'
import { ObjectID } from 'bson'

import { makeCases } from './authCases'

const GLOBAL_DB_NAME = `openbp_global`
const TEST_NAMESPACE_NAME = "iamoauthtestnamespace"
const NAMESPACE_DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`


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

/*
        // Provided token allows to access scopes
        OK = 0;

        // Received token has bad format or its signature doesnt match
        TOKEN_INVALID = 1;
        // Most probably token was deleted after its creation
        TOKEN_NOT_FOUND = 2;
        // Token was manually disabled
        TOKEN_DISABLED = 3;
        // Token expired
        TOKEN_EXPIRED = 4;

        // Token has not enought privileges to access specified scopes
        UNAUTHORIZED = 5;

*/

/**
 * @group native/iam/oauth/checkAccess/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns TOKEN_INVALID if token has bad format", async () => { 
        const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
            accessToken: "somethingbad",
            scopes: []
        })
        expect(checkAccessResponse.status).toBe(CheckAccessResponse_Status.TOKEN_INVALID)
    })

    test("Returns TOKEN_INVALID if token was signed by bad key", async () => {
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
        

        const tokenPayload = decodeJWT(tokenResponse.accessToken)
        const token = signJWT(tokenPayload as string, randomBytes(32).toString("hex"))

        const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
            accessToken: token,
            scopes: []
        })

        expect(checkAccessResponse.status).toBe(CheckAccessResponse_Status.TOKEN_INVALID)
    })

    test("Returns TOKEN_NOT_FOUND if token was deleated", async () => {
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

        const tokenDataResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })

        await nativeIAmTokenGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: tokenDataResponse.tokenData?.uuid as string
        })

        const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
            accessToken: tokenResponse.accessToken,
            scopes: []
        })

        expect(checkAccessResponse.status).toBe(CheckAccessResponse_Status.TOKEN_NOT_FOUND)
    })

    test("Returns TOKEN_NOT_FOUND if token namespace was deleated", async () => { 
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

        const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
            accessToken: tokenResponse.accessToken,
            scopes: []
        })

        expect(checkAccessResponse.status).toBe(CheckAccessResponse_Status.TOKEN_NOT_FOUND)
    })
    
    test("Returns TOKEN_DISABLED if token was manually disabled", async () => { 
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

        const tokenDataResponse = await nativeIAmTokenGRPC.RawGet({
            token: tokenResponse.accessToken,
            useCache: false
        })

        await nativeIAmTokenGRPC.Disable({
            namespace: TEST_NAMESPACE_NAME,
            uuid: tokenDataResponse.tokenData?.uuid as string
        })

        const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
            accessToken: tokenResponse.accessToken,
            scopes: []
        })

        expect(checkAccessResponse.status).toBe(CheckAccessResponse_Status.TOKEN_DISABLED)
    })

    describe("Authorization tests", () => {
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
                for(const policy of testCase.has) {
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

                const checkAccessResponse = await nativeIAmOAuthGRPC.CheckAccess({
                    accessToken: tokenResponse.accessToken as string,
                    scopes: testCase.requests
                })

                const expectedStatus = testCase.authorized ? CheckAccessResponse_Status.OK : CheckAccessResponse_Status.UNAUTHORIZED
                expect(checkAccessResponse.status).toBe(expectedStatus)
            })
        }
    })
})