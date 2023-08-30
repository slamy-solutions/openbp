<template>
    <div class="bg-primary window-height window-width row justify-center items-center">
        <div class="column">
        <div class="row items-center justify-center q-mb-lg">
            <q-img src="logo.svg" width="350px"></q-img>
        </div>
        <div class="row">
            <q-card flat square class="q-pa-lg bg-transparent">
            <q-card-section>
                <q-form class="q-gutter-md">
                <q-input square outlined clearable v-model="namespace" type="text" :label="$t('modules.login.namespaceInput')" :disable="namespaceSelectionBlocked"/>
                <q-input square outlined clearable v-model="username"  type="text" :label="$t('modules.login.usernameInput')" />
                <q-input square outlined clearable v-model="password"  type="password" :label="$t('modules.login.passwordInput')" />
                </q-form>
            </q-card-section>
            <q-card-actions class="q-px-md">
                <q-btn unelevated outline size="lg" class="full-width" :label="$t('modules.login.loginButton')" color="grey-9" :loading="loading" :disabled="loading" @click="login" />
            </q-card-actions>
            <q-card-section>
                <q-btn
                    v-for="provider of availableOAuthProviders" :key="provider.name"
                    class="full-width q-mb-xs"
                    outline
                    :icon="oauthProviderIconByName(provider.name)"
                    size="sm"
                    color="grey-9"
                    :label="provider.name"
                    @click="startOAuthLogin(provider)"
                />
            </q-card-section>
            <q-card-section class="text-center q-pa-none q-mt-sm">
                <p class="text-grey-6">{{ $t('modules.login.hint') }}</p>
            </q-card-section>
            </q-card>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { LocalStorage, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { ref, onMounted, provide } from 'vue'
import { useLoginStore } from '../../stores/login-store'
import { api } from '../../boot/axios'
import { useRouter, useRoute } from 'vue-router'
import { AvailableOAuthProvider } from 'src/boot/api/login'

const loginStore = useLoginStore()
const $q = useQuasar()
const $i18n = useI18n()
const $router = useRouter()
const $route = useRoute()

const loading = ref(false)
const namespace = ref('')
const namespaceSelectionBlocked = ref(false)
const username = ref('')
const password = ref('')

const availableOAuthProviders = ref([] as AvailableOAuthProvider[])

async function login() {
    if (loading.value) {
        // Prevent multiple logins at same time
        return
    }

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.login.loginOperationPendingNotify')
    })
    loading.value = true

    try {
        const loginResponse = await api.login.createTokensWithPassword({
            namespace: namespace.value,
            login: username.value,
            password: password.value
        })

        loginStore.login(namespace.value, loginResponse.accessToken, loginResponse.refreshToken)

        notif({
          type: 'positive',
          message: $i18n.t('modules.login.successfullyLoggedInNotify'),
          timeout: 1000
        })

        await $router.push({
            name: 'home',
            params: {
                currentNamespace: namespace.value === '' ? '_global' : namespace.value
            }
        })
    } catch (e) {

        notif({
          type: 'negative',
          message: $i18n.t('modules.login.failToLoginNotify', {error: String(e)}),
          timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

async function loadOAuthInfo() {
    try {
        const response = await api.login.getAvailableOAuthProviders({ namespace: namespace.value })
        availableOAuthProviders.value = response.providers
    } catch (e) {
        console.error(e)
    }
}

function oauthProviderIconByName(name: string): string {
    switch (name) {
        case 'github': return 'fa-brands fa-github'
        case 'google': return 'fa-brands fa-google'
        case 'facebook': return 'fa-brands fa-facebook'
        case 'twitter': return 'fa-brands fa-twitter'
        case 'microsoft': return 'fa-brands fa-microsoft'
        case 'gitlab': return 'fa-brands fa-gitlab'
        case 'discord': return 'fa-brands fa-discord'
        case 'apple': return 'fa-brands fa-apple'
    }
    return 'info'
}

async function startOAuthLogin(provider: AvailableOAuthProvider) {
    let callBackURL = $router.resolve({
        name: 'login_oauth_finalize',
        params: {
            provider: provider.name,
            currentNamespace: namespace.value === '' ? '_global' : namespace.value
        }
    }).href
    callBackURL = encodeURIComponent(`${window.location.origin}/${callBackURL}`)
    
    const stateBytes = new Uint32Array(16)
    self.crypto.getRandomValues(stateBytes)
    const state = encodeURIComponent(Array.from(stateBytes).map(b => b.toString(16).padStart(2, '0')).join(''))
    LocalStorage.set('login_oauth_state', state)

    const redirectLocation = `${provider.authUrl}?client_id=${provider.clientId}&state=${state}&redirect_uri=${callBackURL}`
    window.location.href = redirectLocation
} 

onMounted(async () => {
    const currentNamespace = $route.params.currentNamespace
    if (currentNamespace !== undefined && currentNamespace !== '_global') {
        namespace.value = currentNamespace as string
        namespaceSelectionBlocked.value = true
    }

    loginStore.tryLoadLoginFromStorage()
    if (loginStore.isLoggedIn) {
        await $router.push({
            name: 'home',

            params: {
                currentNamespace: loginStore.originalNamespace === '' ? '_global' : loginStore.originalNamespace
            }
        })
        return
    }

    await loadOAuthInfo()
})

</script>

<style>
.q-card {
  width: 500px;
}
</style>