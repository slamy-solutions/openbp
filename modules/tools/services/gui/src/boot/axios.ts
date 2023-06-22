import { boot } from 'quasar/wrappers';
import axios, { AxiosInstance } from 'axios';

import api, { API } from './api'
import { APIModuleBase } from './api/model'
import { env } from 'process';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: API
  }
}

export default boot(({ app }) => {
  app.config.globalProperties.$axios = axios.create({});
  app.config.globalProperties.$api = api;

  if (process.env.DEV) {
    APIModuleBase.setBaseURL("http://127.0.0.1:80/api")
  } else {
    APIModuleBase.setBaseURL("/api")
  }
});

export { api };
