import { client as mongoClient } from '../../../../tools/system/mongo'
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
 * @group native/namespace/create/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Adds value to cache, when cache enabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Exists({ name: name, useCache: true })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).not.toBeNull()
    })
    test("Doesnt add value to cache, when cache disabled", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Exists({ name: name, useCache: false })
        const response = await cacheClient.get(`native_namespace_data_${name}`)
        expect(response).toBeNull()
    })
    test("Gets value from cache, when cache enabled. Doesnt get value from cache, when cache disabled.", async () => {
        const name = "testname"
        await grpc.Ensure({ name })
        await grpc.Get({ name, useCache: true })

        // Change value in database to contain wrong data without invalidating cache
        await mongoClient.db('openerp_global').collection<{ name: string }>('namespace').deleteOne({ name })

        const cachedResponse = await grpc.Exists({ name, useCache: true })
        expect(cachedResponse.exist).toBeTruthy()
        const uncachedResponse = await grpc.Exists({ name, useCache: false })
        expect(uncachedResponse.exist).toBeFalsy()
    })
})

/**
 * @group native/namescape/get/blackbox
 * @group blackbox
 */
 describe("Blackbox", async () => {
     test("Returs actual namespace status (cache disabled)", async () => {
        const name = "testname" 
        let response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeFalsy()
        await grpc.Ensure({ name })
        response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeTruthy()
        await grpc.Delete({ name })
        response = await grpc.Exists({ name, useCache: false })
        expect(response.exist).toBeFalsy()
     })

     test("Returs actual namespace status (cache enabled)", async () => {
        const name = "testname" 
        let response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeFalsy()
        await grpc.Ensure({ name })
        response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeTruthy()
        await grpc.Delete({ name })
        response = await grpc.Exists({ name, useCache: true })
        expect(response.exist).toBeFalsy()
     })
 })