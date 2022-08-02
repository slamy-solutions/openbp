import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { identityClient as nativeIAmIdentityGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const TEST_NAMESPACE_NAME = "iamidentitytestnamespace"
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
 * @group native/iam/identity/get/whitebox
 * @group whitebox
 */
 describe("Whitebox", () => {
    test("Gets value from cache if cache enabled", async () => {
        const name1 = randomBytes(32).toString("hex")
        const name2 = randomBytes(32).toString("hex")
        
        const createResponse = await nativeIAmIdentityGRPC.Create({
            namespace: "",
            name: name1,
            initiallyActive: false
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_identity').updateOne({"_id": id}, {"$set": {
            "name": name2
        }})

        const response = await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        expect(response.identity?.name).toBe(name1)
    })

    test("Gets value from DB if cache disabled", async () => {
        const name1 = randomBytes(32).toString("hex")
        const name2 = randomBytes(32).toString("hex")
        
        const createResponse = await nativeIAmIdentityGRPC.Create({
            namespace: "",
            name: name1,
            initiallyActive: false
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)

        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })

        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_identity').updateOne({"_id": id}, {"$set": {
            "name": name2,
        }})

        const response = await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: false })
        expect(response.identity?.name).toBe(name2)
    })

    test("Puts value in cache if cache enabled", async () => {
        const createResponse = await nativeIAmIdentityGRPC.Create({
            namespace: "",
            name: "",
            initiallyActive: false
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)
        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: true })
        const existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(1)
    })

    test("Doesnt put value in cache if cache disabled", async () => {
        const createResponse = await nativeIAmIdentityGRPC.Create({
            namespace: "",
            name: "",
            initiallyActive: false
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)
        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: id.toHexString(), useCache: false })
        const existResponse = await cacheClient.exists(`native_iam_identity_data__${id.toHexString()}`)
        expect(existResponse).toBe(0)
    })
})

/**
 * @group native/iam/identity/get/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Gets actual data", async () => {
        const name = randomBytes(32).toString("hex")
        const initiallyActive = true
        const createResponse = await nativeIAmIdentityGRPC.Create({
            namespace: TEST_NAMESPACE_NAME,
            name,
            initiallyActive
        })
        const id = ObjectId.createFromHexString(createResponse.identity?.uuid as string)
        const response = await nativeIAmIdentityGRPC.Get({ namespace: TEST_NAMESPACE_NAME, uuid: id.toHexString(), useCache: false })
        expect(response.identity?.namespace).toBe(TEST_NAMESPACE_NAME)
        expect(response.identity?.uuid).toBe(id.toHexString())
        expect(response.identity?.name).toBe(name)
        expect(response.identity?.active).toStrictEqual(initiallyActive)
        expect(response.identity?.policies).toStrictEqual([])
    })

    test("Failes with INVALID_ARGUMENT if uuid has bad format", async () => {
        try {
            await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: 'invalid',
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.INVALID_ARGUMENT)
        }
    })

    test("Failes with NOT_FOUND error if identity doesnt exist", async () => {
        try {
            await nativeIAmIdentityGRPC.Get({
                namespace: TEST_NAMESPACE_NAME,
                uuid: new ObjectId().toHexString(),
                useCache: false
            })
            fail()
        } catch(e) {
            expect((e as GRPCRequestError)?.code).toBe(Status.NOT_FOUND)
        }
    })
})