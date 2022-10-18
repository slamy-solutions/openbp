import { connect, close, oauthClient, configClient, identityClient, policyClient, authenticationPasswordClient } from './grpc'

export const authenticationPasswordGRPC = authenticationPasswordClient
export const oauthGRPC = oauthClient
export const clientGRPC = configClient
export const identityGRPC = identityClient
export const policyGRPC = policyClient

export async function setup() {
    await connect()
}

export async function teardown() {
    await close()
}