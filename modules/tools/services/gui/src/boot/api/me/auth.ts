import { APIModuleBase } from "../model"

export type OAuthProviderName = 'github' | 'gitlab' | 'microsoft' | 'google' | 'discord' | 'facebook' | 'apple' | 'twitter' | 'oidc' | 'oidc2' | 'oidc3' | 'instagram'

export interface AvailableOAuthProvider {
    name: string
    clientId: string
    authUrl: string
}
export interface ConfiguredOAuthProvider {
    name: OAuthProviderName
    userId: string
}

export interface GetInfoResponse {
    password: {
        enabled: boolean
    }
    oauth: {
        availableProviders: Array<AvailableOAuthProvider>
        configuredProviders: Array<ConfiguredOAuthProvider>
    }
}

export interface CreateOrUpdatePasswordRequest {
    password: string
}

export interface FinalizeOAuthRegistrationRequest {
    provider: OAuthProviderName
    code: string
}

export interface FinalizeOAuthRegistrationResponse {
    status: 'OK' | 'ALREADY_REGISTERED' | 'PROVIDER_DISABLED' | 'ERROR_WHILE_RETRIEVING_AUTH_TOKEN' | 'ERROR_WHILE_FETCHING_USER_DETAILS'
}

export interface DeleteOAuthRegistrationRequest {
    provider: OAuthProviderName
}

export class AuthAPI extends APIModuleBase {
    async getInfo(): Promise<GetInfoResponse> {
        const response = await AuthAPI._axios.get<GetInfoResponse>('/me/auth')
        return response.data
    }

    async createOrUpdatePassword(params: CreateOrUpdatePasswordRequest): Promise<void> {
        await AuthAPI._axios.put('/me/auth/password', params)
    }

    async deletePassword(): Promise<void> {
        await AuthAPI._axios.delete('/me/auth/password')
    }

    async finalizeOAuthRegistration(params: FinalizeOAuthRegistrationRequest): Promise<FinalizeOAuthRegistrationResponse> {
        const repsonse = await AuthAPI._axios.post<FinalizeOAuthRegistrationResponse>('/me/auth/oauth/provider', params)
        return repsonse.data
    }

    async deleteOAuthRegistration(params: DeleteOAuthRegistrationRequest): Promise<void> {
        await AuthAPI._axios.delete('/me/auth/oauth/provider', { data: params })
    }
}