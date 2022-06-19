import Redis, {  } from "ioredis";


export const client = new Redis(process.env.SYSTEM_CACHE_URL || "redis://system_cache")

export async function connect() {
    if (client.status === "ready" || client.status === 'connecting' || client.status === "reconnecting") {
        console.log("Already connected to redis")
    } else {
        await client.connect()
        console.log(`Successfully connected to cache`);
    }
}

export async function close() {
    await client.quit()
    console.log(`Successfully closed connection to cache`);
}
