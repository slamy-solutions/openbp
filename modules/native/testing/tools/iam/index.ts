import { connect, close, authClient, configClient, identityClient, policyClient } from './grpc'

export const authGRPC = authClient
export const clientGRPC = configClient
export const identityGRPC = identityClient
export const policyGRPC = policyClient

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}