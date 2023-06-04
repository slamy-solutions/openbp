import { APIModuleBase } from "../model"

export interface CreateTokensWithPasswordRequest {
    login: string
    password: string
}
export interface CreateTokensWithPasswordResponse {
    accessToken: string
    refreshToken: string
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