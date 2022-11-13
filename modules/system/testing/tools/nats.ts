import { connect as connectToNats, NatsConnection } from 'nats'

let client: NatsConnection | null = null

export function getClient() {
    return client as NatsConnection
}

export async function connect() {
    client = await connectToNats()
}

export async function close() {
    await client?.drain()
}