import { Status } from '@grpc/grpc-js/build/src/constants'

import { client as mongoClient } from '../../../../system/testing/tools/mongo'
import { client as cacheClient } from '../../../../system/testing/tools/cache'
import { RequestError as GRPCRequestError } from '../../../../system/libs/ts/grpc'
import { client as grpc } from '../../tools/namespace/grpc'

beforeAll(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

/**
 * @group native/namescape/delete/whitebox
 * @group whitebox
 */
describe("Whitebox", async () => {
    test("Value is deleted from the database", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        let entry = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').findOne({ name })
        expect(entry).not.toBeNull()
        await grpc.Delete({ name })
        entry = await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').findOne({ name })
        expect(entry).toBeNull()
    })

    test("Cache for namespace data is deleted from redis", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })
        let response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
        await grpc.Delete({ name })
        response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })

    test("Cache for namespaces list is deleted from redis", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        let response = await cacheClient.get("native_namespace_list")
        expect(response).not.toBeNull()
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        await grpc.Delete({ name })
        response = await cacheClient.get("native_namespace_list")
        expect(response).toBeNull()
    })

    test("Database for namespace is deleted from MongoDB server", async () => {
        const name = "testname"
        await grpc.Ensure({ name })

        const collectionName = "testingcollection"
        await mongoClient.db(`openerp_namespace_${name}`).collection<{ name: string }>(collectionName).insertOne({ name: "testinsert" })

        await grpc.Delete({ name })

        await mongoClient.db(`openerp_namespace_${name}`).listCollections({ name: collectionName }).forEach(() => {
            fail()
        })
    })
})

/**
 * @group native/namescape/delete/blackbox
 * @group blackbox
 */
 describe("Blackbox", async () => {
    test("Can\'t be accessed after deletion", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        expect(async () => await grpc.Get({ name, useCache: false })).not.toThrow()
        await grpc.Delete({ name })
        
        try {
            await grpc.Get({ name, useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
 })