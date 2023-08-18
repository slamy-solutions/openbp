export interface Config {
    urls: {
        namespace: string,
        keyvaluestorage: string
        iam: string
    }
}

export const makeDefaultConfig: ()=>Config = () => {
    return {
        urls: {
            namespace: process.env.NATIVE_NAMESPACE_URL || "native_namespace:80",
            keyvaluestorage: process.env.NATIVE_KEYVALUESTORAGE_URL || "native_keyvaluestorage:80",
            iam: process.env.NATIVE_IAM_URL || "native_iam:80",
        }
    }
}