import { defineStore } from 'pinia'

export interface ModulesState {
  iot: boolean
  crm: boolean
  erp: boolean
}

export const useModulesStore = defineStore('modules', {
  state: () => ({
    _informationLoaded: false,
    _iot: false,
    _crm: false,
    _erp: false,
  }),
  getters: {
    informationLoaded: (state) => state._informationLoaded,
    iot: (state) => state._iot,
    crm: (state) => state._crm,
    erp: (state) => state._erp,
  },
  actions: {
    updateModulesState(newState: ModulesState) {
        this._informationLoaded = true
        this._iot = newState.iot
        this._crm = newState.crm
        this._erp = newState.erp
    }
  }
})