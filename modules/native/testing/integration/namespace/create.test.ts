import { client as mongoClient } from '../../../../tools/system/mongo'
import { RequestError as GRPCRequestError } from '../../../../tools/system/grpc'
import { client as cacheClient } from '../../../../tools/system/cache'

import { grpc as namespaceGRPC } from '../../../../tools/native/namespace'
import { GetAllNamespacesResponse } from '../../../../tools/native/namespace/proto/namespace'

beforeAll(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

/**
 * @group native/namespace/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in database", async () => {
        const name = "customname"
        await namespaceGRPC.Ensure({ name })
        const entry = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').findOne({ name })
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })

    test("Doesnt duplicates in database multiple calls", async () => {
        const name = "customname"
        let count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(0)
        await namespaceGRPC.Ensure({ name })
        count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(1)
        await namespaceGRPC.Ensure({ name })
        count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(1)
    })

    test("Doesnt create entry in cache", async () => {
        const namespace1 = "namespace1"
        await namespaceGRPC.Ensure({ name: namespace1 })
        const existingKeys = await cacheClient.exists("native_namespace_list")
        expect(existingKeys).toBe(0)
    })

    test("Removes cache on list", async () => {
        await namespaceGRPC.Ensure({ name: "namespace1" })
        await namespaceGRPC.Get({ name: "namespace1", useCache: true })

        let existingKeys = await cacheClient.exists("native_namespace_list")
        expect(existingKeys).toBe(1)

        await namespaceGRPC.Ensure({ name: "namespace2" })
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
        const response = await namespaceGRPC.Ensure({ name })
        expect(response.namespace?.name === name).toBeTruthy()
    })

    test("Ok if already created", async () => {
        const name = "customname"
        await namespaceGRPC.Ensure({ name })
        const response = await namespaceGRPC.Ensure({ name })
        expect(response.namespace?.name).toBe(name)
    })

    test("Can create several and they will not overwrite", async () => {
        const name1 = "space1"
        await namespaceGRPC.Ensure({ name: name1 })
        const name2 = "space2"
        await namespaceGRPC.Ensure({ name: name2 })

        const allNamespaces = [] as Array<GetAllNamespacesResponse>
        await namespaceGRPC.GetAll({ useCache: false }).forEach((n) => { allNamespaces.push(n) })

        expect(allNamespaces.length).toBe(2)
        expect(allNamespaces.findIndex((x) => x.namespace?.name === name1) >= 0).toBeTruthy()
        expect(allNamespaces.findIndex((x) => x.namespace?.name === name2) >= 0).toBeTruthy()
    })

    describe("Fails with bad naming", async () => {
        const testf = async (name: string) => {
            try {
                await namespaceGRPC.Ensure({ name })
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