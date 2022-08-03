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
 * @group native/iam/identity/authenrication/password/createOrUpdate/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Doesnt store passwords in plain form", async () => {
        const identity = new ObjectId().toHexString()
        const password = "123"

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password
        })

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_authentication_password").findOne<{ password: string }>({ identity })
        expect(entry).not.toBeNull()
        expect(entry?.password).not.toBeNull()
        expect(entry?.password).not.toBeUndefined()
        expect(entry?.password.length).toBeGreaterThan(0)
        expect(entry?.password).not.toBe(password)
    })
    test("Creates entry in global DB on creation", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_authentication_password").findOne<{ identity: string }>({ identity })
        expect(entry).not.toBeNull()
        expect(entry?.identity).toBe(identity)
    })
    test("Creates entry in namespace DB on creation", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "123"
        })

        const entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_authentication_password").findOne<{ identity: string }>({ identity })
        expect(entry).not.toBeNull()
        expect(entry?.identity).toBe(identity)
    })
    test("Updates entry in global DB on update", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        let entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_authentication_password").findOne<{ _id: ObjectId, password: string }>({ identity })
        const _id = entry?._id as ObjectId
        const oldPassword = entry?.password as string

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "12345"
        })

        entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_authentication_password").findOne<{ _id: ObjectId, password: string }>({ _id })
        expect(entry?.password).not.toBe(oldPassword)
    })
    test("Updates entry in namespace DB on update", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "123"
        })

        let entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_authentication_password").findOne<{ _id: ObjectId, password: string }>({ identity })
        const _id = entry?._id as ObjectId
        const oldPassword = entry?.password as string

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: TEST_NAMESPACE_NAME,
            password: "12345"
        })

        entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_authentication_password").findOne<{ _id: ObjectId, password: string }>({ _id })
        expect(entry?.password).not.toBe(oldPassword)
    })
})

/**
 * @group native/iam/identity/authenrication/password/createOrUpdate/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Exist after creation", async () => {
        const identity = new ObjectId().toHexString()

        await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
            identity,
            namespace: "",
            password: "123"
        })

        const existResponse = await nativeIAmAuthenticationPasswordGRPC.Exist({
            namespace: "",
            identity
        })
        expect(existResponse.exist).toBeTruthy()
    })
    test("Failes with FAILED_PRECONDITION if namespace doesnt exist", async () => {
        try {
            await nativeIAmAuthenticationPasswordGRPC.CreateOrUpdate({
                identity: new ObjectId().toHexString(),
                password: "123",
                namespace: NAMESPACE_DB_NAME + "invalid"
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.FAILED_PRECONDITION)
        }
    })
})