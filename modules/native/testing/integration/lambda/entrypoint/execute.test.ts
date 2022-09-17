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
import { Lambda, AMQPLambdaTaskRequest, AMQPLambdaTaskResponse, ExecuteLambdaResponse } from '../../../tools/lambda/proto/lambda'


const TEST_NAMESPACE_NAME = "lambdanamespace"
const BUNDLE_DB_NAME = `openbp_global`
const INFO_DB_NAME = `openbp_namespace_${TEST_NAMESPACE_NAME}`

const RESPONSE_EXCHANGE = "native_lambda_entrypoint_output"

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
 * @group native/lambda/entrypoint/execute/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates task on AMQP queue and returns result", async () => {
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

        const executeTask = lambdaEntrypointClient.Execute({
            timeout: 3000,
            namespace: TEST_NAMESPACE_NAME,
            lambda: uuid,
            data: Buffer.from(JSON.stringify({ something: "ok" }))
        })

        let receivedMessage: ConsumeMessage | null = null
        const c = await chanel.consume(queue.queue, (msg) => {receivedMessage = msg}, { noAck: true })
        let waits = 100
        while (receivedMessage === null && waits > 0) {
            await new Promise((resolve) => setTimeout(resolve, 10))
            if (receivedMessage !== null) break
            waits -= 1
        }
        await chanel.cancel(c.consumerTag)
        expect(waits).not.toBe(0)
        expect(receivedMessage).not.toBeNull()
        const m = receivedMessage as unknown as ConsumeMessage
        const replyTo = m.properties.replyTo
        const correlationId = m.properties.correlationId

        const lambdaresponse = AMQPLambdaTaskResponse.fromPartial({
            statusCode: 0,
            message: "OK",
            data: Buffer.from(JSON.stringify({"resp": "successtest"}))
        })
        const lambdaResponseBytes = AMQPLambdaTaskResponse.encode(lambdaresponse).finish()

        chanel.publish(RESPONSE_EXCHANGE, replyTo, Buffer.from(lambdaResponseBytes), { correlationId })

        const executionResponse = await executeTask
        const responseData = JSON.parse(executionResponse.result.toString()) as { "resp": string }
        expect(responseData.resp).toBe("successtest")
    })
})

/**
 * @group native/lambda/entrypoint/execute/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Fails with NOT_FOUND error if lambda doesnt exist", async () => {
        try {
            await lambdaEntrypointClient.Execute({
                timeout: 3000,
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