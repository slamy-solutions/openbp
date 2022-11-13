import { connect, close, managerClient } from './grpc'

export const managerGrpc = managerClient

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}