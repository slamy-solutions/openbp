import { client as mongoClient, connect as connectToMongo, close as closeMongo } from '../../../../../system/testing/tools/mongo'
import { client as cacheClient, connect as connectToCache, close as closeCache } from '../../../../../system/testing/tools/cache'
import { configClient as iamConfigGRPC, connect as connectToNativeIAMConfig, close as closeNativeIAMConfig } from '../../../tools/iam/grpc'
import { client as nativeKeyValueStorageClient, connect as connectToNativeKeyValueStorage, close as closeNativeKeyValueStorage } from '../../../tools/keyvaluestorage/grpc'

const GLOBAL_DB_NAME = `openbp_global`

beforeAll(async ()=>{
    await connectToMongo()
    await connectToCache()
    await connectToNativeIAMConfig()
    await connectToNativeKeyValueStorage()
})

afterEach(async ()=>{
    await mongoClient.db(GLOBAL_DB_NAME).collection('native_keyvaluestorage').deleteMany({})
    await cacheClient.flushall()
})

afterAll(async ()=>{
    await closeMongo()
    await closeCache()
    await closeNativeIAMConfig()
    await closeNativeKeyValueStorage()
})

/**
 * @group native/iam/config/set/whitebox
 * @group whitebox
 */
describe("Whitebox", () => {
    test("Creates entry in native_keyvaluestorage", async () => {
        await iamConfigGRPC.Set({
            configuration: {
                accessTokenTTL: 2048,
                refreshTokenTTL: 1000*60*24,
                passwordAuth: {
                    enabled: true,
                    allowRegistration: true
                },
                googleOAuth2: {
                    allowRegistration: false,
                    clientId: "google123",
                    clientSecret: "",
                    enabled: false
                },
                facebookOAuth2: {
                    allowRegistration: true,
                    clientId: "facebook123",
                    clientSecret: "",
                    enabled: false
                },
                githubOAuth2: {
                    allowRegistration: false,
                    clientId: "github123",
                    clientSecret: "",
                    enabled: false
                },
                gitlabOAuth2: {
                    allowRegistration: false,
                    clientId: "gitlab123",
                    clientSecret: "",
                    enabled: false
                }
            }
        })

        const response = await nativeKeyValueStorageClient.Get({ namespace: "", key: "native_iam_config_global", useCache: false })
        expect(response.value).not.toBeNull()
    })
})

/**
 * @group native/iam/config/set/blackbox
 * @group blackbox
 */
 describe("Blackbox", () => {
    test("Overrides value on serveral sets", async () => {
        const ttl1 = 123
        const ttl2 = 234

        await iamConfigGRPC.Set({
            configuration: {
                accessTokenTTL: 2048,
                refreshTokenTTL: ttl1,
                passwordAuth: {
                    enabled: true,
                    allowRegistration: true
                },
                googleOAuth2: {
                    allowRegistration: false,
                    clientId: "google123",
                    clientSecret: "",
                    enabled: false
                },
                facebookOAuth2: {
                    allowRegistration: true,
                    clientId: "facebook123",
                    clientSecret: "",
                    enabled: false
                },
                githubOAuth2: {
                    allowRegistration: false,
                    clientId: "github123",
                    clientSecret: "",
                    enabled: false
                },
                gitlabOAuth2: {
                    allowRegistration: false,
                    clientId: "gitlab123",
                    clientSecret: "",
                    enabled: false
                }
            }
        })


        await iamConfigGRPC.Set({
            configuration: {
                accessTokenTTL: 2048,
                refreshTokenTTL: ttl2,
                passwordAuth: {
                    enabled: true,
                    allowRegistration: true
                },
                googleOAuth2: {
                    allowRegistration: false,
                    clientId: "google123",
                    clientSecret: "",
                    enabled: false
                },
                facebookOAuth2: {
                    allowRegistration: true,
                    clientId: "facebook123",
                    clientSecret: "",
                    enabled: false
                },
                githubOAuth2: {
                    allowRegistration: false,
                    clientId: "github123",
                    clientSecret: "",
                    enabled: false
                },
                gitlabOAuth2: {
                    allowRegistration: false,
                    clientId: "gitlab123",
                    clientSecret: "",
                    enabled: false
                }
            }
        })

        const response = await iamConfigGRPC.Get({useCache: false})
        expect(response.configuration?.refreshTokenTTL).toBe(ttl2)
    })
})