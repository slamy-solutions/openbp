import { Client } from '../../../../system/libs/ts/grpc'

import { IAMAuthServiceClientImpl } from './proto/auth'
import { IAMConfigServiceClientImpl } from './proto/configuration'
import { IAMIdentityServiceClientImpl } from './proto/identity'
import { IAMPolicyServiceClientImpl } from './proto/policy'

const grpcAuthClient = new Client("native_iam_auth:80")
export const authClient = new IAMAuthServiceClientImpl(grpcAuthClient)

const grpcConfigClient = new Client("native_iam_config:80")
export const configClient = new IAMConfigServiceClientImpl(grpcConfigClient)

const grpcIdentityClient = new Client("native_iam_identity:80")
export const identityClient = new IAMIdentityServiceClientImpl(grpcIdentityClient)

const grpcPolicyClient = new Client("native_iam_policy:80")
export const policyClient = new IAMPolicyServiceClientImpl(grpcPolicyClient)

export async function connect() {
    // await grpcAuthClient.connect()
    await grpcConfigClient.connect()
    // await grpcIdentityClient.connect()
    await grpcPolicyClient.connect()
}

export async function close() {
    // grpcAuthClient.close()
    grpcConfigClient.close()
    // grpcIdentityClient.close()
    grpcPolicyClient.close()
}