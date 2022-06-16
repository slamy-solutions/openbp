
import { MongoClient } from 'mongodb'

export const client = new MongoClient("");

export async function connect() {
    await client.connect();
    console.log(`Successfully connected to MongoDB`);
}

export async function close() {
    await client.close()
    console.log(`Successfully closed connection to MongoDB`);
}