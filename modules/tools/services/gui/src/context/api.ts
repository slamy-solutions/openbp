import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { BaseQueryFn, FetchArgs, FetchBaseQueryError, FetchBaseQueryMeta, fetchBaseQuery, createApi } from '@reduxjs/toolkit/query/react'
import { RootState } from '../store'


import { logout } from '../features/auth/auth.slice'

const LOCAL_STORAGE_KEY = 'CONTEXT_API'

interface APIState {
    baseURL: string
}


const savedData = localStorage.getItem(LOCAL_STORAGE_KEY)
const initialState = (savedData) ? JSON.parse(savedData) as APIState : { baseURL: '/api' }

export const apiSlice = createSlice({
    name: 'context/api',
    initialState,
    reducers: {
        setBaseUrl(state, action: PayloadAction<{url: string}>) {
            state.baseURL = action.payload.url
            localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(state))
        }
    }
})

export const baseQueryWithoutAuth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (args, api, extraOptions) => {
    const state = api.getState() as RootState
    const baseUrl = state['context/api'].baseURL;
    const rawBaseQuery = fetchBaseQuery({ baseUrl })
    return rawBaseQuery(args, api, extraOptions);
}

export const baseQueryWithAuth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (args, api, extraOptions) => {
    const state = api.getState() as RootState
    
    const baseUrl = state['context/api'].baseURL;
    const rawBaseQuery = fetchBaseQuery({
        baseUrl,
        credentials: 'include',
        prepareHeaders: (headers) => {
            const loggedIn = state['features/auth'].loggedIn

            if (loggedIn) {
                const token = state['features/auth'].accessToken
                headers.set('authorization', `Bearer ${token}`)
            }

            return headers
        }
    })
    return rawBaseQuery(args, api, extraOptions);
};

export const baseQueryWithReauth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (args, api, extraOptions) => {
    const state = api.getState() as RootState
    let result = await baseQueryWithAuth(args, api, extraOptions)

    if (result?.error?.status === 401 && state['features/auth'].loggedIn) {
        console.log("[context/api]: Received 401 error on request with access token. Trying to refresh token.")

        //TODO: 
        //const refreshResult = await baseQuery

        //api.dispatch(AuthApi.endpoints.refreshToken())
    
        console.log("[context/api]: Failed to refresh token. Logging out.")   
        api.dispatch(logout()) 
    }

    return result
}