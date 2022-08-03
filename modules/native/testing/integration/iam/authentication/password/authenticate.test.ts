import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../../tools/namespace/grpc'
import { authenticationPasswordClient as nativeIAmAuthenticationPasswordGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const TEST_NAMESPACE_NAME = "authpasswordtestnamespace"
const NAMESPACE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`


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
    await closeCache()
    await closeNativeIAM()
    await closeNativeNamespace()
    await cacheClient.flushall()
})

/**
 * @group native/iam/identity/authenrication/password/authenticate/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Retuns false for bad password", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Authenticate({
            identity,
            namespace: "",
            password: "456"
        })

        expect(response.authenticated).toBeFalsy()
    })
    test("Returns false for bad identity", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Authenticate({
            identity: new ObjectId().toHexString(),
            namespace: "",
            password: "123"
        })

        expect(response.authenticated).toBeFalsy()
    })
    test("Returns false for bad namespace", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Authenticate({
            identity,
            namespace: TEST_NAMESPACE_NAME + "inv",
            password: "123"
        })

        expect(response.authenticated).toBeFalsy()
    })
    test("Returns true when everything is ok", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const response = await nativeIAmAuthenticationPasswordGRPC.Authenticate({
            identity,
            namespace: "",
            password: "123"
        })

        expect(response.authenticated).toBeTruthy()
    })
})