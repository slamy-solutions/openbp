import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'

export interface NamespaceState {
    currentNamespace: string
    visibleNamespaces: string[]
}

export const useNamespaceStore = defineStore('namespace', {
  state: () => ({
    _informationLoaded: false,
    _currentNamespace: "",
    _visibleNamespaces: [] as Array<string>,
  }),
  getters: {
    loaded: (state) => state._informationLoaded,
    currentNamespace: (state) => state._currentNamespace,
    visibleNamespaces: (state) => state._visibleNamespaces,
  },
  actions: {
    tryGetFromLocalStorage() {

    },
    updateBootstrapState(newState: NamespaceState) {
        this._informationLoaded = true
        this._currentNamespace = newState.currentNamespace
        this._visibleNamespaces = newState.visibleNamespaces
    }
  }
})