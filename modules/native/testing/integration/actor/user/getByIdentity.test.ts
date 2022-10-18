import { randomBytes } from 'crypto'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { userClient as nativeActorUserGRPC, connect as connectToNativeActor, close as closeNativeActor } from '../../../tools/actor/grpc'
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
 * @group native/actor/user/getByIdentity/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Gets data from cache if cache enabled. Doesnt get data from cache if cache disabled", async () => {
        const avatar = randomBytes(16).toString("hex")
        const createResponse = await nativeActorUserGRPC.Create({
            login: randomBytes(10).toString("hex"),
            fullName: "",
            avatar,
            email: ""
        })
        const login = createResponse.user?.login as string
        const identity = createResponse.user?.identity as string

        await nativeActorUserGRPC.GetByIdentity({
            useCache: true,
            identity
        })

        const newAvatar = randomBytes(16).toString("hex")
        const updateResponse = await mongoClient.db(GLOBAL_DB_NAME).collection('native_actor_user').updateOne({ login }, {"$set": { avatar: newAvatar }})
        expect(updateResponse.modifiedCount).toBe(1)

        const cachedResponse = await nativeActorUserGRPC.GetByIdentity({
            useCache: true,
            identity
        })
        expect(cachedResponse.user?.avatar).toBe(avatar)

        const notCachedResponse = await nativeActorUserGRPC.GetByIdentity({
            useCache: false,
            identity
        })
        expect(notCachedResponse.user?.avatar).toBe(newAvatar)
    })
})

/**
 * @group native/actor/user/getByIdentity/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Gets actual user data", async () => {
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
        const identity = createResponse.user?.identity as string

        const getResponse = await nativeActorUserGRPC.GetByIdentity({
            useCache: false,
            identity
        })

        expect(getResponse.user?.uuid).toBe(uuid)
        expect(getResponse.user?.login).toBe(login)
        expect(getResponse.user?.fullName).toBe(fullName)
        expect(getResponse.user?.avatar).toBe(avatar)
        expect(getResponse.user?.email).toBe(email)
    })

    test("Cached response gets actual user data", async () => {
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
        const identity = createResponse.user?.identity as string

        await nativeActorUserGRPC.GetByIdentity({ useCache: true, identity })
        const getResponse = await nativeActorUserGRPC.GetByIdentity({
            useCache: true,
            identity
        })

        expect(getResponse.user?.uuid).toBe(uuid)
        expect(getResponse.user?.login).toBe(login)
        expect(getResponse.user?.fullName).toBe(fullName)
        expect(getResponse.user?.avatar).toBe(avatar)
        expect(getResponse.user?.email).toBe(email)
    })

    test("Fails with NOT_FOUND error if user with specified identity doesnt exist",  async () => { 
        try {
            await nativeActorUserGRPC.GetByIdentity({ useCache: false, identity: randomBytes(16).toString('hex')})
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})