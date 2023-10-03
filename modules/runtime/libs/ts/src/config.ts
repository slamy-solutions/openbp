export interface Config {
    urls: {
        manager: string,
    }
}

export const makeDefaultConfig: ()=>Config = () => {
    return {
        urls: {
            manager: process.env.RUNTIME_MANAGER_URL || "runtime_manager:80",
        }
    }
}