import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
    login: string
    loggedIn: boolean
    accessToken: string
    refreshToken: string
}
const initialState = {
    login: "",
    loggedIn: false,
    accessToken: "",
    refreshToken: ""
} as AuthState

const LOCAL_STORAGE_SAVE_KEY = "AUTH_DATA"

const authSlice = createSlice({
    name: 'features/auth',
    initialState,
    reducers: {
        login (state, action: PayloadAction<{login: string, accessToken: string, refreshToken: string}>) {
            state.login = action.payload.login
            state.accessToken = action.payload.accessToken
            state.refreshToken = action.payload.refreshToken
            state.loggedIn = true
            localStorage.setItem(LOCAL_STORAGE_SAVE_KEY, JSON.stringify(state))
        },
        logout (state) {
            state.login = ""
            state.accessToken = ""
            state.refreshToken = ""
            state.loggedIn = false
            localStorage.removeItem(LOCAL_STORAGE_SAVE_KEY)
        }, 
        load (state) {
            const data = localStorage.getItem(LOCAL_STORAGE_SAVE_KEY)
            if (data != null) {
                const dataObj = JSON.parse(data) as AuthState
                state.login = dataObj.login
                state.accessToken = dataObj.accessToken
                state.refreshToken = dataObj.refreshToken
                state.loggedIn = true
            }
        }
    }
})

export const { login, logout, load } = authSlice.actions
export default authSlice