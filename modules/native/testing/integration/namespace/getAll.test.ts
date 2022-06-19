import { Status } from '@grpc/grpc-js/build/src/constants'

import { client as mongoClient } from '../../../../system/testing/tools/mongo'
import { client as cacheClient } from '../../../../system/testing/tools/cache'

import { grpc } from '../../tools/namespace'

beforeAll(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

afterEach(async ()=>{
    await mongoClient.db('openerp_global').collection('namespace').drop()
    await cacheClient.flushall()
})

/**
 * @group native/namescape/getAll/whitebox
 * @group whitebox
 */
 describe("Whitebox", async () => {
    test("Value is added to the cache on cache enabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        const response = await cacheClient.get("native_namespace_list")
        expect(response).not.toBeNull()
    })

    test("Value is not added to the cache on cache disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)
        const response = await cacheClient.get("native_namespace_list")
        expect(response).toBeNull()
    })

    test("Value is returned from cache when cache enabled and not from cache when disabled", async () => {
        const namespaces = new Array<string>(10).map((_, index) => `testname${index}`)
        await Promise.all(namespaces.map((name) => grpc.Ensure({ name })))
        
        // load to cache
        await grpc.GetAll({ useCache: true }).forEach(() => undefined)

        // Change value in database to contain wrong data without invalidating cache
        const newName = "newtestname"
        await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').updateOne({ name: "testname1" }, { "$set": { name: newName } })

        const cachedResponse = new Array<string>()
        await grpc.GetAll({ useCache: true }).forEach((r) => cachedResponse.push(r.namespace?.name as string))
        if (cachedResponse.length !== namespaces.length) fail()
        cachedResponse.forEach((namespace) => expect(namespaces.indexOf(namespace)).toBeGreaterThanOrEqual(0))
        namespaces.forEach((namespace) => expect(cachedResponse.indexOf(namespace)).toBeGreaterThanOrEqual(0))
        expect(cachedResponse.indexOf(newName)).toBeLessThan(0)
        
        namespaces[0] = newName
        const uncachedResponse = new Array<string>()
        await grpc.GetAll({ useCache: true }).forEach((r) => uncachedResponse.push(r.namespace?.name as string))
        uncachedResponse.forEach((namespace) => expect(namespaces.indexOf(namespace)).toBeGreaterThanOrEqual(0))
        namespaces.forEach((namespace) => expect(uncachedResponse.indexOf(namespace)).toBeGreaterThanOrEqual(0))
    })
})

/**
 * @group native/namescape/getAll/blackbox
 * @group blackbox
 */
 describe("Blackbox", async () => {
    test("Returs actual namespaces list (cache disabled)", async () => {
        const name = "testname" 
        let response = new Array<string>()
        await grpc.GetAll({ useCache: false }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(0)
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: false }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(1)
        expect(response[0]).toBe(name)
        await grpc.Delete({ name })
        await grpc.GetAll({ useCache: false }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(0)
    })

    test("Returs actual namespaces list (cache enabled)", async () => {
        const name = "testname" 
        let response = new Array<string>()
        await grpc.GetAll({ useCache: true }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(0)
        await grpc.Ensure({ name })
        await grpc.GetAll({ useCache: true }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(1)
        expect(response[0]).toBe(name)
        await grpc.Delete({ name })
        await grpc.GetAll({ useCache: true }).forEach((r) => response.push(r.namespace?.name as string))
        expect(response.length).toBe(0)
    })
})