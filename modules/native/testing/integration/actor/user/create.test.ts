import { randomBytes } from 'crypto'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { userClient as nativeActorUserGRPC, connect as connectToNativeActor, close as closeNativeActor } from '../../../tools/actor/grpc'
import { identityClient as nativeIAmIdentityGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'
import { GetAllNamespacesResponse } from '../../../tools/namespace/proto/namespace'
import { ObjectID } from 'bson'
import { Status } from '@grpc/grpc-js/build/src/constants'

const GLOBAL_DB_NAME = 'openbp_global'

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
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
    await closeNativeNamespace()
    await closeNativeActor()
    await closeNativeIAM()
})

/**
 * @group native/actor/user/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in database", async () => { 
        const login = randomBytes(32).toString('hex')
        const response = await nativeActorUserGRPC.Create({
            login,
            fullName: "",
            avatar: "",
            email: ""
        })
        const id = ObjectID.createFromHexString(response.user?.uuid as string)
        
        const dbResponse = await mongoClient.db(GLOBAL_DB_NAME).collection("native_actor_user").findOne({"_id": id})
        expect(dbResponse).not.toBeNull()
    })

    test("Creates native_iam_identity", async () => { 
        const userResponse = await nativeActorUserGRPC.Create({
            login: randomBytes(10).toString('hex'),
            fullName: "",
            avatar: "",
            email: ""
        })
        const identity = userResponse.user?.identity as string

        const identityResponse = await nativeIAmIdentityGRPC.Get({
            namespace: "",
            useCache: false,
            uuid: identity
        })
        expect(identityResponse.identity?.uuid).toBe(identity)
    })
    
    test("Doesnt duplicate entries in DB with same login", async () => { 
        const login = randomBytes(16).toString('hex')
        await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
        try {
            await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.ALREADY_EXISTS) fail()
        }
        
        const dbResponse = await mongoClient.db(GLOBAL_DB_NAME).collection("native_actor_user").count({ login })
        expect(dbResponse).toBe(1)
    })

    test("Doesnt create new identity when failes with ALREADY_EXISTS error", async () => { 
        const login = randomBytes(16).toString('hex')
        await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})

        const identityCount = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").count({ })
        try {
            await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.ALREADY_EXISTS) fail()
        }
        const newIdentityCount = await mongoClient.db(GLOBAL_DB_NAME).collection("native_iam_identity").count({ })
        expect(newIdentityCount).toBe(identityCount)
    })
})

/**
 * @group native/actor/user/create/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Creates user if everything is OK", async () => { 
        const login = randomBytes(16).toString('hex')
        const createResponse = await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
        await nativeActorUserGRPC.Get({useCache: false, uuid: createResponse.user?.uuid as string})
    })

    test("Returns user data as the response of the request", async () => { 
        const login = randomBytes(10).toString('hex')
        const fullName = randomBytes(10).toString('hex')
        const avatar = randomBytes(10).toString('hex')
        const email = randomBytes(10).toString('hex') + "@mail.com"

        const userResponse = await nativeActorUserGRPC.Create({
            login,
            fullName,
            avatar,
            email
        })
        expect(userResponse.user?.login).toBe(login)
        expect(userResponse.user?.fullName).toBe(fullName)
        expect(userResponse.user?.avatar).toBe(avatar)
        expect(userResponse.user?.email).toBe(email)
    })

    test("Failes with ALREADY_EXISTS error if user with same login already exist", async () => { 
        const login = randomBytes(16).toString('hex')
        await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
        try {
            await nativeActorUserGRPC.Create({ login, fullName: "", avatar: "", email: ""})
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.ALREADY_EXISTS) fail()
        }
    })
})