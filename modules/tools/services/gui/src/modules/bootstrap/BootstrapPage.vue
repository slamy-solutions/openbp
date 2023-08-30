<template>
    <div class="bg-primary window-height window-width row justify-center items-center">
        <div class="column">
        <div class="row items-center justify-center q-mb-lg">
            <q-img src="logo.svg" width="350px"></q-img>
        </div>
        <div class="row items-center justify-center">
            <q-card square bordered class="q-pa-lg shadow-1">
            <q-card-section class="text-center q-pa-none">
                <h3 class="q-ma-lg">{{ $t('modules.bootstrap.header') }}</h3>
                <p class="text-grey-6">{{ $t('modules.bootstrap.subheader') }}</p>
            </q-card-section>
            <q-stepper
                v-model="step"
                ref="stepper"
                active-color="secondary"
                inactive-color="dark"
                done-color="primary"
                vertical
                animated
            >

                <q-step
                    icon="settings"
                    :name="0"    
                    :title="$t('modules.bootstrap.steps.status.label')"
                >
                    <q-spinner size="sm"/>
                </q-step>

                <q-step
                    icon="key"
                    :name="1"
                    :title="$t('modules.bootstrap.steps.vault.label')"
                >
                    <q-card-section>
                        <q-form class="q-gutter-md">
                            <q-input square filled clearable v-model="unsealPassword" type="password" :label="$t('modules.bootstrap.steps.vault.passwordInput')" />
                        </q-form>
                    </q-card-section>
                    <q-card-actions class="q-px-md">
                        <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.bootstrap.steps.vault.unsealButton')" :loading="loading" :disabled="loading" @click="unsealVault" />
                    </q-card-actions>
                    <q-card-section class="text-center q-pa-none">
                        <p class="text-grey-6">{{ $t('modules.bootstrap.steps.vault.unsealHint') }}</p>
                    </q-card-section>
                </q-step>

                <q-step
                    icon="person"
                    :name="2"
                    
                    :title="$t('modules.bootstrap.steps.rootUser.label')"
                >   
                    <div v-if="!bootstrapStore.rootUserCreationBlocked">
                        <q-card-section>
                        <q-form class="q-gutter-md">
                            <q-input square filled clearable v-model="rootUsername" type="text" :label="$t('modules.bootstrap.steps.rootUser.usernameInput')" />
                            <q-input square filled clearable v-model="rootPassword" type="password" :label="$t('modules.bootstrap.steps.rootUser.passwordInput')" />
                        </q-form>
                        </q-card-section>
                        <q-card-actions class="q-px-md">
                            <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.bootstrap.steps.rootUser.createButton')" :loading="loading" :disabled="loading" @click="createRootUser" />
                        </q-card-actions>
                        <q-card-section class="text-center q-pa-none">
                            <p class="text-grey-6">{{ $t('modules.bootstrap.steps.rootUser.createHint') }}</p>
                        </q-card-section>
                    </div>
                    <div v-else>
                        <q-card-section class="text-center q-pa-none">
                            <h4 class="q-ma-lg text-negative">{{ $t('modules.bootstrap.steps.rootUser.blockedHeader') }}</h4>
                            <p class="text-negative">{{ $t('modules.bootstrap.steps.rootUser.blockedHint') }}</p>
                        </q-card-section>
                    </div>
                </q-step>

            </q-stepper>
            </q-card>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { ref, onMounted } from 'vue'
import { useBootstrapStore } from '../../stores/bootstrap-store'
import { useLoginStore } from  '../../stores/login-store' 
import { api } from '../../boot/axios'
import { useRouter } from 'vue-router'

const step = ref(2)

const bootstrapStore = useBootstrapStore()
const loginStore = useLoginStore()
const $router = useRouter()
const $q = useQuasar()
const $i18n = useI18n()

const loading = ref(false)

const unsealPassword = ref('')

const rootUsername = ref('')
const rootPassword = ref('')

async function loadStatus() {
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.bootstrap.steps.status.getStatusOperationNotify')
    })

    try {
        const status = await api.bootstrap.getStatus()
        bootstrapStore.updateBootstrapState(status)
        if (status.vaultSealed) {
            step.value = 1
        } else if (!status.rootUserCreated) {
            step.value = 2
        } else {
            await $router.push({ name: "login", params: { currentNamespace: '_global' } })
        }

        notif()
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.bootstrap.steps.status.getStatusFailNotify', {error: String(e)}),
            timeout: 5000
        })
    }
}

async function unsealVault() {
    //TODO: add functionality
    await loadStatus()
}

async function createRootUser() {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.bootstrap.steps.rootUser.createOperationNotify')
    })


    try {
        await api.bootstrap.createRootUser({
            login: rootUsername.value,
            password: rootPassword.value
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.bootstrap.steps.rootUser.createSuccessNotify'),
            timeout: 5000
        })

        // Clear old login data if it presents (mostly for dev)
        loginStore.logout()
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.bootstrap.steps.rootUser.createFailNotify', {error: String(e)}),
            timeout: 5000
        })
    } finally {
        loading.value = false
        await loadStatus()
    }
}

onMounted(async () => {
    if (bootstrapStore.bootstrapped) {
        await $router.push({name: 'login', params: { currentNamespace: '_global' }})
    } else {
        await loadStatus()
    }
})

</script>

<style>
.q-card {
  width: 90%;
  max-width: 900px;
}
</style>