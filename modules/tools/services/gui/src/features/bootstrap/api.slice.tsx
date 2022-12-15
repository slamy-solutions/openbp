import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { baseQueryWithoutAuth } from '../../context/api'

export interface GetStatusRequest {}
export interface GetStatusResponse {
    fullyBootstrapped: boolean
}

export interface InitRootUserRequest {
    login: string
    password: string
}
export interface InitRootUserResponse {}

const bootstrapApi = createApi({
    reducerPath: 'features/bootstrap/api',
    baseQuery: baseQueryWithoutAuth,
    endpoints: build => ({
        getStatus: build.query<GetStatusResponse, GetStatusRequest>({
            query: () => ({
                url: 'bootstrap/status',
                method: 'GET'
            }),
        }),
        initRootUser: build.mutation<InitRootUserResponse, InitRootUserRequest>({
            query: (body) => ({
                url: 'bootstrap/rootUser',
                method: 'POST',
                body
            })
        })
    })
})

export const { useGetStatusQuery, useLazyGetStatusQuery, useInitRootUserMutation } = bootstrapApi
export default bootstrapApi