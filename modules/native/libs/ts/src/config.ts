export interface Config {
    urls: {
        namespace: string,
        keyvaluestorage: string
        iam: {
            authentication: {
                password: string
            }
            policy: string
            token: string
            oauth: string
            identity: string
        },
        actor: {
            user: string
        }
    }
}

export const makeDefaultConfig: ()=>Config = () => {
    return {
        urls: {
            namespace: process.env.NATIVE_NAMESPACE_URL || "native_namespace:80",
            keyvaluestorage: process.env.NATIVE_KEYVALUESTORAGE_URL || "native_keyvaluestorage:80",
            iam: {
                authentication: {
                    password: process.env.NATIVE_IAM_AUTHENTICATION_PASSWORD_URL || "native_iam_authentication_password:80"
                },
                policy: process.env.NATIVE_IAM_POLICY_URL || "native_iam_policy:80",
                token: process.env.NATIVE_IAM_TOKEN_URL || "native_iam_token:80",
                oauth: process.env.NATIVE_IAM_OAUTH_URL || "native_iam_oauth:80",
                identity: process.env.NATIVE_IAM_IDENTITY_URL || "native_iam_identity:80",
            },
            actor: {
                user: process.env.NATIVE_ACTOR_USER_URL || "native_actor_user:80"
            }
        }
    }
}