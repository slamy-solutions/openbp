export interface Config {
    urls: {
        core: string,
    }
}

export const makeDefaultConfig: ()=>Config = () => {
    return {
        urls: {
            core: process.env.ERP_CORE_URL || "erp_core:80",
        }
    }
}