import { configureStore } from '@reduxjs/toolkit'

import { apiSlice } from './context/api'

import AuthSlice from './features/auth/auth.slice'
import AuthApi from './features/auth/api.slice'

import BootstrapSlice from './features/bootstrap/bootstrap.slice'
import BootstrapApi from './features/bootstrap/api.slice'

export const store = configureStore({
    reducer: {
        [apiSlice.name]: apiSlice.reducer,

        [AuthSlice.name]: AuthSlice.reducer,
        [AuthApi.reducerPath]: AuthApi.reducer,

        [BootstrapSlice.name]: BootstrapSlice.reducer,
        [BootstrapApi.reducerPath]: BootstrapApi.reducer
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware()
        .concat(AuthApi.middleware)
        .concat(BootstrapApi.middleware),
    devTools: true
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch