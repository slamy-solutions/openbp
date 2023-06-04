import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'

const _LOCAL_STORAGE_KEY = 'LOGIN_DATA'

interface LocalStorageLoginData {
    username: string
    accessToken: string
    refreshToken: string
}

export const useLoginStore = defineStore('login', {
  state: () => ({
    _isLoggedIn: false,
    _username: '',
    _accessToken: '',
    _refreshToken: ''
  }),
  getters: {
    isLoggedIn: (state) => state._isLoggedIn,
    username: (state) => state._username,
    accessToken: (state) => state._accessToken,
    refreshToken: (state) => state._refreshToken
  },
  actions: {
    tryLoadLoginFromStorage () {
      if (LocalStorage.has(_LOCAL_STORAGE_KEY)) {
        this._isLoggedIn = true
        const data = LocalStorage.getItem(_LOCAL_STORAGE_KEY) as LocalStorageLoginData
        this._username = data.username
        this._accessToken = data.accessToken
        this._refreshToken = data.refreshToken
      }
    },

    updateAccessToken(newAccessToken: string) {
        this._accessToken = newAccessToken
    },

    login (username: string, accessToken: string, refreshToken: string) {
        this._isLoggedIn = true
        this._username = username
        this._accessToken = accessToken
        this._refreshToken = refreshToken
        LocalStorage.set(_LOCAL_STORAGE_KEY, {
            username: username,
            accessToken: accessToken,
            refreshToken: refreshToken
        } as LocalStorageLoginData)
    },

    logout () {
      this._isLoggedIn = false
      this._username = ''
      this._accessToken = ''
      this._refreshToken = ''
      LocalStorage.remove(_LOCAL_STORAGE_KEY)
    }
  }
})