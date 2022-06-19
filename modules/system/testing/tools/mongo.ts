
import { Collection, MongoClient } from 'mongodb'



export const client = new MongoClient(process.env.SYSTEM_DB_URL || "mongodb://root:example@system_db/admin");

export async function connect() {
    await client.connect();
    console.log(`Successfully connected to MongoDB`);
}

export async function close() {
    await client.close()
    console.log(`Successfully closed connection to MongoDB`);
}