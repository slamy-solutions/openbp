import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'

export interface BootstrapState {
    vaultSealed: boolean
    rootUserCreated: boolean
    rootUserCreationBlocked: boolean
}

export const useBootstrapStore = defineStore('bootstrap', {
  state: () => ({
    _informationLoaded: false,
    _vaultSealed: false,
    _rootUserCreated: false,
    _rootUserCreationBlocked: false
  }),
  getters: {
    bootstrapped: (state) => {
        return state._informationLoaded && !state._vaultSealed && state._rootUserCreated
    },
    vaultSealed: (state) => state._vaultSealed,
    rootUserCreated: (state) => state._rootUserCreated,
    rootUserCreationBlocked: (state) => state._rootUserCreationBlocked,
  },
  actions: {
    updateBootstrapState(newState: BootstrapState) {
        this._informationLoaded = true
        this._vaultSealed = newState.vaultSealed
        this._rootUserCreated = newState.rootUserCreated
        this._rootUserCreationBlocked = newState.rootUserCreationBlocked
    }
  }
})