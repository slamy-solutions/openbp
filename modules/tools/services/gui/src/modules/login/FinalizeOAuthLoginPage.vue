<template>
    <q-spinner></q-spinner>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { LocalStorage, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from 'src/boot/api';
import { onMounted } from 'vue';
import { useLoginStore } from 'src/stores/login-store';

const $q = useQuasar()
const $route = useRoute()
const $router = useRouter()
const $i18n = useI18n()

const loginStore = useLoginStore()

const currentNamespace = $route.params.currentNamespace === "_global" ?  '' : $route.params.currentNamespace as string
// Fix vue-router bug
function findGetParameter(parameterName: string): string {
    var result = null,
        tmp = [];
    location.search
        .substr(1)
        .split("&")
        .forEach(function (item) {
          tmp = item.split("=");
          if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
        });
    return result as unknown as string;
}

async function handleAuth() {
    const code = findGetParameter('code')
    const state = findGetParameter('state')

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.login.oauth.finalize.finalizeOperationNotify'),
    })

    const savedState = LocalStorage.getItem('login_oauth_state') as string
    if (savedState !== state) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.login.oauth.finalize.finalizeInvalidStateNotify'),
            timeout: 5000
        })
        await $router.push({ name: 'login', params: { currentNamespace: $route.params.currentNamespace }})
    }


    try {
        const response = await api.login.createTokenWithOAuth({
            namespace: currentNamespace,
            provider: $route.params.provider as string,
            code,
        })

        loginStore.login(currentNamespace, response.accessToken, response.refreshToken)

        notif({
            type: 'positive',
            message: $i18n.t('modules.login.oauth.finalize.finalizeSuccessNotify'),
            timeout: 1000
        })
        await $router.push({
            name: 'home',
            params: {
                currentNamespace: $route.params.currentNamespace
            }
        })
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.login.oauth.finalize.finalizeFailNotify', { error: e }),
            timeout: 5000
        })
        await $router.push({
            name: 'login',
            params: {
                currentNamespace: $route.params.currentNamespace
            }
        })
    } finally {
        await $router.push({ name: 'login' })
    }
}

onMounted(async () => {
    await handleAuth()
})

</script>