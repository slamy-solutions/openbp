import { randomBytes } from 'crypto'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { userClient as nativeActorUserGRPC, connect as connectToNativeActor, close as closeNativeActor } from '../../../tools/actor/grpc'
import { ObjectID } from 'bson'
import { Status } from '@grpc/grpc-js/build/src/constants'

const GLOBAL_DB_NAME = 'openbp_global'

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeActor()
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
})

/**
 * @group native/actor/user/update/whitebox
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

        await nativeActorUserGRPC.Update({ uuid, avatar: "1", email: "", fullName: "name", login})

        expect(await cacheClient.get(`native_actor_user_data_uuid_${uuid}`)).toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_login_${login}`)).toBeNull()
        expect(await cacheClient.get(`native_actor_user_data_identity_${identity}`)).toBeNull()
    })

    test("Updates data in DB", async () => {
        const login = randomBytes(10).toString('hex')
        const fullName = randomBytes(10).toString('hex')
        const avatar = randomBytes(10).toString('hex')
        const email = randomBytes(10).toString('hex') + "@mail.com"

        const createResponse = await nativeActorUserGRPC.Create({
            login,
            fullName,
            avatar,
            email
        })
        const uuid = createResponse.user?.uuid as string

        const newLogin = randomBytes(10).toString('hex')
        const newFullName = randomBytes(10).toString('hex')
        const newAvatar = randomBytes(10).toString('hex')
        const newEmail = randomBytes(10).toString('hex') + "@mail.com"
        await nativeActorUserGRPC.Update({
            uuid,
            login: newLogin,
            fullName: newFullName,
            avatar: newAvatar,
            email: newEmail
        })

        interface UserInDB {
            login: string
            fullName: string
            avatar: string
            email: string
        }

        const dbEntry = await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').findOne<UserInDB>({"_id": ObjectID.createFromHexString(uuid)})
        expect(dbEntry?.login).toBe(newLogin)
        expect(dbEntry?.fullName).toBe(newFullName)
        expect(dbEntry?.avatar).toBe(newAvatar)
        expect(dbEntry?.email).toBe(newEmail)
    })
})

/**
 * @group native/actor/user/update/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns new data as result of the update", async () => {
        const login = randomBytes(10).toString('hex')
        const fullName = randomBytes(10).toString('hex')
        const avatar = randomBytes(10).toString('hex')
        const email = randomBytes(10).toString('hex') + "@mail.com"

        const createResponse = await nativeActorUserGRPC.Create({
            login,
            fullName,
            avatar,
            email
        })
        const uuid = createResponse.user?.uuid as string

        const newLogin = randomBytes(10).toString('hex')
        const newFullName = randomBytes(10).toString('hex')
        const newAvatar = randomBytes(10).toString('hex')
        const newEmail = randomBytes(10).toString('hex') + "@mail.com"
        const updateResponse = await nativeActorUserGRPC.Update({
            uuid,
            login: newLogin,
            fullName: newFullName,
            avatar: newAvatar,
            email: newEmail
        })

        expect(updateResponse.user?.uuid).toBe(uuid)
        expect(updateResponse.user?.login).toBe(newLogin)
        expect(updateResponse.user?.fullName).toBe(newFullName)
        expect(updateResponse.user?.avatar).toBe(newAvatar)
        expect(updateResponse.user?.email).toBe(newEmail)
    })

    test("After update, 'get' operation returns updated data", async () => {
        const login = randomBytes(10).toString('hex')
        const fullName = randomBytes(10).toString('hex')
        const avatar = randomBytes(10).toString('hex')
        const email = randomBytes(10).toString('hex') + "@mail.com"

        const createResponse = await nativeActorUserGRPC.Create({
            login,
            fullName,
            avatar,
            email
        })
        const uuid = createResponse.user?.uuid as string

        const newLogin = randomBytes(10).toString('hex')
        const newFullName = randomBytes(10).toString('hex')
        const newAvatar = randomBytes(10).toString('hex')
        const newEmail = randomBytes(10).toString('hex') + "@mail.com"
        await nativeActorUserGRPC.Update({
            uuid,
            login: newLogin,
            fullName: newFullName,
            avatar: newAvatar,
            email: newEmail
        })

        const getResponse = await nativeActorUserGRPC.Get({
            useCache: false,
            uuid
        })

        expect(getResponse.user?.uuid).toBe(uuid)
        expect(getResponse.user?.login).toBe(newLogin)
        expect(getResponse.user?.fullName).toBe(newFullName)
        expect(getResponse.user?.avatar).toBe(newAvatar)
        expect(getResponse.user?.email).toBe(newEmail)
    })

    test("Fails with NOT_FOUND error if user for update (uuid) doesnt exist", async () => {
        try {
            await nativeActorUserGRPC.Update({
                uuid: new ObjectID().toHexString(),
                login: randomBytes(10).toString('hex'),
                fullName: randomBytes(10).toString('hex'),
                avatar: randomBytes(10).toString('hex'),
                email: randomBytes(10).toString('hex')
            })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })

    test("Fails with INVALID_ARGUMENT error if uuid has bad format", async () => {
        try {
            await nativeActorUserGRPC.Update({
                uuid: "badformat",
                login: randomBytes(10).toString('hex'),
                fullName: randomBytes(10).toString('hex'),
                avatar: randomBytes(10).toString('hex'),
                email: randomBytes(10).toString('hex')
            })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.INVALID_ARGUMENT) fail()
        }
    })

    test("Fails with ALREADY_EXISTS error if new login for user already used by other user", async () => {
        const create1Response = await nativeActorUserGRPC.Create({ login: 'login1', avatar: "", email: "", fullName: "" })
        const create2Response = await nativeActorUserGRPC.Create({ login: 'login2', avatar: "", email: "", fullName: "" })

        try {
            await nativeActorUserGRPC.Update({
                uuid: create1Response.user?.uuid as string,
                login: create2Response.user?.login as string,
                fullName: randomBytes(10).toString('hex'),
                avatar: randomBytes(10).toString('hex'),
                email: randomBytes(10).toString('hex')
            })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.ALREADY_EXISTS) fail()
        }
    })
})