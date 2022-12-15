import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { baseQueryWithoutAuth } from '../../context/api'

export interface LoginWithPasswordRequest {
    login: string
    password: string
}
export interface LoginWithPasswordResponse {
    accessToken: string
    refreshToken: string
}

export interface RefreshTokenRequest {
    refreshToken: string
}
export interface RefreshTokenResponse {
    accessToken: string
}

const authApi = createApi({
    reducerPath: 'features/auth/api',
    baseQuery: baseQueryWithoutAuth,
    endpoints: build => ({
        loginWithPassword: build.mutation<LoginWithPasswordResponse, LoginWithPasswordRequest>({
            query: (body) => ({
                url: 'auth/login/password',
                method: 'POST',
                body
            }),
        }),
        refreshToken: build.mutation<RefreshTokenResponse, RefreshTokenRequest>({
            query: (body) => ({
                url: 'auth/token/refresh',
                method: 'POST',
                body
            })
        })
    })
})

export const { useLoginWithPasswordMutation, useRefreshTokenMutation } = authApi
export default authApi