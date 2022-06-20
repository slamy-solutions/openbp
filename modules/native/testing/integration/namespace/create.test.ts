import { Status } from '@grpc/grpc-js/build/src/constants'

import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'
import { client as grpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../tools/namespace/grpc'
import { GetAllNamespacesResponse } from '../../tools/namespace/proto/namespace'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeNamespace()
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('namespace').deleteMany({})
    } catch (e) {
        if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) throw e
    }
    await cacheClient.flushall()
})

afterEach(async ()=>{
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('namespace').deleteMany({})
    } catch (e) {
        if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) throw e
    }
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeNamespace()
})

/**
 * @group native/namespace/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in database", async () => {
        const name = "customname"
        await grpc.Ensure({ name })
        const entry = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('namespace').findOne({ name })
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })

    test("Doesnt duplicates in database multiple calls", async () => {
        const name = "customname"
        let count = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(0)
        await grpc.Ensure({ name })
        count = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(1)
        await grpc.Ensure({ name })
        count = await mongoClient.db(GLOBAL_DB_NAME).collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(1)
    })

    test("Doesnt create entry in cache", async () => {
        const namespace1 = "namespace1"
        await grpc.Ensure({ name: namespace1 })
        const existingKeys = await cacheClient.exists("native_namespace_list")
        expect(existingKeys).toBe(0)
    })

    test("Removes cache on list", async () => {
        await grpc.Ensure({ name: "namespace1" })
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)

        let existingKeys = await cacheClient.exists("native_namespace_list")
        expect(existingKeys).toBe(1)

        await grpc.Ensure({ name: "namespace2" })
        existingKeys = await cacheClient.exists("native_namespace_list")
        expect(existingKeys).toBe(0)
    })
})

/**
 * @group native/namespace/create
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Returns created", async () => {
        const name = "customname"
        const response = await grpc.Ensure({ name })
        expect(response.namespace?.name === name).toBeTruthy()
    })

    test("Ok if already created", async () => {
        const name = "customname"
        await grpc.Ensure({ name })
        const response = await grpc.Ensure({ name })
        expect(response.namespace?.name).toBe(name)
    })

    test("Can create several and they will not overwrite", async () => {
        const name1 = "space1"
        await grpc.Ensure({ name: name1 })
        const name2 = "space2"
        await grpc.Ensure({ name: name2 })

        const allNamespaces = [] as Array<GetAllNamespacesResponse>
        await grpc.GetAll({ useCache: false }).forEach((n) => { allNamespaces.push(n) })

        expect(allNamespaces.length).toBe(2)
        expect(allNamespaces.findIndex((x) => x.namespace?.name === name1) >= 0).toBeTruthy()
        expect(allNamespaces.findIndex((x) => x.namespace?.name === name2) >= 0).toBeTruthy()
    })

    describe("Fails with bad naming", () => {
        const testf = async (name: string) => {
            try {
                await grpc.Ensure({ name })
                fail()
            } catch (e) {
                if ((e as GRPCRequestError)?.code !== 3) fail()
            }
        }
        const names = [
            "$name", "na%me", "%", "-", "=", "$in", "$set", " sfuisd", "ds sd", "asdvsdf&asd", "@", ""
        ] as Array<string>
        for (const name of names) {
            test(name, async () => await testf(name))
        }
    })
})