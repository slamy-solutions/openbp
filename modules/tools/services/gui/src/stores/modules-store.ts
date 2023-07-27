import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'

export interface ModulesState {
  iot: boolean
}

export const useModulesStore = defineStore('modules', {
  state: () => ({
    _informationLoaded: false,
    _iot: false,
  }),
  getters: {
    informationLoaded: (state) => state._informationLoaded,
    iot: (state) => state._iot,
  },
  actions: {
    updateModulesState(newState: ModulesState) {
        this._informationLoaded = true
        this._iot = newState.iot
    }
  }
})