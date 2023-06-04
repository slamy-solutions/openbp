import axios, { AxiosError, AxiosInstance } from 'axios';
import { useLoginStore } from 'src/stores/login-store';
import { useRouter } from 'vue-router';

declare module 'axios' {
    interface InternalAxiosRequestConfig<D> {
      _retry?: boolean;
    }
}

export class APIModuleBase{
    protected static _axios = APIModuleBase.init()

    private static init(): AxiosInstance {
        const _axios = axios.create({})

        // Interceptor for injecting authorization headers
        _axios.interceptors.request.use((config)=>{
            const loginStore = useLoginStore()
            if (loginStore.isLoggedIn && config.headers.Authorization === undefined) {
                config.headers.Authorization = `Bearer ${loginStore.accessToken}`
            }
            return config
        })

        // Interceptor that handles unauthenticated responses and retries them with refreshed auth token.
        _axios.interceptors.response.use(
            (response) => response,
            async (error: AxiosError) => {
                const loginStore = useLoginStore()
                const originalRequest = error.config

                if (error.response?.status === 401 && originalRequest !== undefined && !originalRequest._retry && loginStore.refreshToken != "") {
                    originalRequest._retry = true
                    try {
                        
                        const refreshResponse = await axios.post<{accessToken: string}>(`${_axios.defaults.baseURL}/auth/token/refresh`, {refreshToken: loginStore.refreshToken})
                        loginStore.updateAccessToken(refreshResponse.data.accessToken)
                        originalRequest.headers.Authorization = `Bearer ${refreshResponse.data.accessToken}`;
                        return _axios(originalRequest)
                    } catch {
                        console.log("Logged out")
                        loginStore.logout()
                        //await useRouter().push({name: "logout"})
                    }
                }

                return Promise.reject(error)
            }
        )

        return _axios
    }

    public static setBaseURL(baseURL: string) {
        this._axios.defaults.baseURL = baseURL
    }
}