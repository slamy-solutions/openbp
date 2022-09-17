import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { identityClient as nativeIAmIdentityGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `openbp_global`
const TEST_NAMESPACE_NAME = "iamidentitytestnamespace"
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
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_identity').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_identity').deleteMany({})
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
 * @group native/iam/identity/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Deletes entry from global DB", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        let entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).not.toBeNull()

        await nativeIAmIdentityGRPC.Delete({
            namespace: "",
            uuid: id.toHexString()
        })

        entry = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).toBeNull()
    })

    test("Deletes entry from namespace DB", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        let entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).not.toBeNull()

        await nativeIAmIdentityGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString()
        })

        entry = await mongoClient.db(NAMESPACE_DB_NAME).collection("native_iam_identity").findOne<{ name: string }>({"_id": id})
        expect(entry).toBeNull()
    })

    test("Deletes entry from global cache", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({
            namespace: "",
            uuid: id.toHexString(),
            useCache: true
        })

        let existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(1)

        await nativeIAmIdentityGRPC.Delete({
            namespace: "",
            uuid: id.toHexString()
        })

        existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })

    test("Deletes entry from namespace cache", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: TEST_NAMESPACE_NAME
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString(),
            useCache: true
        })

        let existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(1)

        await nativeIAmIdentityGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: id.toHexString()
        })

        existResponse = await cacheClient.exists(`native_iam_identity_data_${TEST_NAMESPACE_NAME}_${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/identity/delete/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmIdentityGRPC.Delete({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Identity can not be geted after deletion", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({
            namespace: "",
            uuid: id.toHexString(),
            useCache: false
        })

        await nativeIAmIdentityGRPC.Delete({
            namespace: "",
            uuid: id.toHexString()
        })

        try {
            await nativeIAmIdentityGRPC.Get({
                namespace: "",
                uuid: id.toHexString(),
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })

    test("Several deletions are ok", async () => {
        const response = await nativeIAmIdentityGRPC.Create({
            name: "123123",
            initiallyActive: false,
            namespace: ""
        })
        const id = ObjectId.createFromHexString(response.identity?.uuid as string)

        for (let i = 0; i < 5; i += 1) {
            await nativeIAmIdentityGRPC.Delete({
                namespace: "",
                uuid: id.toHexString()
            })
        }
    })
    
    test("Deleting non existing identity is ok", async () => {
        await nativeIAmIdentityGRPC.Delete({
            namespace: "",
            uuid: new ObjectId().toHexString()
        })
        await nativeIAmIdentityGRPC.Delete({
            namespace: TEST_NAMESPACE_NAME,
            uuid: new ObjectId().toHexString()
        })
    })
})