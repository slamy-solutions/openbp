import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'

interface BootstrapState {
    fullyBootstrapped: boolean | null
    rootUserInitialized: boolean
}
const initialState = {
    fullyBootstrapped: null,
    rootUserInitialized: false
} as BootstrapState

const bootstrapSlice = createSlice({
    name: 'features/bootstrap',
    initialState,
    reducers: {
        setState (state, action: PayloadAction<BootstrapState>) {
            state.fullyBootstrapped = action.payload.fullyBootstrapped
            state.rootUserInitialized = action.payload.rootUserInitialized
        }
    }
})

export const { setState } = bootstrapSlice.actions
export default bootstrapSlice