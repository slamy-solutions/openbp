import { randomBytes } from 'crypto'
import { Status } from '@grpc/grpc-js/build/src/constants'
import { Binary, GridFSBucket, ObjectId } from 'mongodb'
import { observable, Observable } from 'rxjs'
import { ConsumeMessage } from 'amqplib'

import { RequestError as GRPCRequestError } from '../../../../../system/libs/ts/grpc'
import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { getClient as getAMQPClient, connect as connectAMQP, close as closeAMQP } from '../../../../../system/testing/tools/rabbitmq'
import { client as namespaceGrpc, connect as connectToNativeNamespace, close as closeNativeNamespace } from '../../../tools/namespace/grpc'
import { managerClient as lambdaManagerGrpc, entrypointClient as lambdaEntrypointClient, connect as connectToNativeLambda, close as closeNativeLambda } from '../../../tools/lambda/grpc'
import { Lambda, AMQPLambdaTaskRequest } from '../../../tools/lambda/proto/lambda'


const TEST_NAMESPACE_NAME = "lambdanamespace"
const BUNDLE_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}global`
const INFO_DB_NAME = `${process.env.SYSTEM_DB_PREFIX || "openerp_"}namespace_${TEST_NAMESPACE_NAME}`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectAMQP()
    await connectToNativeNamespace()
    await connectToNativeLambda()
})

beforeEach(async () => {
    await namespaceGrpc.Ensure({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
    try {
        await mongoClient.db(BUNDLE_DB_NAME).collection("native_lambda_manager_bundle").deleteMany({})
    } catch {}
})

afterEach(async ()=>{
    await namespaceGrpc.Delete({ name: TEST_NAMESPACE_NAME })
    await cacheClient.flushall()
    try {
        await mongoClient.db(BUNDLE_DB_NAME).collection("native_lambda_manager_bundle").deleteMany({})
    } catch {}
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeAMQP()
    await closeNativeNamespace()
    await closeNativeLambda()
})

/**
 * @group native/lambda/entrypoint/call/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates task on AMQP queue", async () => {
        const uuid = "customlambda" + randomBytes(20).toString("hex")
        const bundle = randomBytes(32)
        const data = randomBytes(32)
        const runtime = randomBytes(16).toString("hex")

        await lambdaManagerGrpc.Create({
            namespace: TEST_NAMESPACE_NAME,
            uuid,
            bundle,
            data,
            ensureExactlyOneDelivery: true,
            runtime
        })

        const chanel = await getAMQPClient().createChannel()
        const exchange = await chanel.assertExchange("native_lambda_entrypoint_input", "direct")
        const queue = await chanel.assertQueue(`native_lambda_entrypoint_input_${runtime}`, { autoDelete: true, durable: false })
        await chanel.bindQueue(queue.queue, exchange.exchange, runtime)

        await lambdaEntrypointClient.Call({
            namespace: TEST_NAMESPACE_NAME,
            lambda: uuid,
            data: Buffer.from(JSON.stringify({ something: "ok" }))
        })

        let receivedMessage: ConsumeMessage | null = null
        const c = await chanel.consume(queue.queue, (msg) => {receivedMessage = msg}, { noAck: true })
        let waits = 10
        while (receivedMessage === null && waits > 0) {
            await new Promise((resolve) => setTimeout(resolve, 100))
            if (receivedMessage !== null) break
            waits -= 1
        }
        await chanel.cancel(c.consumerTag)
        expect(waits).not.toBe(0)
        expect(receivedMessage).not.toBeNull()

        const m = receivedMessage as unknown as ConsumeMessage
        const task = AMQPLambdaTaskRequest.decode(m.content)
        expect(task.lambda?.uuid).toBe(uuid)
        const amqpTaskData =  JSON.parse(task.data.toString('utf-8')) as { something: string }
        expect(amqpTaskData.something).toBe("ok")
    })
})

/**
 * @group native/lambda/entrypoint/call/blackbox
 * @group blackbox
 */
describe("Blackbox", () => {
    test("Fails with NOT_FOUND error if lambda doesnt exist", async () => {
        try {
            await lambdaEntrypointClient.Call({
                namespace: TEST_NAMESPACE_NAME,
                lambda: 'undefinedlambda',
                data: Buffer.from(JSON.stringify({ something: "ok" }))
            })
            fail()
        } catch (e) {
            if ((e as GRPCRequestError)?.code !== Status.NOT_FOUND) fail()
        }
    })
})