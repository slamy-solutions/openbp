<template>
    <q-page padding>
        <q-card flat class="bg-transparent" style="max-width: 600px;">
            <div class="text-h6">{{ $t('modules.me.auth.password.header') }}</div>
            <q-separator />
            <q-card-section class="q-gutter-sm">
                <div class="text-caption">{{ $t('modules.me.auth.password.caption') }}</div>
                <q-input v-model="newPasswordInput" :label="$t('modules.me.auth.password.passwordInput')" type="password" outlined square class="fit" :loading="loading"/>
                
                <q-btn outline :label="$t('modules.me.auth.password.setOrUpdateButton')" class="fit" :loading="loading" :disable="loading" @click="updatePassword" />
                
                <div v-if="!passwordEnabled" class="text-center">{{ $t('modules.me.auth.password.disabledMessage') }}</div>
                <q-btn v-if="passwordEnabled" outline :label="$t('modules.me.auth.password.disableButton')" class="fit" color="negative" :loading="loading" :disable="loading" @click="disablePassword" />
                <div v-if="passwordEnabled" class="text-grey text-caption text-center">{{ $t('modules.me.auth.password.disableCaption') }}</div>
            </q-card-section>
            
            <div class="text-h6">Third-party services</div>
            <q-separator />
            <q-card-section class="q-gutter-sm" >
                <div class="text-caption">You can bind thrid-party services to your account to be able to login with them. List of the available services can be configured by the namespace administrator</div>
                <div class="text-h6 text-bold text-center text-info" v-if="availableOAuthProviders.length == 0">There are no configured thrid-party services in this namespace.</div>
            
                <q-list separator >
                    <q-item v-for="provider in availableOAuthProviders" :key="provider.name" style="border: #ccc solid 1px;" :clickable="!isOAuthProviderConfigured(provider.name)" :v-ripple="!isOAuthProviderConfigured(provider.name)" @click="startOAuthRegistration(provider)">
                        <q-item-section avatar>
                            <q-icon :name="oauthProviderIconByName(provider.name)" :color="isOAuthProviderConfigured(provider.name) ? 'positive':''" size="sm"></q-icon>
                        </q-item-section>
                        <q-item-section>
                            <q-item-label>{{ provider.name }}</q-item-label>
                        </q-item-section>
                        <q-item-section side>
                            <q-btn v-if="isOAuthProviderConfigured(provider.name)" icon="delete" size="sm" flat round></q-btn>
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-card-section>

            <div class="text-h6">Two-factor authentification</div>
            <q-separator />
        </q-card>
        
    </q-page>
</template>

<script setup lang="ts">
import api from 'src/boot/api';
import { AvailableOAuthProvider, ConfiguredOAuthProvider } from 'src/boot/api/me/auth';
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { LocalStorage, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'

const $q = useQuasar() 
const $i18n = useI18n()
const $router = useRouter()

const loading = ref(false)

// Password
const passwordEnabled = ref(false)
const newPasswordInput = ref('')

async function updatePassword() {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.me.auth.password.updateOperationNotify'),
    })
    try {
        await api.me.auth.createOrUpdatePassword({
            password: newPasswordInput.value
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.me.auth.password.updateSuccessNotify'),
            timeout: 5000
        })
        newPasswordInput.value = ''
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.me.auth.password.updateFailNotify', { error: e }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

async function disablePassword() {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.me.auth.password.deleteOperationNotify'),
    })
    try {
        await api.me.auth.deletePassword()
        notif({
            type: 'positive',
            message: $i18n.t('modules.me.auth.password.deleteSuccessNotify'),
            timeout: 5000
        })
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.me.auth.password.deleteFailNotify', { error: e }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

// OAuth
const availableOAuthProviders = ref([] as Array<AvailableOAuthProvider>)
const configuredOAuthProviders = ref([] as Array<ConfiguredOAuthProvider>)

function isOAuthProviderConfigured(name: string): boolean {
    return configuredOAuthProviders.value.findIndex(p => p.name === name) !== -1
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

async function deleteOAuthRegistration(provider: string) {

}


async function startOAuthRegistration(provider: AvailableOAuthProvider) {
    let callBackURL = $router.resolve({
        name: 'me_auth_oauth_finalize',
        params: {
            provider: provider.name
        }
    }).href
    callBackURL = encodeURIComponent(`${window.location.origin}/${callBackURL}`)
    
    const stateBytes = new Uint32Array(16)
    self.crypto.getRandomValues(stateBytes)
    const state = encodeURIComponent(Array.from(stateBytes).map(b => b.toString(16).padStart(2, '0')).join(''))
    LocalStorage.set('me_auth_oauth_state', state)

    const redirectLocation = `${provider.authUrl}?client_id=${provider.clientId}&state=${state}&redirect_uri=${callBackURL}`
    window.location.href = redirectLocation
}

async function loadData() {
    loading.value = true
    try {
        const response = await api.me.auth.getInfo()
        passwordEnabled.value = response.password.enabled
        availableOAuthProviders.value = response.oauth.availableProviders
        configuredOAuthProviders.value = response.oauth.configuredProviders
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadData()
})

</script>