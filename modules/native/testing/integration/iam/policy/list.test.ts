import { randomBytes } from 'crypto'
import { ObjectId } from 'mongodb'
import { Status } from '@grpc/grpc-js/build/src/constants'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { client as nativeNamespaceGRPC, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { policyClient as nativeIAmPolicyGRPC, connect as connectToNativeIAM, close as closeNativeIAM } from '../../../tools/iam/grpc'

const GLOBAL_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const TEST_NAMESPACE_NAME = "liampolicytestnamespace"
const NAMESPACE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`


beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeIAM()
    await connectToNativeNamespace()
    
})

beforeEach(async () => {
    await nativeNamespaceGRPC.Ensure({ name: TEST_NAMESPACE_NAME })
})

afterEach(async ()=>{
    try {
        await mongoClient.db(GLOBAL_DB_NAME).collection('native_iam_policy').deleteMany({})
    } catch {}
    try {
        await mongoClient.db(NAMESPACE_DB_NAME).collection('native_iam_policy').deleteMany({})
    } catch {}
    await nativeNamespaceGRPC.Delete({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeIAM()
    await closeNativeNamespace()
})

/**
 * @group native/iam/policy/list/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Returns full list of policies", async () => {
        const policies = new Array<string>(32).fill("").map((_, index) => "policy"+index)
        await Promise.all(policies.map((p) => nativeIAmPolicyGRPC.Create({
            name: p,
            namespace: "",
            actions: [],
            resources: []
        })))
        
        const returnedPolicies = [] as Array<string>
        await nativeIAmPolicyGRPC.List({ namespace: "", limit: 0, skip: 0 }).forEach((p) => {
            returnedPolicies.push(p.policy?.name as string)
        })

        expect(policies.sort()).toStrictEqual(returnedPolicies.sort())
    })

    describe("Returns list of policies with skip and limit", () => {        
        const testf = async (r: {skip: number, limit: number}) => {
            const policies = new Array<string>(10).fill("").map((_, index) => "policy"+index)
            await Promise.all(policies.map((p) => nativeIAmPolicyGRPC.Create({
                name: p,
                namespace: "",
                actions: [],
                resources: []
            })))

            const allPolicies = [] as Array<string>
            await nativeIAmPolicyGRPC.List({ namespace: "", limit: 999, skip: 0 }).forEach((p) => {
                allPolicies.push(p.policy?.name as string)
            })
            
            const returnedPolicies = [] as Array<string>
            await nativeIAmPolicyGRPC.List({ namespace: "", limit: r.limit, skip: r.skip }).forEach((p) => {
                returnedPolicies.push(p.policy?.name as string)
            })

            let policySlice = allPolicies.slice(
                Math.min(allPolicies.length, r.skip),
                Math.min(allPolicies.length, r.skip+(r.limit === 0 ? 9999 : r.limit))
            )
            expect(returnedPolicies).toStrictEqual(policySlice)
        }
        const ranges = [
            {
                skip: 1,
                limit: 0
            },
            {
                skip: 1,
                limit: 1
            },
            {
                skip: 999,
                limit: 0
            },
        ] as Array<{skip: number, limit: number}>
        for (const queryRange of ranges) {
            test(`skip: ${queryRange.skip}; limit: ${queryRange.limit}`, async () => await testf(queryRange))
        }
    })
})