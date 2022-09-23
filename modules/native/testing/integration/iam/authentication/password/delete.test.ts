import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../../tools/namespace/grpc'
import { authenticationPasswordClient as nativeIAmAuthenticationPasswordGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `openbp_global`
const TEST_NAMESPACE_NAME = "authpasswordtestnamespace"
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
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_authentication_password').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_authentication_password').deleteMany({})
    } catch {}
    await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })
})

afterAll(async ()=>{
    await closeMongo()
    await cacheClient.flushall()
    await closeCache()
    await closeNativeIAM()
    await closeNativeNamespace()
})

/**
 * @group native/iam/identity/authenrication/password/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Deletes entry in global DB", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        await nativeIAmAuthenticationPasswordGRPC.Delete({
            identity,
            namespace: ""
        })

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_authentication_password").findOne({ identity })
        expect(entry).toBeNull()
    })
    test("Deletes entry in namespace DB", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "123"
        })

        await nativeIAmAuthenticationPasswordGRPC.Delete({
            identity,
            namespace: TEST_NAMESPACE_NAME
        })

        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_authentication_password").findOne({ identity })
        expect(entry).toBeNull()
    })
})

/**
 * @group native/iam/identity/authenrication/password/delete/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Non exist after deletion", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "123"
        })

        await nativeIAmAuthenticationPasswordGRPC.Delete({
            identity,
            namespace: TEST_NAMESPACE_NAME
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Exist({
            namespace: TEST_NAMESPACE_NAME,
            identity
        })
        expect(response.exist).toBeFalsy()
    })
    test("Several deletions of same password are ok", async () => {
        const identity1 = new ObjectId().toHexString()
        const identity2 = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity: identity1,
            namespace: "",
            password: "123"
        })

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity: identity2,
            namespace: "",
            password: "123"
        })

        for(let i = 0; i<5; i++) {
            await nativeIAmAuthenticationPasswordGRPC.Delete({
                identity: identity1,
                namespace: ""
            })

            const response1 = await nativeIAmAuthenticationPasswordGRPC.Exist({
                identity: identity1,
                namespace: ""
            })
            expect(response1.exist).toBeFalsy()
            const response2 = await nativeIAmAuthenticationPasswordGRPC.Exist({
                identity: identity2,
                namespace: ""
            })
            expect(response2.exist).toBeTruthy()
        }
    })
 })