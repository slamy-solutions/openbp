<template>
    <div class="bg-primary window-height window-width row justify-center items-center">
        <div class="column">
        <div class="row items-center justify-center q-mb-lg">
            <q-img src="logo.svg" width="350px"></q-img>
        </div>
        <div class="row">
            <q-card square bordered class="q-pa-lg shadow-1">
            <q-card-section>
                <q-form class="q-gutter-md">
                <q-input square filled clearable v-model="username" type="text" :label="$t('modules.login.usernameInput')" />
                <q-input square filled clearable v-model="password" type="password" :label="$t('modules.login.passwordInput')" />
                </q-form>
            </q-card-section>
            <q-card-actions class="q-px-md">
                <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.login.loginButton')" :loading="loading" :disabled="loading" @click="login" />
            </q-card-actions>
            <q-card-section class="text-center q-pa-none">
                <p class="text-grey-6">{{ $t('modules.login.hint') }}</p>
            </q-card-section>
            </q-card>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { ref, onMounted } from 'vue'
import { useLoginStore } from '../../stores/login-store'
import { api } from '../../boot/axios'
import { useRouter } from 'vue-router'

const loginStore = useLoginStore()
const $q = useQuasar()
const $i18n = useI18n()
const $router = useRouter()

const loading = ref(false)
const username = ref('')
const password = ref('')

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
            login: username.value,
            password: password.value
        })

        loginStore.login(username.value, loginResponse.accessToken, loginResponse.refreshToken)

        notif({
          type: 'positive',
          message: $i18n.t('modules.login.successfullyLoggedInNotify'),
          timeout: 1000
        })

        await $router.push({
            name: 'home'
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

onMounted(async () => {
    loginStore.tryLoadLoginFromStorage()
    if (loginStore.isLoggedIn) {
        await $router.push({
            name: 'home'
        })
    }
})

</script>

<style>
.q-card {
  width: 500px;
}
</style>