import { APIModuleBase } from "../../model"

export type OAuthProviderName = 'github' | 'gitlab' | 'microsoft' | 'google' | 'discord' | 'facebook' | 'apple' | 'twitter' | 'oidc' | 'oidc2' | 'oidc3' | 'instagram'

export interface OAuthProviderConfig {
    namespace: string
    enabled: boolean
    name: OAuthProviderName
    clientID: string
    clientSecret: string
    authURL: string
    tokenURL: string
    userApiURL: string
}

export interface GetProvidersConfigsRequest {
    namespace: string
}
export interface GetProvidersConfigsResponse {
    providers: OAuthProviderConfig[]
}

export interface UpdateProviderConfigRequest {
    namespace: string
    enabled: boolean
    name: OAuthProviderName
    clientID: string
    clientSecret: string
    authURL: string
    tokenURL: string
    userApiURL: string
}
export interface UpdateProviderConfigResponse {}

export class OAuthAPI extends APIModuleBase {
    // Get all providers configs
    async getProvidersConfigs(params: GetProvidersConfigsRequest): Promise<GetProvidersConfigsResponse> {
        const response = await OAuthAPI._axios.get<GetProvidersConfigsResponse>('/accessControl/config/oauth/provider/list', { params })
        return response.data
    } 

    // Update provider config
    async updateProviderConfig(params: UpdateProviderConfigRequest): Promise<UpdateProviderConfigResponse> {
        const response = await OAuthAPI._axios.patch<UpdateProviderConfigResponse>('/accessControl/config/oauth/provider', params)
        return response.data
    }
}