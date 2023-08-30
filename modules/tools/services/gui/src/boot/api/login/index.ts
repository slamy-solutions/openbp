import { APIModuleBase } from "../model"

export interface CreateTokensWithPasswordRequest {
    namespace: string
    login: string
    password: string
}
export interface CreateTokensWithPasswordResponse {
    accessToken: string
    refreshToken: string
}

export interface CreateTokenWithOAuthRequest {
    namespace: string
    provider: string
    code: string
}
export interface CreateTokenWithOAuthResponse {
    accessToken: string
    refreshToken: string
}

export interface AvailableOAuthProvider {
    name: string
    clientId: string
    authUrl: string
}

export interface GetAvailableOAuthProvidersRequest {
    namespace: string
}
export interface GetAvailableOAuthProvidersResponse {
    providers: Array<AvailableOAuthProvider>
}

export interface RefreshTokenRequest {
    refreshToken: string
}
export interface RefreshTokenResponse {
    accessToken: string
}

export interface ValidateTokenRequest {
    token: string
}
export interface ValidateTokenResponse {
    valid: boolean
}

export class LoginAPI extends APIModuleBase {

    // Try to create access and refresh token with credentials
    async createTokensWithPassword(credentials: CreateTokensWithPasswordRequest): Promise<CreateTokensWithPasswordResponse> {
        const response = await LoginAPI._axios.post<CreateTokensWithPasswordResponse>('/auth/login/password', credentials)
        return response.data
    }
    
    // Try to create access and refresh token with OAuth
    async createTokenWithOAuth(credentials: CreateTokenWithOAuthRequest): Promise<CreateTokenWithOAuthResponse> {
        const response = await LoginAPI._axios.post<CreateTokenWithOAuthResponse>('/auth/login/oauth', credentials)
        return response.data
    }

    // Get available OAuth providers for namespace
    async getAvailableOAuthProviders(params: GetAvailableOAuthProvidersRequest): Promise<GetAvailableOAuthProvidersResponse> {
        const response = await LoginAPI._axios.get<GetAvailableOAuthProvidersResponse>('/auth/login/oauth/providers', { params })
        return response.data
    }

    // Try to refresh token and get new access token 
    async refreshToken(params: RefreshTokenRequest): Promise<string> {
        const response = await LoginAPI._axios.post<RefreshTokenResponse>('/auth/token/refresh', params)
        return response.data.accessToken
    }

    // Check it token is ok
    async validateToken(params: ValidateTokenRequest): Promise<boolean> {
        const response = await LoginAPI._axios.post<ValidateTokenResponse>('/auth/token/validate', params)
        return response.data.valid
    }
}