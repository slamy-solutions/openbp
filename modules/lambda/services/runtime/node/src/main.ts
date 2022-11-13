import { connect as connectToAMQP } from 'amqplib'
import {} from './pod'

const EXCHANGE_NAME = ""
const QUEUE_NAME = ""

async function main() {
    const amqpConnection = await connectToAMQP("")
    const amqpChanel = await amqpConnection.createChannel()
    await amqpChanel.assertExchange(EXCHANGE_NAME, "direct", { durable: true, autoDelete: false })
    await amqpChanel.assertQueue(QUEUE_NAME, { durable: true })
    await amqpChanel.consume(QUEUE_NAME, () => {}, { noAck: false })
    console.log("Start consuming messages")
}

await main()
