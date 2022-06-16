import { Status } from '@grpc/grpc-js/build/src/constants'

import { client as mongoClient } from '../../../../tools/system/mongo'
import { RequestError as GRPCRequestError } from '../../../../tools/system/grpc'
import { client as cacheClient } from '../../../../tools/system/cache'

import { grpc } from '../../../../tools/native/namespace'

beforeAll(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

/**
 * @group native/namescape/get/whitebox
 * @group whitebox
 */
describe("Whitebox", async () => {
    test("Value is added to the cache on cache enabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name: name, useCache: true })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
    })

    test("Value is not added to the cache on cache disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name: name, useCache: false })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })

    test("Value is returned from cache when cache enabled and not from cache when disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })

        // Change value in database to contain wrong data without invalidating cache
        const newName = "newtestname"
        await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').updateOne({ name }, { "$set": { name: newName } })

        const cachedResponse = await grpc.Get({ name, useCache: true })
        expect(cachedResponse.namespace?.name).toBe(name)
        const uncachedResponse = await grpc.Get({ name, useCache: false })
        expect(uncachedResponse.namespace?.name).toBe(newName)
    })
})

/**
 * @group native/namescape/get/blackbox
 * @group blackbox
 */
describe("Blackbox", async () => {
    test("Get value same as created", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        const response = await grpc.Get({ name: name, useCache: false })
        expect(response.namespace?.name).toBe(name)
    })

    test("Get value for asked namespace", async () => {
        const namespaces = new Array<string>(16).map((_v, index) => `namespace${index}`)
        await Promise.all(namespaces.map((namespace) => grpc.Ensure({ name: namespace })))
        const results = await Promise.all(namespaces.map((namespace) => grpc.Get({ name: namespace, useCache: true })))
        const resultsNames = results.map(result => result.namespace?.name)
        expect(namespaces).toBe(resultsNames)
    })

    test("Returns error if not found", async () => {
        try {
            await grpc.Get({ name: "somevalue", useCache: false })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})