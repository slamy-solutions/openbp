import { Client } from '../../../../system/libs/ts/grpc'

import { IAMOAuthServiceClientImpl } from './proto/oauth'
import { IAMConfigServiceClientImpl } from './proto/configuration'
import { IAMIdentityServiceClientImpl } from './proto/identity'
import { IAMPolicyServiceClientImpl } from './proto/policy'
import { IAMTokenServiceClientImpl } from './proto/token'

import { IAMAuthenticationPasswordServiceClientImpl } from './proto/authentication/password'

const grpcOAuthClient = new Client("native_iam_oauth:80")
export const oauthClient = new IAMOAuthServiceClientImpl(grpcOAuthClient)

const grpcConfigClient = new Client("native_iam_config:80")
export const configClient = new IAMConfigServiceClientImpl(grpcConfigClient)

const grpcIdentityClient = new Client("native_iam_identity:80")
export const identityClient = new IAMIdentityServiceClientImpl(grpcIdentityClient)

const grpcPolicyClient = new Client("native_iam_policy:80")
export const policyClient = new IAMPolicyServiceClientImpl(grpcPolicyClient)

const grpcTokenClient = new Client("native_iam_token:80")
export const tokenClient = new IAMTokenServiceClientImpl(grpcTokenClient)

// Authentication
const grpcAuthenticationPasswordClient = new Client("native_iam_authentication_password:80")
export const authenticationPasswordClient = new IAMAuthenticationPasswordServiceClientImpl(grpcAuthenticationPasswordClient)

export async function connect() {
    await grpcOAuthClient.connect()
    await grpcConfigClient.connect()
    await grpcIdentityClient.connect()
    await grpcPolicyClient.connect()
    await grpcTokenClient.connect()

    await grpcAuthenticationPasswordClient.connect()
}

export async function close() {
    grpcOAuthClient.close()
    grpcConfigClient.close()
    grpcIdentityClient.close()
    grpcPolicyClient.close()
    grpcTokenClient.close()

    grpcAuthenticationPasswordClient.close()
}