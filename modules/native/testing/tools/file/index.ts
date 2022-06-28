import { connect, close, client } from './grpc'
import { TestFile } from './testfile'

export const grpc = client

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}