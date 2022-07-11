import { Connection, connect as amqpConnect } from 'amqplib'

let connection: Connection | null = null

export function getClient() {
    if (connection === null) {
        throw new Error("Not connected to amqp")
    }
    return connection
}

export async function connect() {
    connection = await amqpConnect(process.env.SYSTEM_RABBITMQ_URL || "amqp://system_rabbitmq:5672")
}

export async function close() {
    await connection?.close()
    connection = null
}