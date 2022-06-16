import Redis from "ioredis";

export const client = new Redis({ lazyConnect: true })

export async function connect() {
    await client.connect()
}

export async function close() {
    await client.quit()
}
