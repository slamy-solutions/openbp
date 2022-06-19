import { client as mongoClient } from '../../../../system/testing/tools/mongo'
import { client as cacheClient } from '../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'
import { client as grpc } from '../../tools/namespace/grpc'
import { GetAllNamespacesResponse } from '../../tools/namespace/proto/namespace'

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
        await grpc.Ensure({ name })
        const entry = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').findOne({ name })
        expect(entry).not.toBeNull()
        expect(entry?.name).toBe(name)
    })

    test("Doesnt duplicates in database multiple calls", async () => {
        const name = "customname"
        let count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(0)
        await grpc.Ensure({ name })
        count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
        expect(count).toBe(1)
        await grpc.Ensure({ name })
        count = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').countDocuments()
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
        await grpc.Get({ name: "namespace1", useCache: true })

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

    describe("Fails with bad naming", async () => {
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