import { randomBytes } from 'crypto'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { userClient as nativeActorUserGRPC, connect as connectToNativeActor, close as closeNativeActor } from '../../../tools/actor/grpc'
import { identityClient as nativeIAmIdentityGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'
import { ObjectID } from 'bson'
import { Status } from '@grpc/grpc-js/build/src/constants'

const GLOBAL_DB_NAME = 'openbp_global'

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeActor()
    await connectToNativeIAM()
    await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').deleteMany({})
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').deleteMany({})
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeActor()
    await closeNativeIAM()
})

/**
 * @group native/actor/user/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Clears cache on update", async () => {
        const createResponse = await nativeActorUserGRPC.Create({ login: randomBytes(10).toString('hex'), fullName: "", avatar: "", email: "" })
        const uuid = createResponse.user?.uuid as string
        const login = createResponse.user?.login as string
        const identity = createResponse.user?.identity as string

        await nativeActorUserGRPC.Get({ uuid, useCache: true })
        await nativeActorUserGRPC.GetByLogin({ login, useCache: true })
        await nativeActorUserGRPC.GetByIdentity({ identity, useCache: true })

        expect(await cacheClient.get(`native_actor_user_data_uuid_${uuid}`)).not.toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_login_${login}`)).not.toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_identity_${identity}`)).not.toBeNull()

        await nativeActorUserGRPC.Delete({ uuid })

        expect(await cacheClient.get(`native_actor_user_data_uuid_${uuid}`)).toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_login_${login}`)).toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_identity_${identity}`)).toBeNull()
    })

    test("Deletes data from the DB", async () => {
        const createResponse = await nativeActorUserGRPC.Create({ login: randomBytes(10).toString('hex'), fullName: "", avatar: "", email: "" })
        const uuid = createResponse.user?.uuid as string

        const countBeforeDelete = await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').count({"_id": ObjectID.createFromHexString(uuid)})
        expect(countBeforeDelete).toBe(1)

        await nativeActorUserGRPC.Delete({ uuid })

        const countAfterDelete = await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').count({"_id": ObjectID.createFromHexString(uuid)})
        expect(countAfterDelete).toBe(0)
    })

    test("Deletes identity", async () => {
        const createResponse = await nativeActorUserGRPC.Create({ login: randomBytes(10).toString('hex'), fullName: "", avatar: "", email: "" })
        const uuid = createResponse.user?.uuid as string
        const identity = createResponse.user?.identity as string

        await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: identity, useCache: false })

        await nativeActorUserGRPC.Delete({ uuid })

        try {
            await nativeIAmIdentityGRPC.Get({ namespace: "", uuid: identity, useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})

/**
 * @group native/actor/user/update/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("After deletion, 'get' operation failes with NOT_FOUND error", async () => { 
        const createResponse = await nativeActorUserGRPC.Create({ login: randomBytes(10).toString('hex'), fullName: "", avatar: "", email: "" })
        const uuid = createResponse.user?.uuid as string

        await nativeActorUserGRPC.Delete({ uuid })
        try {
            await nativeActorUserGRPC.Get({ uuid, useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })

    test("Several deletion of same user is OK", async () => { 
        const createResponse = await nativeActorUserGRPC.Create({ login: randomBytes(10).toString('hex'), fullName: "", avatar: "", email: "" })
        const uuid = createResponse.user?.uuid as string

        await nativeActorUserGRPC.Delete({ uuid })
        await nativeActorUserGRPC.Delete({ uuid })
        await nativeActorUserGRPC.Delete({ uuid })
    })

    test("Failes with INVALID_ARGUMENT if user uuid has bad format", async () => { 
        try {
            await nativeActorUserGRPC.Delete({ uuid: "somebadvalue" })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.INVALID_ARGUMENT) fail()
        }
    })
})