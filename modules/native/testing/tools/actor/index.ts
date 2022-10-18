import { connect, close, userClient } from './grpc'

export const userGRPC = userClient

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}