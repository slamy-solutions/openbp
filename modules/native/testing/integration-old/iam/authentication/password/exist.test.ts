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
 * @group native/iam/identity/authenrication/password/authenticate/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Returns true when exist", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Exist({
            identity,
            namespace: ""
        })

        expect(response.exist).toBeTruthy()
    })
    test("Returns false when bad namespace", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Exist({
            identity,
            namespace: TEST_NAMESPACE_NAME + "invvv"
        })

        expect(response.exist).toBeFalsy()
    })
    test("Returns false when wrong namespace", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Exist({
            identity,
            namespace: TEST_NAMESPACE_NAME
        })

        expect(response.exist).toBeFalsy()
    })
    test("Returns false when bad identity", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Exist({
            identity: new ObjectId().toHexString(),
            namespace: ""
        })

        expect(response.exist).toBeFalsy()
    })
})