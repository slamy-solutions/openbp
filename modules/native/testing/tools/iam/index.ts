import { connect, close, configClient } from './grpc'

export const clientGRPC = configClient

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}